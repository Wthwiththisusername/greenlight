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
			Email:    "zh.arumandes@gmail.com",
			Password: "pa55word",
			wantCode: http.StatusCreated,
		},
		{
			name:     "test for wrong input",
			Email:    "zh.arumandes@gmail.com",
			Password: "pa55word",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "failed Validation",
			Email:    "zh.arumandesgmail.com",
			Password: "pa55word",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "email not found",
			Email:    "notfound@gmail.com",
			Password: "pa55word",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "password didn't match",
			Email:    "zh.arumandes@gmail.com",
			Password: "pa55word1",
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
