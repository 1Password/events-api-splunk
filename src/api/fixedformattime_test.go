package api_test

import (
	"testing"
	"time"

	"go.1password.io/eventsapi-splunk/api"
)

func TestFixedFormatTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		t       api.FixedFormatTime
		want    string
		wantErr bool
	}{
		{
			name:    "Marshal without sub-second",
			t:       api.FixedFormatTime(mustParseTime(time.RFC3339Nano, "2021-07-07T13:51:12Z")),
			want:    `"2021-07-07T13:51:12Z"`,
			wantErr: false,
		},
		{
			name:    "Marshal with sub-second",
			t:       api.FixedFormatTime(mustParseTime(time.RFC3339Nano, "2021-07-07T13:51:12.623Z")),
			want:    `"2021-07-07T13:51:12Z"`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got := string(got); got != tt.want {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixedFormatTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		t       api.FixedFormatTime
		json    string
		wantErr bool
	}{
		{
			name:    "Unmarshal without sub-second",
			t:       api.FixedFormatTime{},
			json:    `"2021-07-07T13:51:12Z"`,
			wantErr: false,
		},
		{
			name:    "Unmarshal with sub-second",
			t:       api.FixedFormatTime{},
			json:    `"2021-07-07T13:51:12.623Z"`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalJSON([]byte(tt.json)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func mustParseTime(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}
