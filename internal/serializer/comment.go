package serializer

type CommentSerializer struct {
	UserID  int    `json:"userId" validate:"required"`
	PostID  int    `json:"postId" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type CommentUpdateSerializer struct {
	CommentID int    `json:"commentId" validate:"required"`
	UserID    int    `json:"userId" validate:"required"`
	PostID    int    `json:"postId" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

type CommentDeleteSerializer struct {
	CommentID int `json:"commentId" validate:"required"`
	UserID    int `json:"userId" validate:"required"`
	PostID    int `json:"postId" validate:"required"`
}
