package mapper

import (
	"api-service/internal/domain/auth"
	"net/http"
)

func MapErrorToHTTPStatus(err error) int {
	switch err {
	case auth.ErrCodeInvalid,
		auth.ErrCodeNotIssued:
		return http.StatusBadRequest
	case auth.ErrCodeExpired:
		return http.StatusGone
	case auth.ErrCodeUsed:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
