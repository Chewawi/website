package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"chewawi_web/src/middleware"
	"chewawi_web/src/models"

	"github.com/go-chi/chi/v5"
)

// LoginHandler handles the GET /login route
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Prepare template data
	data := TemplateData{
		Title: "Login",
	}

	// First, render the content template
	contentTmpl, err := template.ParseFiles("src/views/admin/login.html")
	if err != nil {
		log.Printf("Content template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute content template to a buffer
	var contentBuffer strings.Builder
	err = contentTmpl.ExecuteTemplate(&contentBuffer, "login", data)
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

// LoginSubmitHandler handles the POST /login route
func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Printf("Form parsing error: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get form values
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Authenticate user
	if !middleware.Authenticate(username, password) {
		// Authentication failed
		// Prepare template data
		data := TemplateData{
			Title: "Login",
			Error: "Invalid username or password",
		}

		// First, render the content template
		contentTmpl, err := template.ParseFiles("src/views/admin/login.html")
		if err != nil {
			log.Printf("Content template parsing error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute content template to a buffer
		var contentBuffer strings.Builder
		err = contentTmpl.ExecuteTemplate(&contentBuffer, "login", data)
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
		return
	}

	// Authentication successful, generate token
	token, err := middleware.GenerateToken(username)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Redirect to admin dashboard
	http.Redirect(w, r, "/owner", http.StatusSeeOther)
}

// LogoutHandler handles the POST /logout route
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DashboardHandler handles the GET /owner route
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get all posts
	posts, err := models.GetAllPosts()
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := TemplateData{
		Title:   "Admin Dashboard",
		Posts:   posts,
		IsAdmin: true,
	}

	// First, render the content template
	contentTmpl, err := template.ParseFiles("src/views/admin/dashboard.html")
	if err != nil {
		log.Printf("Content template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute content template to a buffer
	var contentBuffer strings.Builder
	err = contentTmpl.ExecuteTemplate(&contentBuffer, "dashboard", data)
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

// NewPostHandler handles the GET /owner/new route
func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	// Prepare template data
	data := TemplateData{
		Title:   "New Post",
		IsAdmin: true,
	}

	// First, render the content template
	contentTmpl, err := template.ParseFiles("src/views/admin/post_form.html")
	if err != nil {
		log.Printf("Content template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute content template to a buffer
	var contentBuffer strings.Builder
	err = contentTmpl.ExecuteTemplate(&contentBuffer, "post-form", data)
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

// CreatePostHandler handles the POST /owner/new route
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Printf("Form parsing error: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")

	// Validate form
	if title == "" || content == "" {
		// Prepare template data with error
		data := TemplateData{
			Title:   "New Post",
			Error:   "Title and content are required",
			Post:    models.Post{Title: title, Content: content},
			IsAdmin: true,
		}

		// First, render the content template
		contentTmpl, err := template.ParseFiles("src/views/admin/post_form.html")
		if err != nil {
			log.Printf("Content template parsing error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute content template to a buffer
		var contentBuffer strings.Builder
		err = contentTmpl.ExecuteTemplate(&contentBuffer, "post-form", data)
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
		return
	}

	// Create post
	_, err = models.CreatePost(title, content)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to admin dashboard
	http.Redirect(w, r, "/owner", http.StatusSeeOther)
}

// EditPostHandler handles the GET /owner/edit/:slug route
func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get slug from URL
	slug := chi.URLParam(r, "slug")

	// Get post by slug
	post, err := models.GetPostBySlug(slug)
	if err != nil {
		log.Printf("Error getting post: %v", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Prepare template data
	data := TemplateData{
		Title:   "Edit Post",
		Post:    post,
		IsAdmin: true,
	}

	// First, render the content template
	contentTmpl, err := template.ParseFiles("src/views/admin/post_form.html")
	if err != nil {
		log.Printf("Content template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute content template to a buffer
	var contentBuffer strings.Builder
	err = contentTmpl.ExecuteTemplate(&contentBuffer, "post-form", data)
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

// UpdatePostHandler handles the POST /owner/edit/:slug route
func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Get slug from URL
	slug := chi.URLParam(r, "slug")

	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Printf("Form parsing error: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")

	// Validate form
	if title == "" || content == "" {
		// Get original post
		post, err := models.GetPostBySlug(slug)
		if err != nil {
			log.Printf("Error getting post: %v", err)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Update post with form values
		post.Title = title
		post.Content = content

		// Prepare template data with error
		data := TemplateData{
			Title:   "Edit Post",
			Error:   "Title and content are required",
			Post:    post,
			IsAdmin: true,
		}

		// First, render the content template
		contentTmpl, err := template.ParseFiles("src/views/admin/post_form.html")
		if err != nil {
			log.Printf("Content template parsing error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute content template to a buffer
		var contentBuffer strings.Builder
		err = contentTmpl.ExecuteTemplate(&contentBuffer, "post-form", data)
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
		return
	}

	// Update post
	_, err = models.UpdatePost(slug, title, content)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to admin dashboard
	http.Redirect(w, r, "/owner", http.StatusSeeOther)
}

// DeletePostHandler handles the POST /owner/delete/:slug route
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Get slug from URL
	slug := chi.URLParam(r, "slug")

	// Delete post
	err := models.DeletePost(slug)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to admin dashboard
	http.Redirect(w, r, "/owner", http.StatusSeeOther)
}
