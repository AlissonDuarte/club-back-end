package serializer

type PostSerializer struct {
	Title   string `json:"title"`
	Content string `json:"content" validate:"required"`
	UserID  uint   `json:"userID" validate:"required"`
}
