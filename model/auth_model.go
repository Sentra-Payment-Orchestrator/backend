package model

// Define the structure for registration request payload make it required in auth handler
type RegisterRequest struct {
	FirstName            string `json:"first_name" binding:"required"`
	LastName             string `json:"last_name" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
	Email                string `json:"email" binding:"required,email"`
	PhoneNumber          string `json:"phone_number" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
