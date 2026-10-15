package repository

import (
	"fmt"
	"net/url"

	"go.mongodb.org/mongo-driver/v2/bson"

	"rest-api/model"
	"rest-api/query"
)

func buildQuery(version model.FilterVersion, params map[string]string) (filter, sort bson.D, err error) {
	switch version {
	case model.FilterV1:
		if filter, err = buildFilter(params); err != nil {
			return nil, nil, err
		}
		sort, err = buildSort(params)

	case model.FilterV2:
		q, parseErr := ParseGameQuery(params)
		if parseErr != nil {
			return nil, nil, parseErr
		}
		filter, sort = q.Filter(), q.Sort()

	case model.FilterV3:
		values := make(url.Values, len(params))
		for key, value := range params {
			values.Set(key, value)
		}
		q, parseErr := query.ParseGames(values)
		if parseErr != nil {
			return nil, nil, parseErr
		}
		filter, sort = gamesFilter(q), gamesSort(q)

	default:
		return nil, nil, fmt.Errorf("unknown filter version %d", version)
	}
	return filter, sort, err
}
