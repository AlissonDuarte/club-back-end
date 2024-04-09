package models

import "gorm.io/gorm"

// Grupo
type Club struct {
	gorm.Model
	Name        string
	Description string
	Owner       uint
	Users       []*User `gorm:"many2many:user_club;"`
}

func NewGroup(name string, description string, userIds []int, db *gorm.DB) *Club {
	users := []*User{}

	for _, userID := range userIds {
		user := &User{}
		db.First(&user, userID)

		if user.ID != 0 {
			users = append(users, user)
		}
	}

	return &Club{
		Name:        name,
		Description: description,
		Users:       users,
	}
}

func (c *Club) Save(db *gorm.DB) error {
	return db.Create(c).Error
}

func (c *Club) Update(db *gorm.DB) error {
	return db.Save(c).Error
}
