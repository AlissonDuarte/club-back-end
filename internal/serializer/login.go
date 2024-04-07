package serializer

type UserLoginSerializer struct {
	Email	string `json:"email" validate:"required,email"`
	Passwd	string `json:"passwd" validate:"required"`
}
