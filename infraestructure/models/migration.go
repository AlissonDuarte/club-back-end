package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	db.Set("gorm:save_associations", false)
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	err := db.AutoMigrate(
		&User{},
		&Club{},
		&UserUpload{},
		&UserUploadPost{},
		&Post{},
		&Comment{},
		&Author{},
		&Rate{},
		&Book{},
		&Tag{},
	)
	if err != nil {
		return err
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	db.Set("gorm:save_associations", true)
	return nil
}
