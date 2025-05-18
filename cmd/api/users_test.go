package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := newTestApp(t, config{})
	mux := app.mount()
	testToken, _ := app.authenticator.GenerateToken(nil)
	t.Run("should now allowed unath requests",
		func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := executeRequest(req, mux)

			checkResponseCode(t, http.StatusUnauthorized, rr.Code)
		})

	t.Run("should allowed auth requests",
		func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer "+testToken)
			rr := executeRequest(req, mux)

			checkResponseCode(t, http.StatusOK, rr.Code)
		})
}
