package serializer

// campo de senha write only
type UserSerializer struct {
	Name      string `json:"name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	BirthDate string `json:"birth_date" validate:"required"`
	Gender    string `json:"gender" validate:"required"`
	Passwd    string `json:"passwd" validate:"min=8,max=72" `
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required"`
	Bio       string `json:"bio"`
}

type UserUpdateSerializer struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Passwd    string `json:"passwd"`
	Email     string `json:"email" validate:"email"`
	Phone     string `json:"phone"`
	Bio       string `json:"bio"`
}

type UserChangePasswordSerliazer struct {
	Id             int    `json:"user_id" validate:"required"`
	OldPasswd      string `json:"old_passwd" validate:"required"`
	NewPasswd      string `json:"new_passwd" validate:"required"`
	NewPasswdCheck string `json:"new_passwd_check" validate:"required"`
}
