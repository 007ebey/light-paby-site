package models

import (
	"time"
	"word_press/database"
)

type VisionNote struct {
	ID        int
	UserID    int
	Content   string
	CreatedAt time.Time
}

func CreateVisionNote(userID int, content string) error {
	_, err := database.DB.Exec(`
		INSERT INTO vision_notes (user_id, content, created_at)
		VALUES (?, ?, ?)
	`, userID, content, time.Now())

	return err
}

func GetVisionNotesByUser(userID int) ([]VisionNote, error) {

	rows, err := database.DB.Query(`
		SELECT id, user_id, content, created_at
		FROM vision_notes
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []VisionNote

	for rows.Next() {
		var n VisionNote
		rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Content,
			&n.CreatedAt,
		)
		notes = append(notes, n)
	}

	return notes, nil
}