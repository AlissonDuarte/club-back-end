package responses

import "time"

type FeedProfilePicture struct {
	ID       uint   `json:"id"`
	FilePath string `json:"filePath"`
}

type FeedUserData struct {
	ID               uint               `json:"id"`
	Name             string             `json:"name"`
	Username         string             `json:"username"`
	ProfilePictureID uint               `json:"profilePictureID"`
	ProfilePicture   FeedProfilePicture `json:"profilePicture"`
}

type FeedResponse struct {
	ID           uint         `json:"id"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	UserID       uint         `json:"userID"`
	User         FeedUserData `json:"user"`
	ImageID      uint         `json:"imageID"`
	CommentCount int64        `json:"commentCount"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}
