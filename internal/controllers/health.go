package controllers

import (
	"context"
	"gophermart/internal/repository"
	"sync"
)

type HealthControllerImpl struct {
	repos *repository.Repos

	readinessMu sync.RWMutex
	readiness   bool
	livenessMu  sync.RWMutex
	liveness    bool
}

func NewHealthController(repos *repository.Repos) *HealthControllerImpl {
	return &HealthControllerImpl{repos: repos}
}

func (h *HealthControllerImpl) SetReadiness(state bool) {
	h.readinessMu.Lock()
	defer h.readinessMu.Unlock()
	h.readiness = state
}

func (h *HealthControllerImpl) SetLiveness(state bool) {
	h.livenessMu.Lock()
	defer h.livenessMu.Unlock()
	h.liveness = state
}

func (h *HealthControllerImpl) ReadinessState() bool {
	h.readinessMu.RLock()
	defer h.readinessMu.RUnlock()
	return h.readiness
}

func (h *HealthControllerImpl) LivenessState() bool {
	h.livenessMu.RLock()
	defer h.livenessMu.RUnlock()
	return h.liveness
}

func (h *HealthControllerImpl) PingDB() bool {
	return h.repos.HealthRepo.PingDB(context.Background())
}
