package secret

import "strings"

const MaskedValue = "********"

var sensitiveHeaderNames = map[string]struct{}{
	"authorization":       {},
	"proxy-authorization": {},
	"x-api-key":           {},
	"api-key":             {},
	"x-auth-token":        {},
	"x-access-token":      {},
	"x-secret":            {},
	"cookie":              {},
	"set-cookie":          {},
}

func IsSensitiveHeader(name string) bool {
	_, ok := sensitiveHeaderNames[strings.ToLower(strings.TrimSpace(name))]
	return ok
}

func Mask(value string) string {
	if value == "" {
		return ""
	}
	return MaskedValue
}

func IsMasked(value string) bool {
	return value == MaskedValue
}
