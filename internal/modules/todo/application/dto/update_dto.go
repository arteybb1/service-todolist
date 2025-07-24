package dto

type UpdateTodoStatusDTO struct {
	Status bool `json:"status" bson:"status"`
}
