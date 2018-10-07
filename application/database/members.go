package database

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Member struct {
	gorm.Model
	FirstName       string
	LastName        string
	UUID            string
	Rule            string
	Biography       string
	SocialMediaType string
	SocialMediaLink string
}

type Members []Member

func (a *Members) All() (out *Members) {
	out = &Members{}
	db.Model(tables["members"]).Order("created_at desc").Scan(out)
	return
}

func (u *Member) BeforeSave() error {
	tUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	u.UUID = tUUID.String()
	return nil
}

func (u *Member) Save() error {
	return db.Model(u).Save(u).Error
}

func (u *Member) Delete() error {
	return db.Model(u).Delete(u).Error
}

func (u *Member) Get(mode int8) bool {
	query := true
	switch mode {
	case ByID:
		db.Model(u).Find(u, "id = ?", u.ID)
		break
	default:
		query = false
	}
	return query && u.ID > 0
}
