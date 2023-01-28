package data

import (
	"context"
	"database/sql"
	"time"
)

type DailyProgress struct {
	ID        int64     `json:"id"`
	ReadDate  time.Time `json:"read_date"`
	ReadPage  int       `json:"read_page"`
	ReadingID int64     `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Version   int64     `json:"-"`
}

type DailyProgressModel struct {
	DB *sql.DB
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
