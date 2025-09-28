package common

type Similarity[T any] struct {
	Data  T       `json:"data"`
	Score float64 `json:"score"`
}
