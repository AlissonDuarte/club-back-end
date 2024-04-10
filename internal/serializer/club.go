package serializer

type GroupSerializer struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Owner       uint   `json:"owner" validate:"required"`
	Users       []int
}
