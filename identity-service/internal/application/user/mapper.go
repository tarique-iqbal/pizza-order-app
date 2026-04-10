package user

import "identity-service/internal/domain/user"

func MapToResponse(u *user.User) Response {
	return Response{
		ID:     u.ID,
		Email:  u.Email,
		Role:   string(u.Role),
		Status: string(u.Status),
		Name: UserName{
			First: u.FirstName,
			Last:  u.LastName,
		},
		Metadata: UserMetadata{
			LastLogin:   u.LoggedAt,
			MemberSince: u.CreatedAt,
			LastUpdate:  u.UpdatedAt,
		},
	}
}
