package serializer
// campo de senha write only
type UserSerializer struct {
	Name      	string `json:"name" validate:"required"`
	Username  	string `json:"username" validate:"required"`
	BirthDate 	string `json:"birth_date" validate:"required"`
	Gender		string `json:"gender" validate:"required"`
	Passwd		string `json:"passwd" validate:"min=8,max=72" `
	Email     	string `json:"email" validate:"required,email"`
	Phone     	string `json:"phone" validate:"required"`
}

