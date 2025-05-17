package main

import (
	"net/http"
	"project/internal/store"
)

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("Forbidden", "method", r.Method, "path", r.URL.Path)
	writeJSON(w, http.StatusForbidden, "forbidden error")
}

func (app *application) unAuthResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauth error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeJSON(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unAuthBasicResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauth basic error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJSON(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("Internal server error: Method: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	app.jsonResponse(w, http.StatusInternalServerError, "Server error")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Bad request error: Method: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	app.jsonResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Not found error: Method: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	app.jsonResponse(w, http.StatusNotFound, store.ErrNotFound)
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Conflict error: Method %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	app.jsonResponse(w, http.StatusConflict, err.Error())
}
