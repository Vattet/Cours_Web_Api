package service

import (
	"log"
	"math/rand/v2"
	"net/url"
	"strconv"

	"rest-api/model"
)

const (
	defaultLimit = 10
	maxLimit     = 100
	gamesPath    = "/api/v1/games"
)

type GameReader interface {
	FindAll(limit, skip int64, params map[string]string, version model.FilterVersion) ([]model.Game, error)
	Count(params map[string]string, version model.FilterVersion) (int64, error)
}

func ListGames(r GameReader, params map[string]string, limit, skip int64) (model.Page, error) {

	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {

		limit = maxLimit
	}
	if skip < 0 {
		skip = 0
	}

	version := model.FilterVersions[rand.IntN(len(model.FilterVersions))]
	log.Printf("ListGames: filtre %s", version)

	total, err := r.Count(params, version)
	if err != nil {
		return model.Page{}, err
	}

	games, err := r.FindAll(limit, skip, params, version)
	if err != nil {
		return model.Page{}, err
	}

	page := model.Page{
		Items:      games,
		Page:       skip/limit + 1,
		PageSize:   limit,
		TotalItems: total,

		TotalPages: (total + limit - 1) / limit,
		Links:      model.Links{Self: link(params, limit, skip)},
	}

	if skip+limit < total {
		page.Links.Next = link(params, limit, skip+limit)
	}
	if skip > 0 {
		previous := skip - limit
		if previous < 0 {
			previous = 0
		}
		page.Links.Prev = link(params, limit, previous)
	}
	return page, nil
}

func link(params map[string]string, limit, skip int64) string {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}

	values.Set("limit", strconv.FormatInt(limit, 10))
	values.Set("skip", strconv.FormatInt(skip, 10))
	return gamesPath + "?" + values.Encode()
}
