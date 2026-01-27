package actions

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.1password.io/eventsapi-splunk/api"
	"go.1password.io/eventsapi-splunk/store"
)

func StartSignIns(cursorFile string, limit int, startAt *time.Time, eventsAPI *api.EventsAPI) error {
	store, err := store.OpenStore(cursorFile)
	defer store.CloseStore()
	if err != nil {
		return fmt.Errorf("could not read file log file: %w", err)
	}

	cursor, err := store.GetCursor()
	if err != nil {
		return fmt.Errorf("could not read file log file: %w", err)
	}

	// Sets ups notify channel so that we can gracefully shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if cursor == "" {
		slog.Debug("Performing cursor reset")
		body := api.CursorResetRequest{
			Limit:     limit,
			StartTime: startAt,
		}
		res, err := eventsAPI.SignInAttemptsRequest(ctx, body)
		if err != nil {
			return fmt.Errorf("SignInCursorResetRequest request failed: %w", err)
		}
		err = res.PrintEvents()
		if err != nil {
			return fmt.Errorf("PrintEvents failed: %w", err)
		}
		err = store.SaveCursor(res.Cursor)
		if err != nil {
			return fmt.Errorf("could not save cursor: %w", err)
		}
		cursor = res.Cursor
	} else {
		slog.Debug("Using stored cursor")
	}

	for {
		select {
		case <-sigCh:
			slog.Info("Interrupted, shutting down")
			cancel()
			err := store.CloseStore()
			if err != nil {
				return fmt.Errorf("could not close store: %w", err)
			}
			return nil
		default:
			body := api.CursorRequest{Cursor: cursor}
			res, err := eventsAPI.SignInAttemptsRequest(ctx, body)
			if err != nil {
				slog.Error("SignInCursorRequest request failed", "error", err)
				time.Sleep(30 * time.Second)
				continue
			}

			if len(res.Items) == 0 && !res.HasMore {
				// Don't bother printing or storing this cursor,
				// we will reuse the last one until we receive some events
				time.Sleep(10 * time.Second)
				continue
			}

			err = res.PrintEvents()
			if err != nil {
				return fmt.Errorf("PrintEvents failed: %w", err)
			}
			err = store.SaveCursor(res.Cursor)
			if err != nil {
				return fmt.Errorf("SaveCursor failed: %w", err)
			}
			cursor = res.Cursor
		}
	}
}
