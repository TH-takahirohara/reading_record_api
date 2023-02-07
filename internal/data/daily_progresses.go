package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type DailyProgress struct {
	ID        int64     `json:"id"`
	ReadDate  time.Time `json:"readDate"`
	ReadPage  int       `json:"readPage"`
	ReadingID int64     `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Version   int64     `json:"-"`
}

type DailyProgressModel struct {
	DB *sql.DB
}

func (m DailyProgressModel) GetLatest(readingId int64) (*DailyProgress, error) {
	query := `
		SELECT id, read_date, read_page
		FROM daily_progresses
		WHERE reading_id = ?
		ORDER BY read_date DESC
		LIMIT 1
	`

	dailyProgress := DailyProgress{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, readingId).Scan(
		&dailyProgress.ID,
		&dailyProgress.ReadDate,
		&dailyProgress.ReadPage,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		default:
			return nil, err
		}
	}

	return &dailyProgress, nil
}

func (m DailyProgressModel) Insert(dailyProgress *DailyProgress) error {
	query := `
		INSERT INTO daily_progresses (read_date, read_page, reading_id)
		VALUES (?, ?, ?)
	`

	args := []any{dailyProgress.ReadDate, dailyProgress.ReadPage, dailyProgress.ReadingID}

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

	dailyProgress.ID = id
	return nil
}

func (m DailyProgressModel) GetAll(readingID int64) ([]*DailyProgress, error) {
	query := `
		SELECT id, read_date, read_page
		FROM daily_progresses
		WHERE reading_id = ?
		ORDER BY read_date
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, readingID)
	if err != nil {
		return nil, err
	}

	dailyProgresses := []*DailyProgress{}

	defer rows.Close()

	for rows.Next() {
		var dailyProgress DailyProgress

		err := rows.Scan(
			&dailyProgress.ID,
			&dailyProgress.ReadDate,
			&dailyProgress.ReadPage,
		)
		if err != nil {
			return nil, err
		}

		dailyProgresses = append(dailyProgresses, &dailyProgress)
	}

	return dailyProgresses, nil
}

func (m DailyProgressModel) Delete(id int64) error {
	query := `
		DELETE FROM daily_progresses
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
