package responses

import "time"

type PostDataResponse struct {
	ID        uint           `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Updated   bool           `json:"updated"`
	Image     PostImage      `json:"image"`
	User      PostUser       `json:"user"`
	Comments  []PostComments `json:"comments"`
}

type PostImage struct {
	ID       uint   `json:"id"`
	FilePath string `json:"filePath"`
}

type PostUserProfilePicture struct {
	ID       uint   `json:"id"`
	FilePath string `json:"filePath"`
}

type PostUser struct {
	ID             uint                   `json:"id"`
	Name           string                 `json:"name"`
	Username       string                 `json:"username"`
	ProfilePicture PostUserProfilePicture `json:"profile_picture"`
}

type PostComments struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Updated   bool      `json:"updated"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	User      PostUser  `json:"user"`
}
