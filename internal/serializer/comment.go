package serializer

type CommentSerializer struct {
	UserID  int    `json:"userID" validate:"required"`
	PostID  int    `json:"PostID" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type CommentUpdateSerializer struct {
	CommentID int    `json:"commentID" validate:"required"`
	UserID    int    `json:"userID" validate:"required"`
	PostID    int    `json:"PostID" validate:"required"`
	Content   string `json:"content" validate:"required"`
}
