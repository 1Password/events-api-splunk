package actions

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.1password.io/eventsapi-splunk/api"
	"go.1password.io/eventsapi-splunk/store"
)

func StartAuditEvents(cursorFile string, limit int, startAt *time.Time, eventsAPI *api.EventsAPI) error {
	log.Println("Starting AuditEvents...")

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
		log.Println("Performing cursor reset")
		body := api.CursorResetRequest{
			Limit:     limit,
			StartTime: startAt,
		}
		res, err := eventsAPI.AuditEventsRequest(ctx, body)
		if err != nil {
			return fmt.Errorf("AuditEventsRequest request failed: %w", err)
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
		log.Println("Using stored cursor")
	}

	for {
		select {
		case <-sigCh:
			log.Println("Interrupted, shutting down")
			cancel()
			err := store.CloseStore()
			if err != nil {
				return fmt.Errorf("could not close store: %w", err)
			}
			return nil
		default:
			body := api.CursorRequest{Cursor: cursor}
			res, err := eventsAPI.AuditEventsRequest(ctx, body)
			if err != nil {
				log.Printf("AuditEventsRequest request failed: %s\n", err)
				time.Sleep(30 * time.Second)
				continue
			}

			if len(res.AuditEvents) == 0 && !res.HasMore {
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
