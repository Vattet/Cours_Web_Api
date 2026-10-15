package repository

import (
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"

	"rest-api/query"
)

var gameSortFields = map[string]string{
	"title":        "title",
	"release_year": "release_year",
	"price":        "price",
}

func gamesFilter(q query.Games) bson.D {
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
		add("platforms", bson.M{"$all": lowerAll(q.Platforms)})
	}
	if len(q.Developers) > 0 {
		add("developers", bson.M{"$all": q.Developers})
	}
	if len(q.Publishers) > 0 {
		add("publishers", bson.M{"$all": q.Publishers})
	}
	if len(q.Genres) > 0 {
		add("genres.description", bson.M{"$in": exactIgnoreCase(q.Genres)})
	}
	if len(q.GenreIDs) > 0 {
		add("genres.id", bson.M{"$in": q.GenreIDs})
	}
	return filter
}

func gamesSort(q query.Games) bson.D {
	sort := bson.D{}
	for _, f := range q.Sort() {
		direction := 1
		if f.Desc {
			direction = -1
		}
		sort = append(sort, bson.E{Key: gameSortFields[f.Name], Value: direction})
	}
	return append(sort, bson.E{Key: "_id", Value: 1})
}

func lowerAll(items []string) []string {
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = strings.ToLower(item)
	}
	return out
}

func exactIgnoreCase(items []string) []bson.Regex {
	out := make([]bson.Regex, len(items))
	for i, item := range items {
		out[i] = bson.Regex{Pattern: "^" + regexp.QuoteMeta(item) + "$", Options: "i"}
	}
	return out
}
