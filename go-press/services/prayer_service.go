package services

import (
	"errors"
    "log"
	"word_press/models"
)

type PrayerService struct{}

func NewPrayerService() *PrayerService {
	return &PrayerService{}
}

func (s *PrayerService) Start(userID int) (int, error) {

	// Check if already has active session
	active, err := models.GetActiveSession(userID)
	log.Println(err)
	if err != nil {
		return 0, err
	}

	if active != nil {
		// Option 1: return error (strict)
		// return 0, errors.New("session already active")

		// Option 2: auto-close previous (better UX)
		if err := models.EndAllActiveSessions(userID); err != nil {
			return 0, err
		}
	}

	// Create new session
	sessionID, err := models.CreatePrayerSession(userID)
	if err != nil {
		return 0, err
	}

	return sessionID, nil
}

func (s *PrayerService) End(sessionID int) error {

	if sessionID == 0 {
		return errors.New("invalid session id")
	}

	return models.EndPrayerSession(sessionID)
}

func (s *PrayerService) GetTotalTime(userID int) (int, error) {
	return models.GetUserTotalPrayerTime(userID)
}

func (s *PrayerService) GetActive(userID int) (*models.PrayerSession, error) {
	return models.GetActiveSession(userID)
}

func (s *PrayerService) CleanupUserSessions(userID int) error {
	return models.EndAllActiveSessions(userID)
}