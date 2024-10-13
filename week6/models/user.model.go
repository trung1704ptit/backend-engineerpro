package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Age          int64     `json:"age"`
	Password     string    `gorm:"not null"`
	Role         string    `gorm:"type:varchar(255);not null"`
	Provider     string    `gorm:"not null"`
	ProfileImage string    `gorm:"not null"`
	Post         []Post    `gorm:"foreignKey:UserID"`
	Verified     bool      `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`

	// Many-to-Many Relationship for Followers/Following
	Followers []*User `gorm:"many2many:user_followers;joinForeignKey:FollowerID;JoinReferences:FollowingID"` // other users following current user
	Following []*User `gorm:"many2many:user_followers;joinForeignKey:FollowingID;JoinReferences:FollowerID"` // current user following other users
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Age             int64  `json:"age" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	ProfileImage    string `json:"profile_image"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	Age          int64     `json:"age"`
	Role         string    `json:"role,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty"`
	Provider     string    `json:"provider"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Follow struct for the many-to-many relationship
type UserFollower struct {
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	CreatedAt   time.Time `gorm:"not null"`
}
