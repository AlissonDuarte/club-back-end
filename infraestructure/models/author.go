package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name             string
	Resume           string
	Rate             int
	ProfilePictureID uint        `gorm:"default:null"`
	ProfilePicture   *UserUpload `gorm:"foreignKey:ProfilePictureID"`
	Books            []Book      `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Certified        bool
}

func NewAuthor(name string, resume string, rate int, pictureID uint, certified bool, db *gorm.DB) *Author {
	return &Author{
		Name:             name,
		Resume:           resume,
		Rate:             rate,
		ProfilePictureID: pictureID,
		Certified:        certified,
	}
}

func (a *Author) Save(db *gorm.DB) (uint, error) {
	err := db.Create(a).Error
	if err != nil {
		return 0, err
	}
	return a.ID, nil
}

func (a *Author) Update(db *gorm.DB) error {
	err := db.Save(a).Error
	if err != nil {
		return err
	}
	return nil
}
