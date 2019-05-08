package core

import "net/http"

func GetUsernameFromRequest(r *http.Request) string {
	return r.Context().Value("username").(string)
}
