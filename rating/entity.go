package rating

type AverageInput struct {
	ItemId  string
	StartId string
}

type AverageScore struct {
	ScorePoint int
	StarId     string
	Star       int
}

type Average struct {
	RatingId     string
	ItemId       string
	Avg          float64
	AverageScore []AverageScore
}
