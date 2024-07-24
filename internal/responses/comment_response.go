package responses

type UserResponse struct {
	Username       string              `json:"username"`
	ProfilePicture *ProfilePicResponse `json:"profile_picture"`
}

type ProfilePicResponse struct {
	FilePath string `json:"file_path"`
}

type CommentResponse struct {
	CommentID uint         `json:"comment_id"`
	Content   string       `json:"content"`
	CreatedAt string       `json:"created_at"`
	User      UserResponse `json:"user"`
}

type ImageResponse struct {
	FilePath string `json:"file_path"`
}
