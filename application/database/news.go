package database

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type News struct {
	gorm.Model
	UUID     string
	AuthorID uint
	Title    string
	Content  string
}

type NewsList []News

func (a *NewsList) All() (out *NewsList) {
	out = &NewsList{}
	db.Model(tables["news"]).Order("created_at desc").Scan(out)
	return
}

func (u *News) BeforeSave() error {
	tUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	u.UUID = tUUID.String()
	return nil
}

func (u *News) Save() error {
	return db.Model(u).Save(u).Error
}

func (u *News) Delete() error {
	return db.Model(u).Delete(u).Error
}

func (u *News) Get(mode int8) bool {
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
