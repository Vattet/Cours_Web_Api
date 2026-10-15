package query

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"

	"rest-api/model"
)

type Games struct {
	Title       string   `schema:"title"        validate:"max=100"`
	ReleaseYear *int     `schema:"release_year" validate:"omitempty,gte=1950,lte=2100"`
	IsFree      *bool    `schema:"is_free"`
	Platforms   []string `schema:"platforms"    validate:"max=10,dive,max=50"`
	Developers  []string `schema:"developers"   validate:"max=10,dive,max=100"`
	Publishers  []string `schema:"publishers"   validate:"max=10,dive,max=100"`
	Genres      []string `schema:"genre"        validate:"max=10,dive,max=50"`
	GenreIDs    []int    `schema:"genre_id"     validate:"max=10,dive,gt=0"`
	Order       []string `schema:"order"        validate:"max=3,dive,oneof=title -title release_year -release_year price -price"`
}

type SortField struct {
	Name string
	Desc bool
}

func (g Games) Sort() []SortField {
	fields := make([]SortField, len(g.Order))
	for i, raw := range g.Order {
		name, desc := strings.CutPrefix(raw, "-")
		fields[i] = SortField{Name: name, Desc: desc}
	}
	return fields
}

func ParseGames(values url.Values) (Games, error) {
	var g Games
	if err := decoder.Decode(&g, normalize(values)); err != nil {
		return Games{}, fromDecodeError(err)
	}
	if err := validate.Struct(g); err != nil {
		return Games{}, fromValidationError(err)
	}
	if err := checkDuplicateOrder(g.Order); err != nil {
		return Games{}, err
	}
	return g, nil
}

type FieldError struct {
	Param   string `json:"param"`
	Message string `json:"message"`
}

type Error struct {
	Fields []FieldError `json:"errors"`
}

func (e *Error) Error() string {
	parts := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		parts[i] = f.Param + ": " + f.Message
	}
	return model.ErrInvalidFilter.Error() + ": " + strings.Join(parts, "; ")
}

func (e *Error) Unwrap() error { return model.ErrInvalidFilter }

func fromDecodeError(err error) error {
	var multi schema.MultiError
	if !errors.As(err, &multi) {
		return &Error{Fields: []FieldError{{Param: "query", Message: err.Error()}}}
	}

	params := make([]string, 0, len(multi))
	for param := range multi {
		params = append(params, param)
	}
	sort.Strings(params)

	fields := make([]FieldError, len(params))
	for i, param := range params {
		fields[i] = FieldError{Param: baseParam(param), Message: "invalid value type"}
	}
	return &Error{Fields: fields}
}

func fromValidationError(err error) error {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return &Error{Fields: []FieldError{{Param: "query", Message: err.Error()}}}
	}

	fields := make([]FieldError, len(verrs))
	for i, e := range verrs {
		fields[i] = FieldError{Param: baseParam(e.Field()), Message: message(e)}
	}
	return &Error{Fields: fields}
}

func message(e validator.FieldError) string {
	switch e.Tag() {
	case "max":
		if e.Kind() == reflect.Slice {
			return "at most " + e.Param() + " values allowed"
		}
		return "must be at most " + e.Param() + " characters"
	case "gt":
		return "must be greater than " + e.Param()
	case "gte":
		return "must be greater than or equal to " + e.Param()
	case "lte":
		return "must be less than or equal to " + e.Param()
	case "oneof":
		return "must be one of: " + strings.ReplaceAll(e.Param(), " ", ", ")
	default:
		return "failed rule " + e.Tag()
	}
}

func checkDuplicateOrder(order []string) error {
	seen := make(map[string]bool, len(order))
	for _, raw := range order {
		name := strings.TrimPrefix(raw, "-")
		if seen[name] {
			return &Error{Fields: []FieldError{{Param: "order", Message: fmt.Sprintf("%q is used more than once", name)}}}
		}
		seen[name] = true
	}
	return nil
}

func baseParam(name string) string {
	if i := strings.IndexByte(name, '['); i >= 0 {
		return name[:i]
	}
	return name
}

var (
	decoder    = newDecoder()
	validate   = newValidator()
	listParams = paramsOfKind[Games](reflect.Slice)
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

func normalize(values url.Values) url.Values {
	out := make(url.Values, len(values))
	for key, raw := range values {
		for _, value := range raw {
			parts := []string{value}
			if listParams[key] {
				parts = strings.Split(value, ",")
			}
			for _, part := range parts {
				if part = strings.TrimSpace(part); part != "" {
					out[key] = append(out[key], part)
				}
			}
		}
	}
	return out
}

func paramsOfKind[T any](kind reflect.Kind) map[string]bool {
	params := map[string]bool{}
	t := reflect.TypeFor[T]()
	for i := range t.NumField() {
		if f := t.Field(i); f.Type.Kind() == kind {
			params[f.Tag.Get("schema")] = true
		}
	}
	return params
}
