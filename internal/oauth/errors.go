package oauth

import "fmt"

type TokenExchangeError struct {
	StatusCode int
}

func (e TokenExchangeError) Error() string {
	return fmt.Sprintf("failed to exchange code for tokens: status %d", e.StatusCode)
}

type UserInfoError struct {
	StatusCode int
}

func (e UserInfoError) Error() string {
	return fmt.Sprintf("failed to get user info: status %d", e.StatusCode)
}
