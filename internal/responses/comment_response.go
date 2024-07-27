package responses

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type CommentResponse struct {
	CommentID uint         `json:"comment_id"`
	Content   string       `json:"content"`
	Updated   bool         `json:"updated"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
	User      UserResponse `json:"user"`
}
