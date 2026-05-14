package auth

type InvalidRefreshTokenError struct{}

func (InvalidRefreshTokenError) Error() string {
	return "invalid refresh token"
}

type RefreshTokenExpiredError struct{}

func (RefreshTokenExpiredError) Error() string {
	return "refresh token expired"
}

type UserNotFoundError struct{}

func (UserNotFoundError) Error() string {
	return "user not found"
}

var (
	ErrInvalidRefreshToken = InvalidRefreshTokenError{}
	ErrRefreshTokenExpired = RefreshTokenExpiredError{}
	ErrUserNotFound        = UserNotFoundError{}
)
