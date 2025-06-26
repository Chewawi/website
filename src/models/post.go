package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"chewawi_web/src/database"
)

type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}

// GetAllPosts retrieves all posts from the database
func GetAllPosts() ([]Post, error) {
	rows, err := database.DB.Query("SELECT id, title, content, slug, created FROM posts ORDER BY created DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetPostBySlug retrieves a post by its slug
func GetPostBySlug(slug string) (Post, error) {
	var post Post
	err := database.DB.QueryRow("SELECT id, title, content, slug, created FROM posts WHERE slug = $1", slug).
		Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return Post{}, errors.New("post not found")
		}
		return Post{}, err
	}
	return post, nil
}

// CreatePost creates a new post
func CreatePost(title, content string) (Post, error) {
	// Generate slug from title
	slug := generateSlug(title)
	
	// Check if slug already exists
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM posts WHERE slug = $1", slug).Scan(&count)
	if err != nil {
		return Post{}, err
	}
	
	// If slug exists, append a number
	if count > 0 {
		slug = slug + "-" + time.Now().Format("20060102150405")
	}
	
	// Insert post
	var post Post
	err = database.DB.QueryRow(
		"INSERT INTO posts (title, content, slug) VALUES ($1, $2, $3) RETURNING id, title, content, slug, created",
		title, content, slug,
	).Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.Created)
	
	if err != nil {
		return Post{}, err
	}
	
	return post, nil
}

// UpdatePost updates an existing post
func UpdatePost(slug string, title, content string) (Post, error) {
	// Check if post exists
	_, err := GetPostBySlug(slug)
	if err != nil {
		return Post{}, err
	}
	
	// Generate new slug if title changed
	newSlug := generateSlug(title)
	
	// Update post
	var post Post
	err = database.DB.QueryRow(
		"UPDATE posts SET title = $1, content = $2, slug = $3 WHERE slug = $4 RETURNING id, title, content, slug, created",
		title, content, newSlug, slug,
	).Scan(&post.ID, &post.Title, &post.Content, &post.Slug, &post.Created)
	
	if err != nil {
		return Post{}, err
	}
	
	return post, nil
}

// DeletePost deletes a post by its slug
func DeletePost(slug string) error {
	result, err := database.DB.Exec("DELETE FROM posts WHERE slug = $1", slug)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("post not found")
	}
	
	return nil
}

// generateSlug creates a URL-friendly slug from a title
func generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	
	// Remove consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	
	return slug
}