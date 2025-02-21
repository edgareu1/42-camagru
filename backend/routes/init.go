package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"backend/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !utils.IsRequestAuthenticated(r) {
			utils.SendError(w, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func InitRoutes(router *mux.Router) {
	routerWithAuth := router.NewRoute().Subrouter()
	routerWithAuth.Use(AuthMiddleware)

	// Users
	routerWithAuth.HandleFunc("/users", GetUsers).Methods("GET")
	routerWithAuth.HandleFunc("/users/{id}", GetUser).Methods("GET")
	routerWithAuth.HandleFunc("/users/{id}", EditUser).Methods("PATCH")
	router.HandleFunc("/users/{id}/verify", VerifyUser).Methods("POST")
	router.HandleFunc("/users/{id}/renew-password", RenewPassword).Methods("POST")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/sign-in", SignInUser).Methods("POST")
	router.HandleFunc("/users/auth", AuthUser).Methods("POST")
	router.HandleFunc("/users/forgot-password", RequestPasswordChange).Methods("POST")

	// Images
	routerWithAuth.HandleFunc("/users/{id}/images", CreateImage).Methods("POST")
	router.HandleFunc("/images", GetImages).Methods("GET")
	router.HandleFunc("/users/{id}/images", GetUserImages).Methods("GET")
	router.HandleFunc("/images/{id}", GetImage).Methods("GET")
	router.HandleFunc("/images/{id}/details", GetImageDetails).Methods("GET")
	routerWithAuth.HandleFunc("/images/{id}", DeleteImage).Methods("DELETE")

	// Likes
	routerWithAuth.HandleFunc("/images/{id}/like", LikeImage).Methods("POST")
	routerWithAuth.HandleFunc("/images/{id}/unlike", UnlikeImage).Methods("POST")

	// Comments
	router.HandleFunc("/images/{id}/comments", GetComments).Methods("GET")
	routerWithAuth.HandleFunc("/images/{id}/comments", CreateComment).Methods("POST")
}
