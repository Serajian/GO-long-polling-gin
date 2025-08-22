package longpolling

import "sync"

// PollingManager manages pending long-poll clients by ID.
type PollingManager struct {
	mu      sync.Mutex
	clients map[string]chan string
}

func NewPollingManager() *PollingManager {
	return &PollingManager{
		clients: make(map[string]chan string),
	}
}

// RegisterClient creates a waiting channel for the given client ID.
// If a channel already exists, it is closed and replaced.
func (pm *PollingManager) RegisterClient(clientID string) chan string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// If client already exists, close its channel to avoid leak
	if old, exists := pm.clients[clientID]; exists {
		close(old)
	}

	ch := make(chan string, 1)
	pm.clients[clientID] = ch
	return ch
}

// SendMessage delivers a message to a specific client if they are waiting,
// then removes the client entry.
func (pm *PollingManager) SendMessage(clientID, message string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if ch, exists := pm.clients[clientID]; exists {
		select {
		case ch <- message: // try to send
		default: // channel full → drop
		}
		close(ch)
		delete(pm.clients, clientID)
	}
}

// RemoveClient deletes a waiting client (used on timeout/abort cleanup).
func (pm *PollingManager) RemoveClient(clientID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if ch, exists := pm.clients[clientID]; exists {
		close(ch)
		delete(pm.clients, clientID)
	}
}
