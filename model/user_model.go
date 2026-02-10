package model

import "time"

type User struct {
	Id        int64      `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Status    int        `json:"status"`
	LastLogin *time.Time `json:"last_login_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type UserProfile struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateUserpayload struct {
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
	FullName             string `json:"full_name" binding:"required"`
	PhoneNumber          string `json:"phone_number" binding:"required"`
}
