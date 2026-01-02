package api

import (
	"anemone_backend-microservices/internal/auth/services"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
}

func (h *AuthHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/api/v1/auth/register", http.HandlerFunc(h.register))
	mux.Handle("/api/v1/auth/login", http.HandlerFunc(h.login))
	mux.Handle("/api/v1/auth/change-password", http.HandlerFunc(h.changePassword))
	mux.Handle("/api/v1/auth/refresh", http.HandlerFunc(h.refresh))
	mux.Handle("/api/v1/auth/logout", http.HandlerFunc(h.logout))
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, user, err := h.Service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_data":     user,
	})
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	access, refresh, user, err := h.Service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]any{
		"access_token":  access,
		"refresh_token": refresh,
		"user_data":     user,
	})
}

func (h *AuthHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Service.ChangePassword(r.Context(), req.Email, req.OldPassword, req.NewPassword)
	if err != nil {
		if strings.Contains(err.Error(), "invalid old password") || strings.Contains(err.Error(), "user not found") {
			http.Error(w, "Invalid email or old password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Password change failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newAccess, newRefresh, user, err := h.Service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]any{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
		"user_data":     user,
	})
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.Logout(r.Context(), req.RefreshToken); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
