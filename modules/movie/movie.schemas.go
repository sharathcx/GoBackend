package movie

type GenreSchema struct {
	GenreId   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=100"`
}

type RankingSchema struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"required"`
}

type MovieSchema struct {
	MovieID     string        `bson:"movie_id" json:"movie_id"`
	Title       string        `bson:"title" json:"title"`
	PosterPath  string        `bson:"poster_path" json:"poster_path"`
	YoutubeID   string        `bson:"youtube_id" json:"youtube_id"`
	Genre       []GenreSchema `bson:"genre" json:"genre"`
	AdminReview string        `bson:"admin_review" json:"admin_review"`
	Ranking     RankingSchema `bson:"ranking" json:"ranking"`
}

type GetMoviesPayloadSchema struct {
	Limit int `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
}

type AddMoviePayloadSchema struct {
	Title       string        `json:"title" binding:"required,min=2,max=100"`
	PosterPath  string        `json:"poster_path" binding:"required,url"`
	YoutubeID   string        `json:"youtube_id" binding:"required"`
	Genre       []GenreSchema `json:"genre" binding:"required,dive"`
	AdminReview string        `json:"admin_review" binding:"required"`
	Ranking     RankingSchema `json:"ranking" binding:"required"`
}
