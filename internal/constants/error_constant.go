package constants

type ErrorMessage string

const (
	USER_NOT_FOUND        ErrorMessage = "USER_NOT_FOUND"
	INVALID_CREDENTIAL    ErrorMessage = "INVALID_CREDENTIAL"
	INVALID_REFRESH_TOKEN ErrorMessage = "INVALID_REFRESH_TOKEN"
	INVALID_CLAIMS        ErrorMessage = "INVALID_CLAIMS"
	INVALID_ID            ErrorMessage = "INVALID_ID"
	INVALID_ID_FORMAT     ErrorMessage = "INVALID_ID_FORMAT"
	TOKEN_REQUIRED        ErrorMessage = "TOKEN_REQUIRED"
	UNAUTHORIZED          ErrorMessage = "UNAUTHORIZED"
	INTERNAL_ERROR        ErrorMessage = "INTERNAL_ERROR"
)
