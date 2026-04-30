package user

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required,oneof=user owner admin"`
	Code      string `json:"code" binding:"required,len=6,numeric"`
}

type RegisterOwnerRequest struct {
	RegisterRequest
	BusinessName string `json:"businessName" binding:"required,min=2,max=128"`
	VATNumber    string `json:"vatNumber" binding:"required,startswith=DE,len=11,alphanum"`
}

type Response struct {
	ID       uuid.UUID    `json:"id"`
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
