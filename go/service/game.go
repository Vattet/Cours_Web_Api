package service

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net/url"
	"rest-api/model"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
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

type GameWriter interface {
	Create(ctx context.Context, game model.Game) (model.Game, error)
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

var validate = validator.New(validator.WithRequiredStructEnabled())

func CreateGame(ctx context.Context, w GameWriter, in model.CreateGame) (model.Game, error) {
	if err := validate.Struct(in); err != nil {
		return model.Game{}, fmt.Errorf("%w: %s", model.ErrInvalidGame, err)
	}

	if in.IsFree && in.Price != nil && *in.Price > 0 {
		return model.Game{}, fmt.Errorf("%w: a free game cannot have a price", model.ErrInvalidGame)
	}
	if !in.IsFree && in.Price == nil {
		return model.Game{}, fmt.Errorf("%w: price is required for a paid game", model.ErrInvalidGame)
	}

	game := model.Game{
		Title:       strings.TrimSpace(in.Title),
		Description: in.Description,
		IsFree:      in.IsFree,
		Price:       in.Price,
		RequiredAge: in.RequiredAge,
		Image:       in.Image,
		Website:     in.Website,
		Developers:  in.Developers,
		Publishers:  in.Publishers,
		Platforms:   in.Platforms,
		Genres:      in.Genres,
		ReleaseDate: in.ReleaseDate,
	}
	if in.ReleaseDate != nil {
		year, _ := strconv.Atoi((*in.ReleaseDate)[:4]) // format déjà validé : AAAA-MM-JJ
		game.ReleaseYear = &year
	}

	return w.Create(ctx, game)
}
