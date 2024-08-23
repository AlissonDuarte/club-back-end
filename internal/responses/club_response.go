package responses

import "time"

type ClubResponse struct {
	ID          uint                   `json:"ID"`
	CreatedAt   time.Time              `json:"CreatedAt"`
	UpdatedAt   time.Time              `json:"UpdatedAt"`
	Name        string                 `json:"Name"`
	Description string                 `json:"Description"`
	ImageID     uint                   `json:"ImageID"`
	OwnerRefer  ClubOwnerReferResponse `json:"OwnerRefer"`
	Users       []ClubUsersResponse    `json:"Users"`
}

type ClubOwnerReferResponse struct {
	ID        uint   `json:"ID"`
	Username  string `json:"Username"`
	Gender    string `json:"Gender"`
	BirthDate string `json:"BirthDate"`
}

type ClubUsersResponse struct {
	ID        uint   `json:"ID"`
	Username  string `json:"Username"`
	Gender    string `json:"Gender"`
	BirthDate string `json:"BirthDate"`
}

// daqui para baixo é referente a serializacao da resposta da requisicao GET /club/id/post/id
type ClubPostResponse struct {
	ID      uint         `json:"ID"`
	Title   string       `json:"Title"`
	Content string       `json:"Content"`
	User    ClubPostUser `json:"User"`
	ClubID  uint         `json:"ClubID"`
}

type ClubPostUser struct {
	ID             uint                       `json:"ID"`
	Name           string                     `json:"Name"`
	Username       string                     `json:"Username"`
	ProfilePicture ClubPostUserProfilePicture `json:"ProfilePicture"`
}

type ClubPostUserProfilePicture struct {
	ID uint `json:"ID"`
}
