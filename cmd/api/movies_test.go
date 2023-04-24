package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestShowMovie(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid",
			urlPath:  "/v1/movies/1",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent",
			urlPath:  "/v1/movies/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative",
			urlPath:  "/v1/movies/-5",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal",
			urlPath:  "/v1/movies/1.28",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String",
			urlPath:  "/v1/movies/foo",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}

}

func TestCreateMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const (
		validTitle   = "Test Title"
		validYear    = 2023
		validRuntime = "98 mins"
	)

	validGenres := []string{"comedy", "drama"}

	tests := []struct {
		name     string
		Title    string
		Year     int32
		Runtime  string
		Genres   []string
		wantCode int
	}{
		{
			name:     "Valid",
			Title:    validTitle,
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusCreated,
		},
		{
			name:     "Empty",
			Title:    "",
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Unvalid year",
			Title:    validTitle,
			Year:     1500,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "wrong input",
			Title:    validTitle,
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Title   string   `json:"title"`
				Year    int32    `json:"year"`
				Runtime string   `json:"runtime"`
				Genres  []string `json:"genres"`
			}{
				Title:   tt.Title,
				Year:    tt.Year,
				Runtime: tt.Runtime,
				Genres:  tt.Genres,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.postForm(t, "/v1/movies", b)

			assert.Equal(t, code, tt.wantCode)

		})
	}
}

func TestDeleteMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "deleting existing movie",
			urlPath:  "/v1/movies/1",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/v1/movies/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "invalid ID",
			urlPath:  "/v1/movies/s",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.deleteReq(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}

}

func TestUpdateMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const (
		validTitle   = "Test Title"
		validYear    = 2023
		validRuntime = "98 mins"
	)
	validGenres := []string{"comedy", "drama"}

	tests := []struct {
		name     string
		url      string
		Title    string
		Year     int32
		Runtime  string
		Genres   []string
		wantCode int
		wantBody string
	}{
		{
			name:     "Non-existent ID",
			url:      "/v1/movies/2",
			Title:    validTitle,
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "invalid ID",
			url:      "/v1/movies/s",
			Title:    "",
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "wrong input",
			url:      "/v1/movies/1",
			Title:    validTitle,
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "bad json",
			url:      "/v1/movies/1",
			Title:    validTitle,
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty Title",
			url:      "/v1/movies/1",
			Title:    "",
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "unvalid year",
			url:      "/v1/movies/1",
			Title:    validTitle,
			Year:     1500,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Valid",
			url:      "/v1/movies/1",
			Title:    validTitle,
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusOK,
			wantBody: "{\"movie\":{\"id\":1,\"title\":\"Test Title\",\"year\":2021,\"runtime\":\"105 mins\",\"genres\":[\"comedy\",\"drama\"],\"version\":0}}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Title   string   `json:"title"`
				Year    int32    `json:"year"`
				Runtime string   `json:"runtime"`
				Genres  []string `json:"genres"`
			}{
				Title:   tt.Title,
				Year:    tt.Year,
				Runtime: tt.Runtime,
				Genres:  tt.Genres,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "wrong input" {
				b = append(b, 'a')
			} else if tt.name == "bad json" {
				b[1] = ','
			}

			code, _, body := ts.patchReq(t, tt.url, b)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}
}

func TestListMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name    string
		Title   string
		Genres  string
		Filters struct {
			Page     string
			PageSize string
			Sort     string
		}
		wantCode int
		wantBody string
	}{
		{
			name:   "Invalid page input",
			Title:  "Test",
			Genres: "",
			Filters: struct {
				Page     string
				PageSize string
				Sort     string
			}{Page: "p", PageSize: "s", Sort: ""},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:   "Invalid page input page>0",
			Title:  "Test",
			Genres: "",
			Filters: struct {
				Page     string
				PageSize string
				Sort     string
			}{Page: "-1", PageSize: "", Sort: ""},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:   "Valid",
			Title:  "Test",
			Genres: "vampires",
			Filters: struct {
				Page     string
				PageSize string
				Sort     string
			}{Page: "", PageSize: "", Sort: ""},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/v1/movies?title=%s&genres=%s&page=%s&page_size=%s&sort=%s", tt.Title, tt.Genres, tt.Filters.Page, tt.Filters.PageSize, tt.Filters.Sort)

			code, _, body := ts.get(t, url)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}
}
