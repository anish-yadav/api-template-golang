package constants

const (
	IntervalServerError  = "internal sever error"
	TokenExpired         = "jwt token expired"
	InvalidSigningMethod = "invalid signing method, use hmac"
	InvalidToken         = "invalid token"

	StudentNotFound = "student data not present"

	ItemNotFound = "item not found"

	BadRequest       = "malformed request data"
	Conflict         = "data already present"
	PasswordMismatch = "email or password invalid"
	Forbidden        = "user is not authorized to access this resource "
)
