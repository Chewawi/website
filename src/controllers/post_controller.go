package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"chewawi_web/src/models"
	"chewawi_web/src/utils"

	"github.com/go-chi/chi/v5"
)

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Title       string
	Posts       []models.Post
	Post        models.Post
	HTMLContent template.HTML
	Error       string
	IsAdmin     bool
	Content     template.HTML
}

// ListPostsHandler handles the GET /posts route
func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Get all posts
	posts, err := models.GetAllPosts()
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := TemplateData{
		Title: "Blog Posts",
		Posts: posts,
	}

	// First, render the content template
	contentTmpl, err := template.ParseFiles("src/views/posts/list.html")
	if err != nil {
		log.Printf("Content template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute content template to a buffer
	var contentBuffer strings.Builder
	err = contentTmpl.ExecuteTemplate(&contentBuffer, "post-list", data)
	if err != nil {
		log.Printf("Content template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Add the rendered content to the data
	data.Content = template.HTML(contentBuffer.String())

	// Parse layout template
	layoutTmpl, err := template.ParseFiles(
		"src/views/layout.html",
		"src/views/home/hero.html",
		"src/views/home/footer.html",
		"src/views/home/section.html",
	)
	if err != nil {
		log.Printf("Layout template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute layout template
	err = layoutTmpl.Execute(w, data)
	if err != nil {
		log.Printf("Layout template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ViewPostHandler handles the GET /posts/:slug route
func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get slug from URL
	slug := chi.URLParam(r, "slug")

	// Get post by slug
	post, err := models.GetPostBySlug(slug)
	if err != nil {
		log.Printf("Error getting post: %v", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Convert Markdown to HTML
	htmlContent := template.HTML(utils.MarkdownToHTML(post.Content))

	// Prepare template data
	data := TemplateData{
		Title:       post.Title,
		Post:        post,
		HTMLContent: htmlContent,
	}

	// First, render the content template
	contentTmpl, err := template.ParseFiles("src/views/posts/single.html")
	if err != nil {
		log.Printf("Content template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute content template to a buffer
	var contentBuffer strings.Builder
	err = contentTmpl.ExecuteTemplate(&contentBuffer, "single-post", data)
	if err != nil {
		log.Printf("Content template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Add the rendered content to the data
	data.Content = template.HTML(contentBuffer.String())

	// Parse layout template
	layoutTmpl, err := template.ParseFiles(
		"src/views/layout.html",
		"src/views/home/hero.html",
		"src/views/home/footer.html",
		"src/views/home/section.html",
	)
	if err != nil {
		log.Printf("Layout template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute layout template
	err = layoutTmpl.Execute(w, data)
	if err != nil {
		log.Printf("Layout template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// HomeHandler handles the GET / route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get recent posts (limit to 3)
	posts, err := models.GetAllPosts()
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		// Continue without posts
		posts = []models.Post{}
	}

	// Limit to 3 posts for homepage
	if len(posts) > 3 {
		posts = posts[:3]
	}

	// Check if user is admin
	_, isAdmin := r.Context().Value("username").(string)

	// Prepare template data
	data := TemplateData{
		Title:   "Chewawi",
		Posts:   posts,
		IsAdmin: isAdmin,
	}

	// For the home page, we don't pre-render a content template
	// The layout will use the default "section" template

	// Parse layout template
	layoutTmpl, err := template.ParseFiles(
		"src/views/layout.html",
		"src/views/home/hero.html",
		"src/views/home/footer.html",
		"src/views/home/section.html",
	)
	if err != nil {
		log.Printf("Layout template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute layout template
	err = layoutTmpl.Execute(w, data)
	if err != nil {
		log.Printf("Layout template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
