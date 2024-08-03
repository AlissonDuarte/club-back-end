package responses

type UserGroup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type UserGetResponse struct {
	ID               uint        `json:"id"`
	Name             string      `json:"name"`
	Username         string      `json:"username"`
	Gender           string      `json:"gender"`
	BirthDate        string      `json:"birthDate"`
	Email            string      `json:"email"`
	Phone            string      `json:"phone"`
	Bio              string      `json:"bio"`
	ProfilePictureID uint        `json:"profilePictureId"`
	Clubs            []UserGroup `json:"clubs"`
}
