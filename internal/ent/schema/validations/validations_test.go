package validations

import (
	"errors"
	"testing"
)

func TestPorName(t *testing.T) {
	testCases := []struct {
		name          string
		portName      string
		expectedError error
	}{
		{"empty name", "", ErrPortNameLength},
		{"name > 15 chars", "port123456789011", ErrPortNameLength},
		{"valid name", "porty", nil},
		{"name with hyphen", "porty-123", nil},
		{"name with hyphen at beginning", "-porty", ErrPortNameHyphens},
		{"name with hyphen at end", "porty-", ErrPortNameHyphens},
		{"name with adjacent hyphens", "porty--123", ErrPortNameAdjacentHyphens},
		{"name with no letters", "123456789012345", ErrPortNameOneLetter},
		{"name with invalid chars", "porty!", ErrPortNameInvalidChars},
		{"name with no chars", "8080", ErrPortNameOneLetter},
	}

	for _, tt := range testCases {
		// go vet
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := PortName(tt.portName)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
