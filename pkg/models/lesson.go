package models

type Lesson struct {
	ID            string `json:"id" bson:"_id"`
	Title         string `json:"title" bson:"title"`
	Lection       string `json:"lection" bson:"lection"`
	Task          string `json:"task" bson:"task"`
	EstimatedTime int    `json:"time" bson:"time"`
}
