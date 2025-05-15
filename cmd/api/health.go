package main

import (
	"net/http"
)

// HealthCheckAPI  godoc
//
//	@Summary		Health check
//	@Description	Health check
//	@Tags			health
//	@Produce		json
//	@Success		204	{object}	string	"OK"
//	@Failure		500	{object}	error
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "OK",
		"version": "1.0.0",
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
