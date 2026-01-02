package api

import (
	"encoding/json"
	"net/http"

	"anemone_backend-microservices/internal/mail/repository"
	"anemone_backend-microservices/internal/mail/services"

	"github.com/gorilla/mux"
)

type MailHandler struct {
	service *services.MailService
	repo    *repository.Repository
}

func NewMailHandler(s *services.MailService, r *repository.Repository) *MailHandler {
	return &MailHandler{service: s, repo: r}
}

func (h *MailHandler) Register(r *mux.Router, jwtSecret string) {
	api := r.PathPrefix("/api/v1/mail").Subrouter()
	api.Use(Auth(jwtSecret))

	api.HandleFunc("/addresses", h.create).Methods("POST")
	api.HandleFunc("/addresses", h.list).Methods("GET")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(CheckAddressOwner(h.repo))

	protected.HandleFunc("/inbox/{id}", h.inbox).Methods("GET")
	protected.HandleFunc("/addresses/{id}", h.delete).Methods("DELETE")
}

func (h *MailHandler) create(w http.ResponseWriter, r *http.Request) {
	userID, _ := GetUserID(r.Context())

	var req struct{ Address string }
	_ = json.NewDecoder(r.Body).Decode(&req)

	addr, err := h.service.CreateAddress(userID, req.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(addr)
}

func (h *MailHandler) list(w http.ResponseWriter, r *http.Request) {
	userID, _ := GetUserID(r.Context())
	res, _ := h.service.ListAddresses(userID)
	json.NewEncoder(w).Encode(res)
}

func (h *MailHandler) inbox(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(AddressIDKey).(int)
	res, _ := h.service.GetInbox(id)
	json.NewEncoder(w).Encode(res)
}

func (h *MailHandler) delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := GetUserID(r.Context())
	id := r.Context().Value(AddressIDKey).(int)
	_ = h.service.DeleteAddress(id, userID)
	w.WriteHeader(http.StatusNoContent)
}
