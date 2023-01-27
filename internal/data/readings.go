package data

import (
	"context"
	"database/sql"
	"errors"
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

func (m ReadingModel) Get(id int64, userID int64) (*Reading, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, book_name, book_author, total_page_count, current_page, finished, memo, user_id, created_at, updated_at, version
		FROM readings
		WHERE id = ?
	`

	var reading Reading

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&reading.ID,
		&reading.BookName,
		&reading.BookAuthor,
		&reading.TotalPageCount,
		&reading.CurrentPage,
		&reading.Finished,
		&reading.Memo,
		&reading.UserID,
		&reading.CreatedAt,
		&reading.UpdatedAt,
		&reading.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if reading.UserID != userID {
		return nil, ErrNotPermitted
	}

	return &reading, nil
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

func (m ReadingModel) Update(reading *Reading) error {
	query := `
		UPDATE readings SET book_name = ?, book_author = ?, memo = ?, version = version + 1
		WHERE id = ? AND version = ?
	`

	args := []any{
		reading.BookName,
		reading.BookAuthor,
		reading.Memo,
		reading.ID,
		reading.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrEditConflict
	}

	return nil
}

func (m ReadingModel) GetAll(userID int64, filters Filters) ([]*Reading, error) {
	query := `
		SELECT id, book_name, book_author, total_page_count, current_page, finished, memo, user_id, created_at, updated_at, version
		FROM readings
		WHERE user_id = ?
		LIMIT ? OFFSET ?
	`

	args := []any{userID, filters.PageSize, filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	readings := []*Reading{}

	defer rows.Close()

	for rows.Next() {
		var reading Reading

		err := rows.Scan(
			&reading.ID,
			&reading.BookName,
			&reading.BookAuthor,
			&reading.TotalPageCount,
			&reading.CurrentPage,
			&reading.Finished,
			&reading.Memo,
			&reading.UserID,
			&reading.CreatedAt,
			&reading.UpdatedAt,
			&reading.Version,
		)
		if err != nil {
			return nil, err
		}

		readings = append(readings, &reading)
	}

	return readings, nil
}
