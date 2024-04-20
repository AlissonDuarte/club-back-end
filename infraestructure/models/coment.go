package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
	Content string
	Updated bool
}

func (c *Comment) Save(db *gorm.DB) (uint, error) {

	err := db.Create(c).Error
	if err != nil {
		return 0, err
	}

	return c.ID, nil

}

func (p *Comment) Update(db *gorm.DB) error {
	err := db.Save(p).Error

	if err != nil {
		return err
	}

	return nil
}
