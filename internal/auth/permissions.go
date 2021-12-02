package auth

import "net/http"

var PermissionMap = map[string]map[string]string{
	http.MethodGet: {
		"/health":   "",
		"/users":    "",
		"/users/me": "",
	},
	http.MethodPost: {
		"/users":                 "",
		"/users/change-password": "",
		"/users/reset-password":  "",
	},
	http.MethodDelete: {
		"/users": "",
	},
}
