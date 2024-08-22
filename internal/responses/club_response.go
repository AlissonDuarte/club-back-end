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
