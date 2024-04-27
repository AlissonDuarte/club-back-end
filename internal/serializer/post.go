package serializer

type PostSerializer struct {
	Title   string `json:"title"`
	Content string `json:"content" validate:"required"`
	UserID  uint   `json:"userID" validate:"required"`
}

type PostDeleteSerializer struct {
	PostID uint `json:"postID" validate:"required"`
	UserID uint `json:"userID" validate:"required"`
}

type PostClubDeleteSerializer struct {
	PostID uint `json:"postID" validate:"required"`
	UserID uint `json:"userID" validate:"required"`
	ClubID uint `json:"clubID" validate:"required"`
}
