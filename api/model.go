package api

type ErrorResponse struct {
	Message string `json:"message"`
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	LineID   string `json:"line_id" validate:"required,min=1,max=20"`
	Nickname string `json:"nickname" validate:"required,alphanum,min=1,max=20"`
}

type SignUpResponse struct {
	Email    string `json:"email"`
	LineID   string `json:"line_id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}
