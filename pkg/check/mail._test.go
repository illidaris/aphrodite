package check

import (
	"testing"
)

// TestIsValidEmail covers various valid and invalid email cases.
func TestIsValidEmail(t *testing.T) {
	// Define test cases in a table-driven format.
	tests := []struct {
		name     string // Test case name
		email    string // Input email
		expected bool   // Expected result
	}{
		// Valid emails (should return true)
		{"valid_standard", "test@example.com", true},
		{"valid_dot_in_local", "user.name@example.com", true},
		{"valid_plus_in_local", "user+tag@example.org", true},
		{"valid_numeric_local", "123@example.com", true},
		{"valid_dash_in_domain", "user@exa-mple.com", true},
		{"valid_long_tld", "user@example.museum", true},
		{"valid_subdomain", "user@sub.example.com", true},

		// Invalid emails (should return false)
		{"missing_at_symbol", "invalid-email", false},
		{"empty_local_part", "@example.com", false},
		{"empty_string", "", false},
		{"space_prefix", " user@example.com", false},
		{"space_suffix", "user@example.com ", false},
		{"missing_domain", "user@", false},
		{"tld_too_short", "user@example.c", false},
		{"invalid_char_in_domain", "user@exa!mple.com", false},
	}

	// Iterate through each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the function under test.
			got := IsValidEmail(tt.email)
			// Validate the result against expected.
			if got != tt.expected {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.expected)
			}
		})
	}
}
