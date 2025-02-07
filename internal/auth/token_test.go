package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid Bearer Token",
			headers:       http.Header{"Authorization": {"Bearer validtoken"}},
			expectedToken: "validtoken",
			expectError:   false,
		},
		{
			name:        "No Authorization Header",
			headers:     http.Header{},
			expectError: true,
		},
		{
			name:        "Invalid Authorization Header",
			headers:     http.Header{"Authorization": {"InvalidHeader"}},
			expectError: true,
		},
		{
			name:        "Missing Bearer Token",
			headers:     http.Header{"Authorization": {"Bearer "}},
			expectError: true,
		},
		{
			name:        "Non-Bearer Token",
			headers:     http.Header{"Authorization": {"Basic dXNlcjpwYXNz"}},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if token != tt.expectedToken {
				t.Errorf("expected token: %s, got: %s", tt.expectedToken, token)
			}
		})
	}
}
