package dto

type CreateUser struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name string `json:"name" example:"John Doe"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" example:"active"`
}
