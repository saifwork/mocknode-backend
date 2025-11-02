package dtos

// UserRegisterRequestDto represents the payload for user registration
type UserRegisterRequestDto struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserLoginRequestDto represents the payload for user login
type UserLoginRequestDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponseDto is the safe user data returned to the client
type UserResponseDto struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

// AuthResponseDto represents the response returned after login or registration
type AuthResponseDto struct {
	User  UserResponseDto `json:"user"`
	Token string          `json:"token"`
}

// UserUpdateRequestDto handles update profile requests
type UserUpdateRequestDto struct {
	FullName string `json:"fullName,omitempty"`
	Password string `json:"password,omitempty"`
}

// UserChangePasswordDto handles password change requests
type UserChangePasswordDto struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// UserTokenClaims defines JWT claims or session data
type UserTokenClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}
