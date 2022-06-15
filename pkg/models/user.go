package models

type User struct {
	ID       string   `json:"id" bson:"_id"`
	ChatID   int64    `json:"chat_id" bson:"chat_id"`
	Email    string   `json:"email" bson:"email"`
	Password string   `json:"password" bson:"password"`
	Courses  []Course `json:"courses" bson:"courses"`
}
