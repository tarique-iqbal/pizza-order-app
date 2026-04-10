package user

import "time"

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required,oneof=user owner admin"`
	Code      string `json:"code" binding:"required,min=6"`
}

type Response struct {
	ID       int          `json:"id"`
	Name     UserName     `json:"name"`
	Email    string       `json:"email"`
	Role     string       `json:"role"`
	Status   string       `json:"status"`
	Metadata UserMetadata `json:"metadata"`
}

type UserName struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type UserMetadata struct {
	LastLogin   *time.Time `json:"lastLogin"`
	MemberSince time.Time  `json:"memberSince"`
	LastUpdate  *time.Time `json:"lastUpdate"`
}
