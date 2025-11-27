package dto

type UserSignupDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRoleDTO struct {
	Role string `json:"role" binding:"required"`
}

type UpdateUserProfileDTO struct {
	Name    *string `json:"name,omitempty"`
	Email   *string `json:"email,omitempty"`
	Address *string `json:"address,omitempty"`
}
