package handlers

import (
	"testing"
	"time"
)

func TestRemainingDuration(t *testing.T) {
	tests := []struct {
		name             string
		epoch            int64
		expectedDuration time.Duration
	}{
		{
			name:             "Time in the past",
			epoch:            time.Now().Add(-10 * time.Minute).Unix(),
			expectedDuration: -10 * time.Minute,
		},
		{
			name:             "Time in the future",
			epoch:            time.Now().Add(10 * time.Minute).Unix(),
			expectedDuration: 10 * time.Minute,
		},
		{
			name:             "Current time",
			epoch:            time.Now().Unix(),
			expectedDuration: 0,
		},
	}

	tolerance := 1 * time.Second // Allow up to 1 second of tolerance

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualDuration := remainingDuration(tt.epoch)
			if actualDuration > tt.expectedDuration+tolerance || actualDuration < tt.expectedDuration-tolerance {
				t.Errorf("Got %v, expected approximately %v", actualDuration, tt.expectedDuration)
			}
		})
	}
}

func TestGetRowColor(t *testing.T) {
	// Mock remainingDuration for testing
	oldRemainingDuration := remainingDuration
	remainingDuration = func(epoch int64) time.Duration {
		return time.Unix(epoch, 0).Sub(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	}
	defer func() {
		remainingDuration = oldRemainingDuration
	}()

	currentTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

	tests := []struct {
		name     string
		epoch    int64
		expected string
	}{
		{"zero", 0, ""},
		{"expired", currentTime - 10, "red-row"},
		{"1 day", currentTime + 1*24*60*60, "red-row"},       // Converted 1 day to seconds
		{"2 days", currentTime + 2*24*60*60, "red-row"},      // Converted 2 days to seconds
		{"3 days", currentTime + 3*24*60*60, "red-row"},      // Converted 3 days to seconds
		{"4 days", currentTime + 4*24*60*60, "orange-row"},   // Converted 4 days to seconds
		{"5 days", currentTime + 5*24*60*60, "orange-row"},   // Converted 5 days to seconds
		{"29 days", currentTime + 29*24*60*60, "orange-row"}, // Converted 29 days to seconds
		{"30 days", currentTime + 30*24*60*60, "orange-row"}, // Converted 30 days to seconds
		{"31 days", currentTime + 31*24*60*60, "yellow-row"}, // Converted 31 days to seconds
		{"59 days", currentTime + 59*24*60*60, "yellow-row"}, // Converted 59 days to seconds
		{"60 days", currentTime + 60*24*60*60, "yellow-row"}, // Converted 60 days to seconds
		{"61 days", currentTime + 61*24*60*60, ""},           // Converted 61 days to seconds
	}

	for _, tt := range tests {
		actual := getRowColor(tt.epoch)
		if actual != tt.expected {
			t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, actual)
		}
	}
}

func TestEpochToHumanReadable(t *testing.T) {
	// Mock remainingDuration for testing
	oldRemainingDuration := remainingDuration
	remainingDuration = func(epoch int64) time.Duration {
		return time.Unix(epoch, 0).Sub(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	}
	defer func() {
		remainingDuration = oldRemainingDuration
	}()

	tests := []struct {
		name     string
		epoch    int64
		expected string
	}{
		{"zero", 0, "-"},
		{"expired", time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC).Unix(), "now"},
		{"1 day 2 hours 3 minutes 4 seconds", time.Date(2023, 1, 2, 2, 3, 4, 0, time.UTC).Unix(), "1 days, 2 hours, 3 minutes, 4 seconds"},
		{"5 days", time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC).Unix(), "5 days"},
		{"30 seconds", time.Date(2023, 1, 1, 0, 0, 30, 0, time.UTC).Unix(), "30 seconds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := epochToHumanReadable(tt.epoch)
			if actual != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, actual)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		format   string
		expected string
	}{
		{
			"Zero Time",
			time.Time{},
			"2006-01-02",
			"-",
		},
		{
			"Unix zero time",
			time.Unix(0, 0),
			"2006-01-02",
			"-",
		},
		{
			"Valid Time",
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			"2006-01-02",
			"2022-01-01",
		},
		{
			"Custom Format",
			time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			"02-Jan-2006",
			"31-Dec-2023",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := formatTime(tt.time, tt.format)
			if actual != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, actual)
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	testCases := []struct {
		Name        string
		baseTplStr  string
		tplStr      string
		data        interface{}
		expectedStr string
		expectedErr error
	}{
		{
			Name:       "Valid_Case",
			baseTplStr: "{{formatTime .Time \"2006-01-02\"}} {{humanReadable .Epoch}} {{getRowColor .Status}}",
			tplStr:     "",
			data: map[string]interface{}{
				"Time":   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				"Epoch":  int64(1630454400),
				"Status": int64(1630454400),
			},
			expectedStr: "2022-01-01 now red-row",
			expectedErr: nil,
		},
		{
			Name:        "Error_in_base_template",
			baseTplStr:  "{{ .MissingFunction }}",
			tplStr:      "",
			data:        nil,
			expectedStr: "<no value>",
			expectedErr: nil, // The template package often does not return errors for missing fields/functions
		},
		{
			Name:        "Error_in_tplStr_parsing",
			baseTplStr:  "{{ .Value }}",
			tplStr:      "{{ .MissingFunction }}",
			data:        nil,
			expectedStr: "<no value>",
			expectedErr: nil,
		},
		{
			Name:       "Error_during_template_execution",
			baseTplStr: "{{ .MissingField }}",
			tplStr:     "",
			data: map[string]interface{}{
				"ExistingField": "value",
			},
			expectedStr: "<no value>",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := renderTemplate(tc.baseTplStr, tc.tplStr, tc.data)
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}
			}
			if result != tc.expectedStr {
				t.Errorf("expected output %q, got %q", tc.expectedStr, result)
			}
		})
	}
}
