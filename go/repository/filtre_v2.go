package repository

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/v2/bson"

	"rest-api/model"
)

type GameQuery struct {
	Title       string `schema:"title"        validate:"max=100"`
	ReleaseYear *int   `schema:"release_year" validate:"omitempty,gte=1950,lte=2100"`
	IsFree      *bool  `schema:"is_free"` // 0/1, true/false
	Platforms   CSV    `schema:"platforms"`
	Developers  CSV    `schema:"developers"`
	Publishers  CSV    `schema:"publishers"`
	Genres      CSV    `schema:"genre"`
	GenreIDs    IntCSV `schema:"genre_id"     validate:"dive,gt=0"`
	Order       string `schema:"order"     validate:"omitempty,oneof=title -title release_year -release_year price -price"`
}

var (
	decoder  = newDecoder()
	validate = newValidator()
)

func newDecoder() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return d
}

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		return f.Tag.Get("schema")
	})
	return v
}

func ParseGameQuery(params map[string]string) (GameQuery, error) {
	values := make(map[string][]string, len(params))
	for key, value := range params {
		if value = strings.TrimSpace(value); value != "" {
			values[key] = []string{value}
		}
	}

	var q GameQuery
	if err := decoder.Decode(&q, values); err != nil {
		return GameQuery{}, fmt.Errorf("%w: %s", model.ErrInvalidFilter, err)
	}
	if err := validate.Struct(q); err != nil {
		return GameQuery{}, fmt.Errorf("%w: %s", model.ErrInvalidFilter, describe(err))
	}
	return q, nil
}

func (q GameQuery) Filter() bson.D {
	filter := bson.D{}
	add := func(field string, value any) {
		filter = append(filter, bson.E{Key: field, Value: value})
	}

	if q.Title != "" {
		add("title", bson.Regex{Pattern: regexp.QuoteMeta(q.Title), Options: "i"})
	}
	if q.ReleaseYear != nil {
		add("release_year", *q.ReleaseYear)
	}
	if q.IsFree != nil {
		add("is_free", *q.IsFree)
	}
	if len(q.Platforms) > 0 {
		add("platforms", bson.M{"$all": q.Platforms.Lower()})
	}
	if len(q.Developers) > 0 {
		add("developers", bson.M{"$all": q.Developers})
	}
	if len(q.Publishers) > 0 {
		add("publishers", bson.M{"$all": q.Publishers})
	}
	if len(q.Genres) > 0 {
		add("genres.description", bson.M{"$in": q.Genres.ExactIgnoreCase()})
	}
	if len(q.GenreIDs) > 0 {
		add("genres.id", bson.M{"$in": q.GenreIDs})
	}
	return filter
}

func (q GameQuery) Sort() bson.D {
	field, desc := strings.CutPrefix(q.Order, "-")
	if field == "" {
		return bson.D{{Key: "_id", Value: 1}}
	}

	direction := 1
	if desc {
		direction = -1
	}
	return bson.D{{Key: field, Value: direction}, {Key: "_id", Value: 1}}
}

func buildFilterV2(params map[string]string) (bson.D, error) {
	q, err := ParseGameQuery(params)
	if err != nil {
		return nil, err
	}
	return q.Filter(), nil
}

func buildSortV2(params map[string]string) (bson.D, error) {
	q, err := ParseGameQuery(params)
	if err != nil {
		return nil, err
	}
	return q.Sort(), nil
}

func describe(err error) string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err.Error()
	}

	msgs := make([]string, len(verrs))
	for i, e := range verrs {
		msgs[i] = fmt.Sprintf("%s: %s=%s", e.Field(), e.Tag(), e.Param())
	}
	return strings.Join(msgs, ", ")
}

type CSV []string

func (c *CSV) UnmarshalText(text []byte) error {
	*c = nil
	for item := range strings.SplitSeq(string(text), ",") {
		if item = strings.TrimSpace(item); item != "" {
			*c = append(*c, item)
		}
	}
	return nil
}

func (c CSV) Lower() []string {
	out := make([]string, len(c))
	for i, item := range c {
		out[i] = strings.ToLower(item)
	}
	return out
}

func (c CSV) ExactIgnoreCase() []bson.Regex {
	out := make([]bson.Regex, len(c))
	for i, item := range c {
		out[i] = bson.Regex{Pattern: "^" + regexp.QuoteMeta(item) + "$", Options: "i"}
	}
	return out
}

type IntCSV []int

func (c *IntCSV) UnmarshalText(text []byte) error {
	var items CSV
	_ = items.UnmarshalText(text)

	*c = make(IntCSV, 0, len(items))
	for _, item := range items {
		n, err := strconv.Atoi(item)
		if err != nil {
			return fmt.Errorf("%q n'est pas un entier", item)
		}
		*c = append(*c, n)
	}
	return nil
}
