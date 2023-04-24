package main

import (
	"encoding/json"
	"greenlight.bcc/internal/assert"
	"net/http"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		Username string
		Email    string
		Password string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid name",
			Username: "aizhan",
			Email:    "aizhan@gmail.com",
			Password: "12345678",
			wantCode: http.StatusCreated,
		},
		{
			name:     "Wrong input",
			Username: "zhamal",
			Email:    "zhamal@gmail.com",
			Password: "12345678",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "inValid name",
			Username: "",
			Email:    "unknown@gmail.com",
			Password: "12345678",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "inValid email",
			Username: "lashyn",
			Email:    "",
			Password: "12345678",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "inValid password",
			Username: "aizhan",
			Email:    "aizhan@gmail.com",
			Password: "123456",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Duplicated",
			Username: "nargiz",
			Email:    "nargizazat7@gmail.com",
			Password: "12345678",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Name:     tt.Username,
				Email:    tt.Email,
				Password: tt.Password,
			}
			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "Wrong input" {
				b = append(b, 'a')
			}

			code, _, body := ts.postForm(t, "/v1/users", b)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}
}

func TestActivateUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		Token    string
		wantCode int
		wantBody string
	}{
		{
			name:     "inValid",
			Token:    "12",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "wrong input",
			Token:    "12",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Valid",
			Token:    "12345678912345678912345678",
			wantCode: http.StatusOK,
		},
		{
			name:     "ErrRecordNotFound",
			Token:    "12345678912345678912345679",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "cant update",
			Token:    "12345678912345678912345677",
			wantCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := struct {
				Token string `json:"token"`
			}{
				Token: tt.Token,
			}

			b, err := json.Marshal(&input)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "wrong input" {
				b = append(b, 'a')
			}

			code, _, body := ts.putReq(t, "/v1/users/activated", b)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}
}
