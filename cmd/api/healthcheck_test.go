package main

import (
	"greenlight.bcc/internal/assert"
	"net/http"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid",
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.get(t, "/v1/healthcheck")

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}
}
