package models

import (
	"database/sql"
	"time"
	"word_press/database"
	"fmt"
)

type PrayerSession struct {
	ID        int
	UserID    int
	StartTime time.Time
	EndTime   *time.Time
	Duration  int
}

func CreatePrayerSession(userID int) (int, error) {
	res, err := database.DB.Exec(`
		INSERT INTO prayer_sessions (user_id)
		VALUES (?)
	`, userID)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, fmt.Errorf("failed to retrieve inserted session ID")
	}

	return int(id), nil
}

func EndPrayerSession(sessionID int) error {
	_, err := database.DB.Exec(`
		UPDATE prayer_sessions
		SET 
			end_time = CURRENT_TIMESTAMP,
			duration = CAST(
				strftime('%s','now') - strftime('%s', start_time)
				AS INTEGER
			)
		WHERE id = ?
	`, sessionID)

	return err
}

func GetUserTotalPrayerTime(userID int) (int, error) {
	var total sql.NullInt64

	err := database.DB.QueryRow(`
		SELECT SUM(duration)
		FROM prayer_sessions
		WHERE user_id = ?`,
		userID,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	if total.Valid {
		return int(total.Int64), nil
	}
	return 0, nil
}

func GetActiveSession(userID int) (*PrayerSession, error) {
	row := database.DB.QueryRow(`
		SELECT id, user_id, start_time
		FROM prayer_sessions
		WHERE user_id = ? AND end_time IS NULL
		ORDER BY start_time DESC
		LIMIT 1
	`, userID)

	var s PrayerSession

	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.StartTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

func EndAllActiveSessions(userID int) error {
	_, err := database.DB.Exec(`
		UPDATE prayer_sessions
		SET 
			end_time = CURRENT_TIMESTAMP,
			duration = CAST(
				strftime('%s','now') - strftime('%s', start_time)
				AS INTEGER
			)
		WHERE user_id = ? AND end_time IS NULL
	`, userID)

	return err
}