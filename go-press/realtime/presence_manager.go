package realtime

import (
	"sync"
	"time"
	"word_press/models"
)

const numShards = 32

type presenceEntry struct {
	user models.User
	lastSeen int64
}

type shard struct {
	mu sync.RWMutex
	users map[int]*presenceEntry
}

type PresenceManager struct {
	shards   []shard
	ttl      time.Duration
	stopCh   chan struct{}
}

func NewPresenceManager(ttl time.Duration) *PresenceManager {
	pm := &PresenceManager{
	 ttl:    ttl,
	 shards: make([]shard, numShards),
	 stopCh: make(chan struct{}),
    }

	for i := 0; i < numShards; i++ {
	  pm.shards[i] = shard{
	  	users: make(map[int]*presenceEntry),
	  }
    }

	go pm.cleanupLoop()

	return pm
}

func (p *PresenceManager) getShard(userID int) *shard {
	return &p.shards[userID%numShards]
}

func (p *PresenceManager) AddOrUpdateUser(u models.User) {
	s := p.getShard(u.ID)

	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.users[u.ID]; ok {
	  existing.user = u
	  existing.lastSeen = time.Now().UnixNano()
	  return
    }

    s.users[u.ID] = &presenceEntry{
	  user:     u,
	  lastSeen: time.Now().UnixNano(),
    }
}

func (p *PresenceManager) RemoveUser(userID int) {
	s := p.getShard(userID)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, userID)
}

func (p *PresenceManager) GetStats() (active int, countries int, cities int) {
	countrySet := make(map[string]struct{})
	citySet := make(map[string]struct{})

    for i := 0; i < numShards; i++ {
		s := p.shards[i]
		s.mu.RLock()
		for _, entry := range s.users {
			active++

			u := entry.user

			if u.Country != "" {
				countrySet[u.Country] = struct{}{}
			}

			if u.City != "" {
				citySet[u.City] = struct{}{}
			}
		}
		s.mu.RUnlock()
	}

	return active, len(countrySet), len(citySet)
}

func (p *PresenceManager) GetUsers() []models.User {
	total := 0
	for i := 0; i < numShards; i++ {
		s := p.shards[i]
		s.mu.RLock()
		total += len(s.users)
		s.mu.RUnlock()
	}
	users := make([]models.User, 0, total)
	return users
}

func (p *PresenceManager) cleanupLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.cleanup()
		case <-p.stopCh:
			return
		}
	}
}

func (p *PresenceManager) cleanup() {
	now := time.Now().UnixNano()

	for i := 0; i < numShards; i++ {
		s := p.shards[i]
		s.mu.Lock()

		for userID, entry := range s.users {
			if now - entry.lastSeen > p.ttl.Nanoseconds() {
				delete(s.users, userID)
			}
		}

		s.mu.Unlock()
	}
}

func (p *PresenceManager) Stop() {
	close(p.stopCh)
}