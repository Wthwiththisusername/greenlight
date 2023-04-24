package main

import (
	"encoding/json"
	"greenlight.bcc/internal/assert"
	"net/http"
	"testing"
)

func TestCreateToken(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		Email    string
		Password string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid",
			Email:    "nargiz@gmail.com",
			Password: "12345",
			wantCode: http.StatusCreated,
		},
		{
			name:     "wrong input",
			Email:    "nargiz@gmail.com",
			Password: "12345",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "failed Validation",
			Email:    "nargiz@gmail.com",
			Password: "12345",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "email not found",
			Email:    "notfound@gmail.com",
			Password: "12345",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "password didn't match",
			Email:    "nargiz@gmail.com",
			Password: "123456",
			wantCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    tt.Email,
				Password: tt.Password,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.postForm(t, "/v1/tokens/authentication", b)

			assert.Equal(t, code, tt.wantCode)

		})
	}

}
