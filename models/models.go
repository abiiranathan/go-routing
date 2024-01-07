package models

import "time"

type Post struct {
	ID        int       `json:"id" gorm:"primaryKey; autoIncrement; not null"`
	Title     string    `json:"title" gorm:"uniqueIndex:idx_post_title; not null"`
	Slug      string    `json:"slug" gorm:"uniqueIndex:idx_post_slug; not null"`
	Body      string    `json:"body" gorm:"type:text; not null"`
	Author    string    `json:"author" gorm:"default:'Anonymous'"`
	CreatedAt time.Time `json:"created_at"`
}

func (post Post) String() string {
	return post.Title
}
