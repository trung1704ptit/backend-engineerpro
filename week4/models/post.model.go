package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title     string    `gorm:"uniqueIndex;not null" json:"title,omitempty"`
	Content   string    `gorm:"not null" json:"content,omitempty"`
	Image     string    `gorm:"not null" json:"image,omitempty"`
	UserID    uuid.UUID `gorm:"not null" json:"user,omitempty"`
	Comments  []Comment `gorm:"foreignKey:PostID"`
	Likes     []Like    `gorm:"foreignKey:PostID"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreatePostRequest struct {
	Title     string    `json:"title"  binding:"required"`
	Content   string    `json:"content" binding:"required"`
	Image     string    `json:"image" binding:"required"`
	UserID    string    `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UpdatePost struct {
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	Image     string    `json:"image,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Comment struct {
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Content   string    `gorm:"not null" json:"body,omitempty"`
	PostID    uuid.UUID `gorm:"not null" json:"post_id,omitempty"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Like struct {
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	PostID    uuid.UUID `gorm:"not null" json:"post_id,omitempty"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
