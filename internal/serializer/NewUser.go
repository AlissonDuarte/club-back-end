package serializer
// campo de senha write only
type UserSerializer struct {
	Name      	string `json:"name" validate:"required"`
	BirthDate 	string `json:"birth_date" validate:"required"`
	Cep       	string `json:"cep" validate:"required"`
	Passwd		string `json:"passwd" validate:"min=8,max=72" `
	Email     	string `json:"email" validate:"required,email"`
	Phone     	string `json:"phone" validate:"required"`
}
