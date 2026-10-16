package model

type User struct {
	JwtToken string `json:"jwt_token,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Theme struct {
	ID          int    `json:"id" bson:"id"`
	Description string `json:"description" bson:"description"`
}

type Game struct {
	ID              int      `json:"_id" bson:"_id"`
	Title           string   `json:"title" bson:"title"`
	RequiredAge     any      `json:"required_age" bson:"required_age"`
	IsFree          bool     `json:"is_free" bson:"is_free"`
	Description     string   `json:"description" bson:"description"`
	Image           string   `json:"image" bson:"image"`
	Website         *string  `json:"website" bson:"website"`
	Developers      []string `json:"developers,omitempty" bson:"developers,omitempty"`
	Publishers      []string `json:"publishers" bson:"publishers"`
	Platforms       []string `json:"platforms" bson:"platforms"`
	Genres          []Theme  `json:"genres,omitempty" bson:"genres,omitempty"`
	Categories      []Theme  `json:"categories,omitempty" bson:"categories,omitempty"`
	Price           *float64 `json:"price,omitempty" bson:"price,omitempty"`
	Recommendations *int     `json:"recommendations,omitempty" bson:"recommendations,omitempty"`
	MetacriticScore *int     `json:"metacritic_score,omitempty" bson:"metacritic_score,omitempty"`
	ReleaseDate     *string  `json:"release_date" bson:"release_date"`
	ReleaseYear     *int     `json:"release_year" bson:"release_year"`
}
type CreateGame struct {
	Title       string   `json:"title"        validate:"required,max=200"`
	Description string   `json:"description"  validate:"max=5000"`
	IsFree      bool     `json:"is_free"`
	Price       *float64 `json:"price"        validate:"omitempty,gte=0"`
	RequiredAge int      `json:"required_age" validate:"gte=0,lte=21"`
	Image       string   `json:"image"        validate:"omitempty,url"`
	Website     *string  `json:"website"      validate:"omitempty,url"`
	Developers  []string `json:"developers"   validate:"max=20,dive,required,max=100"`
	Publishers  []string `json:"publishers"   validate:"max=20,dive,required,max=100"`
	Platforms   []string `json:"platforms"    validate:"required,min=1,dive,oneof=windows mac linux"`
	Genres      []Theme  `json:"genres"       validate:"max=20"`
	ReleaseDate *string  `json:"release_date" validate:"omitempty,datetime=2006-01-02"`
}
