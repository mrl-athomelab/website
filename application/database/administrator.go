package database

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Administrator struct {
	gorm.Model
	Username string
	Password string
	UUID     string
}

type Administrators []Administrator

func (a *Administrators) All() (out *Administrators) {
	out = &Administrators{}
	db.Model(&Administrator{}).Order("created_at desc").Scan(out)
	return
}

func (u *Administrator) BeforeSave() error {
	tUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	u.UUID = tUUID.String()
	return nil
}

func (u *Administrator) Save() error {
	return db.Model(u).Save(u).Error
}

func (u *Administrator) Delete() error {
	return db.Model(u).Delete(u).Error
}

func (u *Administrator) Get(mode int8) bool {
	query := true
	switch mode {
	case ByUsernamePassword:
		db.Model(u).Find(u, "username = ? AND password = ?", u.Username, u.Password)
		break
	case ByID:
		db.Model(u).Find(u, "id = ?", u.ID)
		break
	default:
		query = false
	}
	return query && u.ID > 0
}
