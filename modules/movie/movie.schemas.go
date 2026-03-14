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
	MovieID     string        `bson:"movie_id" json:"movie_id" validate:"required"`
	Title       string        `bson:"title" json:"title" validate:"required,min=2,max=100"`
	PosterPath  string        `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YoutubeID   string        `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre       []GenreSchema `bson:"genre" json:"genre" validate:"required,dive"`
	AdminReview string        `bson:"admin_review" json:"admin_review" validate:"required"`
	Ranking     RankingSchema `bson:"ranking" json:"ranking" validate:"required"`
}

type GetMoviesPayloadSchema struct {
	Limit int `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
}

type GetMoviePayloadSchema struct {
	MovieID string `uri:"movie_id" binding:"required" json:"movie_id"`
}

type AddMoviePayloadSchema struct {
	Title       string        `bson:"title" json:"title" binding:"required,min=2,max=100"`
	PosterPath  string        `bson:"poster_path" json:"poster_path" binding:"required,url"`
	YoutubeID   string        `bson:"youtube_id" json:"youtube_id" binding:"required"`
	Genre       []GenreSchema `bson:"genre" json:"genre" binding:"required,dive"`
	AdminReview string        `bson:"admin_review" json:"admin_review" binding:"required"`
	Ranking     RankingSchema `bson:"ranking" json:"ranking" binding:"required"`
}
