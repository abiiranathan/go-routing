package api

import (
	"encoding/json"
	"fmt"
	"hello/models"
	"hello/svc"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func sluggify(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "-")
}

func listPosts(s *svc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := s.PostService.GetAll()
		if err != nil {
			SendHTMLError(w, err)
			return
		}
		json.NewEncoder(w).Encode(posts)
	}
}

func createPost(s *svc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse request body
		var post models.Post
		body, err := io.ReadAll(r.Body)
		if err != nil {
			SendHTMLError(w, err)
			return
		}

		// Unmarshal json
		err = json.Unmarshal(body, &post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// generate slug
		post.Slug = sluggify(post.Title)

		// create post
		err = s.PostService.Create(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// send response
		json.NewEncoder(w).Encode(post)
	}
}

func updatePost(s *svc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract id from path
		id := r.PathValue("id")
		postId, err := strconv.Atoi(id)
		if err != nil {
			SendHTMLError(w, fmt.Errorf("invalid post ID"), http.StatusNotFound)
			return
		}

		// fetch post from db
		post, err := s.PostService.Get(postId)
		if err != nil {
			SendHTMLError(w, err, http.StatusNotFound)
			return
		}

		// parse request body
		var payload models.Post
		body, err := io.ReadAll(r.Body)
		if err != nil {
			SendHTMLError(w, err)
			return
		}

		// Unmarshal json
		err = json.Unmarshal(body, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// update post
		// slug should not be updated to avoid breaking links
		post.Title = payload.Title
		post.Body = payload.Body
		post.Author = payload.Author

		// create post
		updatedPost, err := s.PostService.Update(post.ID, &post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// send response
		json.NewEncoder(w).Encode(updatedPost)
	}
}

func deletePost(s *svc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract id from path
		id := r.PathValue("id")
		postId, err := strconv.Atoi(id)
		if err != nil {
			SendHTMLError(w, fmt.Errorf("invalid post ID"), http.StatusNotFound)
			return
		}

		// delete post
		err = s.PostService.Delete(postId)
		if err != nil {
			SendHTMLError(w, err, http.StatusNotFound)
			return
		}
	}
}

func getPost(s *svc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract id from path
		id := r.PathValue("id")
		postId, err := strconv.Atoi(id)
		if err != nil {
			SendHTMLError(w, fmt.Errorf("invalid post ID"), http.StatusNotFound)
			return
		}

		// fetch post
		post, err := s.PostService.Get(postId)
		if err != nil {
			SendHTMLError(w, err, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(post)
	}
}

func getPostBySlug(s *svc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract id from path
		slug := r.PathValue("slug")

		// fetch post by slug(its an index)
		post, err := s.PostService.FindOne(svc.Where("slug = ?", slug))
		if err != nil {
			SendHTMLError(w, err, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(post)
	}
}
