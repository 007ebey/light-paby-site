package services

import (
	"word_press/models"
	"word_press/realtime"
)

type PresenceService struct {
	manager *realtime.PresenceManager
}

type PresenceStats struct {
	ActiveUsers int
	Countries   int
	Cities      int
}

func NewPresenceService(m *realtime.PresenceManager) *PresenceService {
	return &PresenceService{
		manager: m,
	}
}

func (s *PresenceService) Touch(user *models.User) {
	if user == nil {
		return
	}

	s.manager.AddOrUpdateUser(*user)
}

func (s *PresenceService) Leave(userID int) {
	if userID == 0 {
		return
	}

	s.manager.RemoveUser(userID)
}

func (s *PresenceService) GetStats() PresenceStats {
	active, countries, cities := s.manager.GetStats()

	return PresenceStats{
		ActiveUsers: active,
		Countries:   countries,
		Cities:      cities,
	}
}

func (s *PresenceService) GetUsers() []models.User {
	return s.manager.GetUsers()
}