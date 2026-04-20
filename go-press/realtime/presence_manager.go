package realtime

import (
	"sync"
	"time"
	"word_press/models"
)

type PresenceManager struct {
	mu       sync.RWMutex
	users    map[int]models.User
	lastSeen map[int]time.Time

	ttl time.Duration
}

func NewPresenceManager(ttl time.Duration) *PresenceManager {
	pm := &PresenceManager{
		users:    make(map[int]models.User),
		lastSeen: make(map[int]time.Time),
		ttl:      ttl,
	}

	go pm.cleanupLoop()

	return pm
}

func (p *PresenceManager) AddOrUpdateUser(u models.User) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.users[u.ID] = u
	p.lastSeen[u.ID] = time.Now()
}

func (p *PresenceManager) RemoveUser(userID int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.users, userID)
	delete(p.lastSeen, userID)
}

func (p *PresenceManager) GetStats() (active int, countries int, cities int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	countrySet := make(map[string]struct{})
	citySet := make(map[string]struct{})

	for _, u := range p.users {
		if u.Country != "" {
			countrySet[u.Country] = struct{}{}
		}
		if u.City != "" {
			citySet[u.City] = struct{}{}
		}
	}

	return len(p.users), len(countrySet), len(citySet)
}

func (p *PresenceManager) GetUsers() []models.User {
	p.mu.RLock()
	defer p.mu.RUnlock()

	users := make([]models.User, 0, len(p.users))

	for _, u := range p.users {
		users = append(users, u)
	}

	return users
}

func (p *PresenceManager) cleanupLoop() {
	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		p.cleanup()
	}
}

func (p *PresenceManager) cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()

	for userID, last := range p.lastSeen {
		if now.Sub(last) > p.ttl {
			delete(p.users, userID)
			delete(p.lastSeen, userID)
		}
	}
}