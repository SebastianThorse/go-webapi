package models

import "time"

type Post struct {
	Title       string    `json:"title"`
	Description string    `json:"desc"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
	AuthorName  string    `json:"author_name"`
}

type PostRequest struct {
	Title       string `json:"title"`
	Description string `json:"Description"`
	UpdatedAt   string `json:"UpdatedAt"`
	AuthorName  string `json:"AuthorName"`
}
