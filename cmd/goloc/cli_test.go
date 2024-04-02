package main

import (
	"log"
	"testing"
)

func TestValidateCoordinates(t *testing.T) {
	gotLatitude, gotLongitude, err := validateCoordinates("22.23354166,22.54223445")
	
	if err != nil {
		log.Fatal("error during validaion")
		t.Errorf("Error occurred during validation: %v", err)
	}
	
	wantedLatitude, wantedLongitude := 22.23354166, 22.54223445

	if gotLatitude != wantedLatitude || gotLongitude != wantedLongitude {
		log.Fatalf("Expected (%f,%f), got (%f,%f)", wantedLatitude, wantedLongitude, gotLatitude, gotLongitude)
		t.Errorf("Expected (%f,%f), got (%f,%f)", wantedLatitude, wantedLongitude, gotLatitude, gotLongitude)
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		input   string
		want    bool
		message string
	}{
		{"2024-04-01T12:00:00+00:00", true, "Valid ISO 8601 date"},
		{"2024-04-01", false, "Invalid date format"},
		{"12:00:00+00:00", false, "Invalid date format"},
		{"2024-04-01T12:00:00", false, "Invalid date format"},
		{"2024-04-01T12:00:00+00:00Z", false, "Invalid date format"},
		{"2024/04/01T12:00:00+00:00", false, "Invalid date format"},
	}

	for _, test := range tests {
		got := validateDate(test.input)
		if got != test.want {
			t.Errorf("%s: got %t, want %t", test.message, got, test.want)
		}
	}
}