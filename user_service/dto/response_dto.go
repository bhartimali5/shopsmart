package dto

type UserSignupResponseDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UserLoginResponseDTO struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type GetUserProfileResponseDTO struct {
	Email   string `json:"email"`
	Address string `json:"address"`
	Name    string `json:"name"`
	Role    string `json:"role"`
}
