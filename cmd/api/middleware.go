package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) BasicAuthMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unAuthBasicResponse(w, r, fmt.Errorf("missing Authorization header"))
				return
			}
			// parse, decode it
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unAuthBasicResponse(w, r, fmt.Errorf("invalid Authorization header"))
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unAuthBasicResponse(w, r, err)
				return
			}

			creds := strings.SplitN(string(decoded), ":", 2)
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unAuthBasicResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			// check the creds

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unAuthResponse(w, r, fmt.Errorf("missing Authorization header"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unAuthResponse(w, r, fmt.Errorf("invalid Authorization header"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unAuthResponse(w, r, err)
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unAuthResponse(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userId)
		if err != nil {
			app.unAuthResponse(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
