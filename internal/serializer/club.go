package serializer

import (
	"clube/infraestructure/models"
	"clube/internal/responses"
)

type GroupSerializer struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Owner       uint   `json:"owner" validate:"required"`
	Users       []int
}

func ClubGetSerialize(club *models.Club) (*responses.ClubResponse, error) {
	var clubResponse responses.ClubResponse

	clubResponse.ID = club.ID
	clubResponse.Name = club.Name
	clubResponse.Description = club.Description
	clubResponse.ImageID = club.ImageID
	clubResponse.CreatedAt = club.CreatedAt
	clubResponse.UpdatedAt = club.UpdatedAt
	clubResponse.OwnerRefer = responses.ClubOwnerReferResponse{
		ID:        club.OwnerID,
		Username:  club.OwnerRefer.Username,
		Gender:    club.OwnerRefer.Gender,
		BirthDate: club.OwnerRefer.BirthDate,
	}

	for _, user := range club.Users {
		clubResponse.Users = append(clubResponse.Users, responses.ClubUsersResponse{
			ID:        user.ID,
			Username:  user.Username,
			Gender:    user.Gender,
			BirthDate: user.BirthDate,
		})
	}

	return &clubResponse, nil
}

func ClubsGetSerialize(clubs []*models.Club) ([]*responses.ClubResponse, error) {
	var clubsResponse []*responses.ClubResponse

	for _, club := range clubs {
		clubResponse, err := ClubGetSerialize(club)
		if err != nil {
			return nil, err
		}
		clubsResponse = append(clubsResponse, clubResponse)
	}

	return clubsResponse, nil
}
