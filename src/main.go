package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"chewawi_web/src/controllers"
	"chewawi_web/src/database"
	"chewawi_web/src/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	database.InitDB()
	defer database.CloseDB()

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(addUserToContext)

	fileServer := http.FileServer(http.Dir("static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Public routes
	r.Get("/", controllers.HomeHandler)
	r.Get("/posts", controllers.ListPostsHandler)
	r.Get("/posts/{slug}", controllers.ViewPostHandler)

	// Authentication routes
	r.Get("/login", controllers.LoginHandler)
	r.Post("/login", controllers.LoginSubmitHandler)
	r.Post("/logout", controllers.LogoutHandler)

	// Admin routes (protected)
	r.Route("/owner", func(r chi.Router) {
		// Use auth middleware for all /owner routes
		r.Use(middleware.AuthMiddleware)

		r.Get("/", controllers.DashboardHandler)
		r.Get("/new", controllers.NewPostHandler)
		r.Post("/new", controllers.CreatePostHandler)
		r.Get("/edit/{slug}", controllers.EditPostHandler)
		r.Post("/edit/{slug}", controllers.UpdatePostHandler)
		r.Post("/delete/{slug}", controllers.DeletePostHandler)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// addUserToContext middleware adds the username to the request context if the user is authenticated
func addUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == nil {

			claims := &middleware.Claims{}
			token, err := middleware.ParseToken(cookie.Value, claims)

			if err == nil && token.Valid {
				ctx := context.WithValue(r.Context(), "username", claims.Username)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
