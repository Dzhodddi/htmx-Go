package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	SortBy string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	q := r.URL.Query()

	limit := q.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = l
	}

	offset := q.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = o
	}

	sortBy := q.Get("sortBy")
	if sortBy != "" {
		fq.SortBy = sortBy
	}

	tags := q.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := q.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := q.Get("since")
	if since != "" {
		fq.Since = parseTime(since)
	}

	until := q.Get("until")
	if until != "" {
		fq.Until = parseTime(until)
	}
	return fq, nil
}

func parseTime(s string) string {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}
