package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type TemplateData struct {
	Title string
}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static files
	fileServer := http.FileServer(http.Dir("static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Routes
	r.Get("/", homeHandler)

	// Start server
	port := "8081"
	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Chewawi",
	}

	tmpl, err := template.ParseFiles(
		"src/views/layout.html",
		"src/views/components/hero.html",
		"src/views/components/footer.html",
		"src/views/components/section.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
