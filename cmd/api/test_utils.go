package main

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"project/internal/auth"
	"project/internal/store"
	"project/internal/store/cache"
	"testing"
)

func newTestApp(t *testing.T) *application {
	t.Helper()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockCacheStorage()
	testAuth := &auth.TestAuth{}
	logger := zap.NewNop().Sugar()
	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCache,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
