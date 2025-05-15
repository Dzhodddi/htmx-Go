package main

import (
	"errors"
	"net/http"
	"project/internal/store"
	"time"
)

// Get user feed godoc
//
//	@Summary		Fetches user feed
//	@Description	Fetches user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := time.Now().Format("2006-01-02")
	dateWeekBefore := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	filterDefault := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		SortBy: "desc",
		Since:  dateWeekBefore,
		Until:  now,
		Tags:   []string{},
	}
	filterQuery, err := filterDefault.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(filterQuery); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if filterQuery.Since == "" || filterQuery.Until == "" {
		app.badRequestError(w, r, errors.New("Since or Until provided in bad format"))
		return
	}
	posts, err := app.store.Posts.GetUserFeed(ctx, int64(9), filterQuery)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, 200, posts); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
