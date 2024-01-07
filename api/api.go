package api

import (
	"encoding/json"
	"fmt"
	"hello/svc"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriter) Write(data []byte) (n int, err error) {
	return w.ResponseWriter.Write(data)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func SendHTMLError(w http.ResponseWriter, error error, code ...int) {
	statusCode := http.StatusInternalServerError
	if len(code) > 0 {
		statusCode = code[0]
	}

	w.Header().Set("Content-Type", "text/html")
	http.Error(w, error.Error(), statusCode)
}

func SendJSONError(w http.ResponseWriter, err error, key string, code ...int) {
	statusCode := http.StatusInternalServerError
	if len(code) > 0 {
		statusCode = code[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"%s": "%s"}`, key, err.Error())
}

func SendJson(w http.ResponseWriter, data interface{}, code ...int) {
	statusCode := http.StatusOK
	if len(code) > 0 {
		statusCode = code[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func Use(mux *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = mux

	for _, mw := range middlewares {
		s = mw(s)
	}
	return s
}

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a custom response writer
		cw := &ResponseWriter{w, http.StatusOK}

		// call the handler
		h.ServeHTTP(cw, r)

		// log the request
		log.Printf("%s %s %d", r.Method, r.URL.Path, cw.statusCode)
	})
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from home page!")
}

// New creates a new API service.
func New(db *gorm.DB, mux *http.ServeMux) {
	s := svc.NewService(db)

	mux.HandleFunc("/{$}", homePage)

	// List all posts
	mux.HandleFunc("GET /posts", listPosts(s))

	// Create post
	mux.HandleFunc("POST /posts/create", createPost(s))

	// Get post detail
	mux.HandleFunc("GET /posts/{id}", getPost(s))

	// Get post detail by slug
	mux.HandleFunc("GET /posts/slug/{slug}", getPostBySlug(s))

	// Update post
	mux.HandleFunc("PUT /posts/{id}", updatePost(s))

	// Delete post
	mux.HandleFunc("DELETE /posts/{id}", deletePost(s))

}
