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
