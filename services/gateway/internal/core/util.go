package core

import "net/http"

// GetUsernameFromRequest retrieves the username from the current request's context
func GetUsernameFromRequest(r *http.Request) string {
	return r.Context().Value("username").(string)
}
