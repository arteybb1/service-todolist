package dto

type TodoCreateDTO struct {
	UserID string `json:"user_id" bson:"user_id"`
	Title  string `json:"title" bson:"title"`
}
