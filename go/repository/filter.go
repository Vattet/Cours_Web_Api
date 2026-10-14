package repository

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"

	"rest-api/model"
)

type rule struct {
	field string
	build func(string) (any, error)
}

var rules = map[string]rule{
	"title":        {"title", containsIgnoreCase},
	"release_year": {"release_year", exactInt},
	"platforms":    {"platforms", allOfLower},
	"genre":        {"genres.description", anyOfName},
	"genre_id":     {"genres.id", anyOfInt},
	"is_free":      {"is_free", exactBool},
	"developers":   {"developers", allOf},
	"publishers":   {"publishers", allOf},
}

var sortable = map[string]string{
	"title":        "title",
	"release_year": "release_year",
	"price":        "price",
}

func anyOfInt(value string) (any, error) {
	var ids []int

	for _, raw := range split(value, false) {
		number, err := strconv.Atoi(raw)
		if err != nil {
			return nil, model.ErrInvalidFilter
		}
		ids = append(ids, number)
	}
	if len(ids) == 0 {
		return nil, model.ErrInvalidFilter
	}
	return bson.D{{Key: "$in", Value: ids}}, nil
}

func containsIgnoreCase(value string) (any, error) {
	return bson.D{
		{Key: "$regex", Value: regexp.QuoteMeta(value)},
		{Key: "$options", Value: "i"},
	}, nil
}

func exactInt(value string) (any, error) {
	number, err := strconv.Atoi(value)
	if err != nil {
		return nil, model.ErrInvalidFilter
	}
	return number, nil
}

func allOfLower(value string) (any, error) {
	return allValues(value, true)
}

func anyOfName(value string) (any, error) {
	names := split(value, false)
	if len(names) == 0 {
		return nil, model.ErrInvalidFilter
	}
	patterns := make([]bson.Regex, len(names))
	for i, name := range names {
		patterns[i] = bson.Regex{
			Pattern: "^" + regexp.QuoteMeta(name) + "$",
			Options: "i",
		}
	}
	return bson.D{{Key: "$in", Value: patterns}}, nil
}

func allValues(value string, lower bool) (any, error) {
	items := split(value, lower)
	if len(items) == 0 {
		return nil, model.ErrInvalidFilter
	}
	return bson.D{{Key: "$all", Value: items}}, nil
}

func split(value string, lower bool) []string {
	var out []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if lower {
			item = strings.ToLower(item)
		}
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func exactBool(value string) (any, error) {
	boolean, err := strconv.ParseBool(value)

	if err != nil {
		return nil, model.ErrInvalidFilter
	}
	return boolean, err
}
func allOf(value string) (any, error) {
	return allValues(value, false)
}
func buildFilter(params map[string]string) (bson.D, error) {
	filter := bson.D{}
	for param, value := range params {
		r, ok := rules[param]
		if !ok {
			continue
		}
		if strings.TrimSpace(value) == "" {
			continue
		}
		condition, err := r.build(value)

		if err != nil {
			return nil, err
		}
		filter = append(filter, bson.E{Key: r.field, Value: condition})
	}
	return filter, nil
}

func buildSort(params map[string]string) (bson.D, error) {
	ordering := strings.TrimSpace(params["ordering"])

	if ordering == "" {
		return bson.D{{Key: "_id", Value: 1}}, nil
	}

	direction := 1
	if strings.HasPrefix(ordering, "-") {
		direction = -1
		ordering = strings.TrimPrefix(ordering, "-")
	}

	field, ok := sortable[ordering]
	if !ok {
		return nil, fmt.Errorf("%w: cannot sort by %q", model.ErrInvalidFilter, ordering)
	}

	return bson.D{
		{Key: field, Value: direction},
		{Key: "_id", Value: 1},
	}, nil
}
