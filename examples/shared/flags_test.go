package shared

import (
	"os"
	"testing"
)

func TestGetSerialWithFallback(t *testing.T) {
	tests := []struct {
		name      string
		flagValue string
		envValue  string
		want      string
		wantErr   bool
	}{
		{
			name:      "flag takes precedence",
			flagValue: "FLAG123",
			envValue:  "ENV456",
			want:      "FLAG123",
			wantErr:   false,
		},
		{
			name:      "env var used when flag empty",
			flagValue: "",
			envValue:  "ENV456",
			want:      "ENV456",
			wantErr:   false,
		},
		{
			name:      "error when both empty",
			flagValue: "",
			envValue:  "",
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("BS_SERIAL", tt.envValue)
				defer os.Unsetenv("BS_SERIAL")
			} else {
				os.Unsetenv("BS_SERIAL")
			}

			got, err := GetSerialWithFallback(tt.flagValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSerialWithFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSerialWithFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDeviceIDWithFallback(t *testing.T) {
	tests := []struct {
		name      string
		flagValue int
		envValue  string
		want      int
		wantErr   bool
	}{
		{
			name:      "flag takes precedence",
			flagValue: 123,
			envValue:  "456",
			want:      123,
			wantErr:   false,
		},
		{
			name:      "env var used when flag zero",
			flagValue: 0,
			envValue:  "456",
			want:      456,
			wantErr:   false,
		},
		{
			name:      "error when both empty",
			flagValue: 0,
			envValue:  "",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "error when env var invalid",
			flagValue: 0,
			envValue:  "not-a-number",
			want:      0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("BS_DEVICE_ID", tt.envValue)
				defer os.Unsetenv("BS_DEVICE_ID")
			} else {
				os.Unsetenv("BS_DEVICE_ID")
			}

			got, err := GetDeviceIDWithFallback(tt.flagValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeviceIDWithFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDeviceIDWithFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGroupIDWithFallback(t *testing.T) {
	tests := []struct {
		name      string
		flagValue int
		envValue  string
		want      int
		wantErr   bool
	}{
		{
			name:      "flag takes precedence",
			flagValue: 42,
			envValue:  "99",
			want:      42,
			wantErr:   false,
		},
		{
			name:      "env var used when flag zero",
			flagValue: 0,
			envValue:  "99",
			want:      99,
			wantErr:   false,
		},
		{
			name:      "error when both empty",
			flagValue: 0,
			envValue:  "",
			want:      0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("BS_GROUP_ID", tt.envValue)
				defer os.Unsetenv("BS_GROUP_ID")
			} else {
				os.Unsetenv("BS_GROUP_ID")
			}

			got, err := GetGroupIDWithFallback(tt.flagValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGroupIDWithFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGroupIDWithFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}
