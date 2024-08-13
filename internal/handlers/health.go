package handlers

import (
	"gophermart/internal/controllers"
	"gophermart/internal/repository"
	"net/http"
)

type HealthHandlers struct {
	controller controllers.HealthController
}

func NewHealthHandlers(repoGroup *repository.Repos) *HealthHandlers {
	controller := controllers.NewHealthController(repoGroup)
	return &HealthHandlers{
		controller: controller,
	}
}

func (h *HealthHandlers) SetReadiness(state bool) {
	h.controller.SetReadiness(state)
}

func (h *HealthHandlers) SetLiveness(state bool) {
	h.controller.SetLiveness(state)
}

func (h *HealthHandlers) ReadinessState(w http.ResponseWriter, r *http.Request) {
	if h.controller.ReadinessState() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func (h *HealthHandlers) LivenessState(w http.ResponseWriter, r *http.Request) {
	if h.controller.LivenessState() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func (h *HealthHandlers) PingDB(w http.ResponseWriter, r *http.Request) {
	if h.controller.PingDB() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
