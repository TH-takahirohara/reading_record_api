package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrNotPermitted   = errors.New("not permitted")
)

type Models struct {
	Tokens          TokenModel
	Users           UserModel
	Readings        ReadingModel
	DailyProgresses DailyProgressModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tokens:          TokenModel{DB: db},
		Users:           UserModel{DB: db},
		Readings:        ReadingModel{DB: db},
		DailyProgresses: DailyProgressModel{DB: db},
	}
}
