package api

import "net/http"

func permissionMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Tambahkan RBAC project di sini, misalnya:
			// access, err := permissionRepo.GetUserAccessByID(r.Context(), userID)
			next.ServeHTTP(w, r)
		})
	}
}
