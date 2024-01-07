// Uses apigen, a custom service generator for GORM.
// github.com/abiiranathan/apigen
//
// GODEBUG=httpmuxgo121=0 gotip run main.go
// Author: Dr. Abiira Nathan
// Date:   07 January 2024
package main

import (
	"hello/api"
	"hello/models"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Post{}); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	wrapped := api.Use(mux, api.LoggerMiddleware)

	// setup API routes
	api.New(db, mux)

	log.Fatalln(http.ListenAndServe(":8080", wrapped))
}
