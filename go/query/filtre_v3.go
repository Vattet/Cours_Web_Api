package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"rest-api/model"
)

// Les filtres possibles pour GET /api/v1/games
type Games struct {
	Title       string
	ReleaseYear *int  // nil = pas de filtre
	IsFree      *bool // nil = pas de filtre
	Platforms   []string
	Developers  []string
	Publishers  []string
	Genres      []string
	GenreIDs    []int
	Order       []string // ex: "title" ou "-release_year" (le - veut dire décroissant)
}

type SortField struct {
	Name string
	Desc bool
}

var allowedOrder = []string{"title", "release_year", "price"}

func (g Games) Sort() []SortField {
	var fields []SortField
	for _, o := range g.Order {
		if strings.HasPrefix(o, "-") {
			fields = append(fields, SortField{Name: o[1:], Desc: true})
		} else {
			fields = append(fields, SortField{Name: o, Desc: false})
		}
	}
	return fields
}

// ParseGames lit les paramètres de l'URL et vérifie qu'ils sont valides.
func ParseGames(values url.Values) (Games, error) {
	var g Games
	var err error

	g.Title = strings.TrimSpace(values.Get("title"))
	if len(g.Title) > 100 {
		return g, invalid("title", "100 caractères maximum")
	}

	if v := values.Get("release_year"); v != "" {
		year, err := strconv.Atoi(v)
		if err != nil || year < 1950 || year > 2100 {
			return g, invalid("release_year", "doit être une année entre 1950 et 2100")
		}
		g.ReleaseYear = &year
	}

	if v := values.Get("is_free"); v != "" {
		free, err := strconv.ParseBool(v)
		if err != nil {
			return g, invalid("is_free", "doit être true ou false")
		}
		g.IsFree = &free
	}

	if g.Platforms, err = parseList(values, "platforms", 50); err != nil {
		return g, err
	}
	if g.Developers, err = parseList(values, "developers", 100); err != nil {
		return g, err
	}
	if g.Publishers, err = parseList(values, "publishers", 100); err != nil {
		return g, err
	}
	if g.Genres, err = parseList(values, "genre", 50); err != nil {
		return g, err
	}

	ids, err := parseList(values, "genre_id", 20)
	if err != nil {
		return g, err
	}
	for _, s := range ids {
		id, err := strconv.Atoi(s)
		if err != nil || id <= 0 {
			return g, invalid("genre_id", "doit être un entier positif")
		}
		g.GenreIDs = append(g.GenreIDs, id)
	}

	if g.Order, err = parseOrder(values); err != nil {
		return g, err
	}

	return g, nil
}

// parseList transforme "a,b,c" en ["a", "b", "c"] (10 valeurs max).
func parseList(values url.Values, key string, maxLen int) ([]string, error) {
	var list []string
	for _, raw := range values[key] {
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if len(part) > maxLen {
				return nil, invalid(key, fmt.Sprintf("chaque valeur fait %d caractères maximum", maxLen))
			}
			list = append(list, part)
		}
	}
	if len(list) > 10 {
		return nil, invalid(key, "10 valeurs maximum")
	}
	return list, nil
}

// parseOrder vérifie le tri : 3 champs max, connus, et pas deux fois le même.
func parseOrder(values url.Values) ([]string, error) {
	order, err := parseList(values, "order", 20)
	if err != nil {
		return nil, err
	}
	if len(order) > 3 {
		return nil, invalid("order", "3 champs maximum")
	}

	used := map[string]bool{}
	for _, o := range order {
		name := strings.TrimPrefix(o, "-")
		if !contains(allowedOrder, name) {
			return nil, invalid("order", "valeurs possibles : title, release_year, price (avec - devant pour décroissant)")
		}
		if used[name] {
			return nil, invalid("order", name+" est utilisé deux fois")
		}
		used[name] = true
	}
	return order, nil
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// invalid renvoie une erreur qui contient ErrInvalidFilter,
// comme ça le handler peut la reconnaître avec errors.Is.
func invalid(param, message string) error {
	return fmt.Errorf("%w: %s: %s", model.ErrInvalidFilter, param, message)
}
