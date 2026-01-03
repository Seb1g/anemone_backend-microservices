package api

import (
	"anemone_backend-microservices/internal/catechize/model"
	"anemone_backend-microservices/internal/catechize/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Status string `json:"status"`
}

type Handler struct {
	Service *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) Routes(r *mux.Router, JWTsecret string) {
	api := r.PathPrefix("/api/v1/quiz").Subrouter()
	api.Use(AuthMiddleware(JWTsecret))

	api.HandleFunc("/add", h.addResult).Methods("POST")
	api.HandleFunc("/get_all", h.getResults).Methods("GET")
	api.HandleFunc("/clear", h.clearResults).Methods("DELETE")
}

func (h *Handler) addResult(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r.Context())

	var req model.QuizResult
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.Service.AddResult(r.Context(), userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *Handler) getResults(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r.Context())

	results, err := h.Service.GetResults(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func (h *Handler) clearResults(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r.Context())

	if err := h.Service.ClearResultsUser(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondOK(w)
}

func respondOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Response{
		Status: "ok",
	})
}
