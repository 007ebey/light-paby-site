package models

import "word_press/database"

type Slider struct {
	ID int 
	Title string
	Subtitle string 
	Image string
	ButtonText string 
	ButtonURL string
	Status string
}

func CreateSlider(
	title, subtitle, image, buttonText, buttonURL, status string,
) error {
	query := `
	 INSERT INTO sliders
	 (title, subtitle, image, button_text, button_url, status)
	 VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := database.DB.Exec(
		query,
		title,
		subtitle,
		image,
		buttonText,
		buttonURL,
		status,
	)

	return err
}