package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) setupRoutes() *mux.Router {
	m := mux.NewRouter()

	// Health handlers
	m.HandleFunc("/ping", s.health.PingDB)
	m.HandleFunc("/liveness", s.health.LivenessState)
	m.HandleFunc("/readiness", s.health.ReadinessState)

	// Auth handlers
	m.HandleFunc("/api/user/register", s.auth.Register).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
	m.HandleFunc("/api/user/login", s.auth.Login).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")

	//  Orders handlers
	m.HandleFunc("/api/user/orders", s.orders.Create).
		Methods(http.MethodPost).
		Headers("Content-Type", "text/plain")
	m.HandleFunc("/api/user/orders", s.orders.UserOrders).
		Methods(http.MethodGet).
		Headers("Content-Type", "application/json")
	m.HandleFunc("/api/user/orders/{number}", s.orders.UserOrderByNumber).
		Methods(http.MethodGet).
		Headers("Content-Type", "application/json")

	// Balance handlers
	m.HandleFunc("/api/user/balance", s.balance.GetForUser).
		Methods(http.MethodGet).
		Headers("Content-Type", "application/json")

	// Withdrawals handlers
	m.HandleFunc("/api/user/balance/withdraw", s.withdrawals.Create).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
	m.HandleFunc("/api/user/withdrawals", s.withdrawals.UserWithdrawals).
		Methods(http.MethodGet).
		Headers("Content-Type", "application/json")

	return m
}
