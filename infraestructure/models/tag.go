package models

import (
	"errors"

	"gorm.io/gorm"
)

type TagChoices int

const (
	Undefined            TagChoices = 999
	ScienceFiction       TagChoices = 1
	Fantasy              TagChoices = 2
	Mystery              TagChoices = 3
	Thriller             TagChoices = 4
	Horror               TagChoices = 5
	Biography            TagChoices = 6
	Autobiography        TagChoices = 7
	History              TagChoices = 8
	HistoricalFiction    TagChoices = 9
	Poetry               TagChoices = 10
	Drama                TagChoices = 11
	CrimeFiction         TagChoices = 12
	Dystopian            TagChoices = 13
	Utopian              TagChoices = 14
	Adventure            TagChoices = 15
	Teenager             TagChoices = 16
	ChildrenLiterature   TagChoices = 17
	SelfHelp             TagChoices = 18
	LiteraryFiction      TagChoices = 19
	UrbanFantasy         TagChoices = 20
	PsychologicalFiction TagChoices = 21
	Humor                TagChoices = 22
	ReligionAndSpiritual TagChoices = 23
	Memoir               TagChoices = 24
	Erotica              TagChoices = 25
	EspionageFiction     TagChoices = 26
	EpicFantasy          TagChoices = 27
	Essays               TagChoices = 28
	WarFiction           TagChoices = 29
	TechnicalManual      TagChoices = 30
	Sports               TagChoices = 31
	Paranormal           TagChoices = 32
	Noir                 TagChoices = 33
	Mythological         TagChoices = 34
	Mathematics          TagChoices = 35
	Chemistry            TagChoices = 36
	Physics              TagChoices = 37
	Computing            TagChoices = 38
)

var validTagChoices = map[TagChoices]struct{}{
	ScienceFiction:       {},
	Fantasy:              {},
	Mystery:              {},
	Thriller:             {},
	Horror:               {},
	Biography:            {},
	Autobiography:        {},
	History:              {},
	HistoricalFiction:    {},
	Poetry:               {},
	Drama:                {},
	CrimeFiction:         {},
	Dystopian:            {},
	Utopian:              {},
	Adventure:            {},
	Teenager:             {},
	ChildrenLiterature:   {},
	SelfHelp:             {},
	LiteraryFiction:      {},
	UrbanFantasy:         {},
	PsychologicalFiction: {},
	Humor:                {},
	ReligionAndSpiritual: {},
	Memoir:               {},
	Erotica:              {},
	EspionageFiction:     {},
	EpicFantasy:          {},
	Essays:               {},
	WarFiction:           {},
	TechnicalManual:      {},
	Sports:               {},
	Paranormal:           {},
	Noir:                 {},
	Mythological:         {},
	Mathematics:          {},
	Chemistry:            {},
	Physics:              {},
	Computing:            {},
}

type Tag struct {
	gorm.Model
	Name string     `gorm:"not null"`
	Type TagChoices `gorm:"default:999"`
}

func (tag *Tag) BeforeSave(tx *gorm.DB) (err error) {
	if !isValidTagType(tag.Type) {
		return errors.New("ivalid tag type")
	}
	return nil
}

func isValidTagType(tagType TagChoices) bool {
	_, ok := validTagChoices[tagType]
	return ok
}
