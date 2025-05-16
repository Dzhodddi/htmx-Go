package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"project/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=24"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler Go docs
//
//	@Summary		Register a user
//	@Description	Registers new users
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()
	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(ctx, user, hashedToken, app.config.mail.exp)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.badRequestError(w, r, err)
		case errors.Is(err, store.ErrDuplicateUsername):
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	userWithToken := &UserWithToken{
		User:  user,
		Token: plainToken,
	}
	//isProdEnv := app.config.env == "production"
	//activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	//vars := struct {
	//	Username      string
	//	ActivationURL string
	//}{
	//	Username:      user.Username,
	//	ActivationURL: activationURL,
	//}
	//var _ int64
	//_, err = app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	//if err != nil {
	//	app.logger.Errorw("Error sending welcome email", "error", err)
	//
	//	//rollback if email fails
	//	if err := app.store.Users.Delete(ctx, user.ID); err != nil {
	//		app.logger.Errorw("Error deleting user", "error", err)
	//		app.internalServerError(w, r, err)
	//	}
	//	app.internalServerError(w, r, err)
	//	return
	//}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}
