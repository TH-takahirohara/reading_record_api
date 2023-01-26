package data

import (
	"context"
	"database/sql"
	"time"
	"unicode/utf8"

	"github.com/TH-takahirohara/reading_record_api/internal/validator"
)

type Reading struct {
	ID             int64     `json:"id"`
	BookName       string    `json:"book_name"`
	BookAuthor     string    `json:"book_author"`
	TotalPageCount int       `json:"total_page_count"`
	CurrentPage    int       `json:"current_page"`
	Finished       bool      `json:"finished"`
	Memo           string    `json:"memo"`
	UserID         int64     `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	Version        int64     `json:"-"`
}

func ValidateBookName(v *validator.Validator, bookName string) {
	v.Check(bookName != "", "book_name", "値を入力してください")
	v.Check(len([]rune(bookName)) <= 500, "book_name", "500文字以内の文字列を入力してください")
}

func ValidateBookAuthor(v *validator.Validator, bookAuthor string) {
	v.Check(bookAuthor != "", "book_author", "値を入力してください")
	v.Check(len([]rune(bookAuthor)) <= 500, "book_author", "500文字以内の文字列を入力してください")
}

func ValidateTotalPageCount(v *validator.Validator, totalPageCount int) {
	v.Check(totalPageCount > 0, "total_page_count", "0より大きい値を入力してください")
	v.Check(totalPageCount <= 50000, "total_page_count", "50000以下の値を入力してください")
}

func ValidateMemo(v *validator.Validator, memo string) {
	v.Check(utf8.RuneCountInString(memo) <= 10000, "memo", "10000文字以内の文字列を入力してください")
}

func ValidateReading(v *validator.Validator, reading *Reading) {
	ValidateBookName(v, reading.BookName)
	ValidateBookAuthor(v, reading.BookAuthor)
	ValidateTotalPageCount(v, reading.TotalPageCount)
	ValidateMemo(v, reading.Memo)
}

type ReadingModel struct {
	DB *sql.DB
}

func (m ReadingModel) Insert(reading *Reading) error {
	query := `
		INSERT INTO readings (book_name, book_author, total_page_count, user_id)
		VALUES (?, ?, ?, ?)
	`

	args := []any{
		reading.BookName,
		reading.BookAuthor,
		reading.TotalPageCount,
		reading.UserID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	reading.ID = id
	return nil
}
