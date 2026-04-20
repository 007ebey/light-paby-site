package models

import (
	"database/sql"
	"time"

	"word_press/database"
)

type PrayerItem struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}

func CreatePrayerItem(userID int, title, content string) error {
	_, err := database.DB.Exec(`
		INSERT INTO prayer_items (user_id, title, content, created_at)
		VALUES (?, ?, ?, ?)
	`, userID, title, content, time.Now())

	return err
}

func GetPrayerItemsByUser(userID int) ([]PrayerItem, error) {

	rows, err := database.DB.Query(`
		SELECT id, user_id, title, content, created_at
		FROM prayer_items
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]PrayerItem, 0)

	for rows.Next() {
		var p PrayerItem

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func GetPrayerItemByID(id int) (*PrayerItem, error) {

	row := database.DB.QueryRow(`
		SELECT id, user_id, title, content, created_at
		FROM prayer_items
		WHERE id = ?
		LIMIT 1
	`, id)

	var p PrayerItem

	err := row.Scan(
		&p.ID,
		&p.UserID,
		&p.Title,
		&p.Content,
		&p.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func DeletePrayerItem(id int) error {
	_, err := database.DB.Exec(`
		DELETE FROM prayer_items
		WHERE id = ?
	`, id)

	return err
}