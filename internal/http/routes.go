package http

import (
	"gophermart/internal/middlewares"
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
		Methods(http.MethodPost)
	m.HandleFunc("/api/user/login", s.auth.Login).
		Methods(http.MethodPost)

	userAuth := m.PathPrefix("/api/user").Subrouter()

	//  Orders handlers
	userAuth.HandleFunc("/orders", s.orders.Create).
		Methods(http.MethodPost)
	userAuth.HandleFunc("/orders", s.orders.UserOrders).
		Methods(http.MethodGet)
	userAuth.HandleFunc("/orders/{number}", s.orders.UserOrderByNumber).
		Methods(http.MethodGet)

	// Balance handlers
	userAuth.HandleFunc("/balance", s.balance.GetForUser).
		Methods(http.MethodGet)

	// Withdrawals handlers
	userAuth.HandleFunc("/balance/withdraw", s.withdrawals.Create).
		Methods(http.MethodPost)
	userAuth.HandleFunc("/withdrawals", s.withdrawals.UserWithdrawals).
		Methods(http.MethodGet)

	// Middlewares
	userAuth.Use(middlewares.AuthorizationMiddleware)
	m.Use(middlewares.GzipMiddleware)

	return m
}
