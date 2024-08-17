package serializer

import (
	"clube/infraestructure/models"
	"clube/internal/responses"
)

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

func PostGetSerialize(post *models.Post) (*responses.PostDataResponse, error) {
	var responsePost responses.PostDataResponse
	responsePost.ID = post.ID
	responsePost.Title = post.Title
	responsePost.Content = post.Content
	responsePost.Image = responses.PostImage{
		ID:       post.Image.ID,
		FilePath: post.Image.FilePath,
	}
	responsePost.User = responses.PostUser{
		ID:       post.User.ID,
		Name:     post.User.Name,
		Username: post.User.Username,
		ProfilePicture: responses.PostUserProfilePicture{
			ID:       post.User.ProfilePicture.ID,
			FilePath: post.User.ProfilePicture.FilePath,
		},
	}

	for _, comment := range post.Comments {
		responsePost.Comments = append(responsePost.Comments, responses.PostComments{
			ID:        comment.ID,
			Content:   comment.Content,
			UpdatedAt: comment.UpdatedAt,
			CreatedAt: comment.CreatedAt,
			User: responses.PostUser{
				ID:       comment.User.ID,
				Name:     comment.User.Name,
				Username: comment.User.Username,
				ProfilePicture: responses.PostUserProfilePicture{
					ID:       comment.User.ProfilePicture.ID,
					FilePath: comment.User.ProfilePicture.FilePath,
				},
			},
		})
	}

	return &responsePost, nil

}
