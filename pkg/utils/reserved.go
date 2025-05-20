package utils

import "strings"

var ReservedCodes = []string{
	"api", "api/", "metrics", "health",
	"register", "login", "me", "shorten",
	"my", "admin", "verify", "docs", "assets",
}

func IsReservedCode(code string) bool {
	code = strings.ToLower(code)

	for _, reserved := range ReservedCodes {
		if code == reserved || strings.HasPrefix(code, reserved+"/") {
			return true
		}
	}
	return false
}
