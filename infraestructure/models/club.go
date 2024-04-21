package models

import (
	"gorm.io/gorm"
)

// Grupo
type Club struct {
	gorm.Model
	Name        string
	Description string
	OwnerID     uint
	OwnerRefer  *User   `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Users       []*User `gorm:"many2many:user_club;constraint:OnDelete:CASCADE"`
}

func NewClub(name string, description string, userIds []int, owner uint, db *gorm.DB) *Club {
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
		OwnerID:     owner,
	}
}

func (c *Club) Save(db *gorm.DB) error {
	return db.Create(c).Error
}

func (c *Club) Update(db *gorm.DB) error {
	return db.Save(c).Error
}

func ClubGetById(db *gorm.DB, id int) (*Club, error) {
	var club Club
	err := db.Preload("OwnerRefer",
		func(tx *gorm.DB) *gorm.DB {
			return tx.Select("ID", "Username", "Gender", "BirthDate")
		}).Preload("Users", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("ID", "Username", "Gender", "BirthDate")
	}).First(&club, id).Error

	if err != nil {
		return nil, err
	}

	return &club, nil
}
