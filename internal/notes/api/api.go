package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"anemone_backend-microservices/internal/notes/repository"
	"anemone_backend-microservices/internal/notes/services"

	"github.com/gorilla/mux"
)

type response struct {
	Status string `json:"status"`
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func respondOK(w http.ResponseWriter) {
	respondJSON(w, http.StatusOK, response{Status: "ok"})
}

func userIDFromCtx(ctx context.Context) int64 {
	return UserID(ctx)
}

func mustID(r *http.Request) int {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic("invalid id")
	}
	return id
}

type Handler struct {
	Service *services.Service
	Repo    *repository.Repo
}

func NewHandler(s *services.Service, r *repository.Repo) *Handler {
	return &Handler{Service: s, Repo: r}
}

func (h *Handler) Routes(r *mux.Router, jwtSecret string) {
	notes := r.PathPrefix("/api/v1/notes").Subrouter()
	notes.Use(AuthMiddleware(jwtSecret))
	notes.Use(NoteOwner(h.Repo))

	notes.HandleFunc("", h.createNote).Methods("POST")
	notes.HandleFunc("", h.listNotes).Methods("GET")
	notes.HandleFunc("/{id:[0-9]+}", h.getNote).Methods("GET")
	notes.HandleFunc("/{id:[0-9]+}/title", h.updateTitle).Methods("PUT")
	notes.HandleFunc("/{id:[0-9]+}/content", h.updateContent).Methods("PUT")

	deleteNotes := r.PathPrefix("/api/v1/notes/delete").Subrouter()
	deleteNotes.Use(AuthMiddleware(jwtSecret))
	deleteNotes.Use(NoteOwner(h.Repo))

	deleteNotes.HandleFunc("/{id}/soft-delete", h.markDeletedNote).Methods("PUT")
	deleteNotes.HandleFunc("/{id}/soft-undelete", h.unmarkDeletedNote).Methods("PUT")
	deleteNotes.HandleFunc("/trash/clear/{id}", h.deleteAllMarkedNotes).Methods("DELETE")

	apiFolder := r.PathPrefix("/api/v1/folder").Subrouter()
	apiFolder.Use(AuthMiddleware(jwtSecret))
	apiFolder.Use(FolderOwner(h.Repo))

	apiFolder.HandleFunc("/create", h.createFolder).Methods("POST")
	apiFolder.HandleFunc("/{id}", h.getAllFolders).Methods("GET")
	apiFolder.HandleFunc("/update", h.updateTitleFolder).Methods("PUT")
	apiFolder.HandleFunc("/delete/{id}", h.deleteFolder).Methods("DELETE")

	noteToFolderApi := r.PathPrefix("/api/v1/folder").Subrouter()
	noteToFolderApi.Use(AuthMiddleware(jwtSecret))
	noteToFolderApi.Use(FolderOwner(h.Repo))
	noteToFolderApi.Use(NoteOwner(h.Repo))

	apiFolder.HandleFunc("/get_notes/{id}", h.getAllNotesFromFolder).Methods("GET")
	apiFolder.HandleFunc("/add_to_folder", h.addNoteToFolder).Methods("POST")
	apiFolder.HandleFunc("/cancel_from_folder/{id}", h.removeNoteFromFolder).Methods("POST")
}

/* ===================== NOTES ===================== */

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	userID := userIDFromCtx(r.Context())

	note, err := h.Service.CreateNote(r.Context(), userID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) {
	id := mustID(r)

	userID := userIDFromCtx(r.Context())
	note, err := h.Service.GetNote(r.Context(), id, userID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func (h *Handler) listNotes(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r.Context())

	notes, err := h.Service.ListNotes(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notes)
}

func (h *Handler) updateTitle(w http.ResponseWriter, r *http.Request) {
	id := mustID(r)
	userID := userIDFromCtx(r.Context())

	var req struct {
		Title string `json:"title"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	note, err := h.Service.UpdateTitle(r.Context(), id, userID, req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func (h *Handler) updateContent(w http.ResponseWriter, r *http.Request) {
	id := mustID(r)
	userID := userIDFromCtx(r.Context())

	var req struct {
		Content string `json:"content"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	note, err := h.Service.UpdateContent(r.Context(), id, userID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(note)
}

/* ===================== DELETE ===================== */

func (h *Handler) markDeletedNote(w http.ResponseWriter, r *http.Request) {
	id := mustID(r)
	userID := userIDFromCtx(r.Context())

	if err := h.Service.MarkDeletedNote(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	respondOK(w)
}

func (h *Handler) unmarkDeletedNote(w http.ResponseWriter, r *http.Request) {
	id := mustID(r)
	userID := userIDFromCtx(r.Context())

	if err := h.Service.UnmarkDeletedNote(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	respondOK(w)
}

func (h *Handler) deleteAllMarkedNotes(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r.Context())

	if err := h.Service.DeleteAllMarkedNotes(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondOK(w)
}

/* ===================== FOLDERS ===================== */

func (h *Handler) createFolder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := userIDFromCtx(r.Context())

	folder, err := h.Service.CreateFolder(r.Context(), userID, req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, folder)
}

func (h *Handler) getAllFolders(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r.Context())

	folders, err := h.Service.GetAllFolders(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, folders)
}

func (h *Handler) getAllNotesFromFolder(w http.ResponseWriter, r *http.Request) {
	folderID := mustID(r)
	userID := userIDFromCtx(r.Context())

	notes, err := h.Service.GetAllNotesFromFolder(r.Context(), folderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, notes)
}

func (h *Handler) addNoteToFolder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NoteID   int `json:"note_id"`
		FolderID int `json:"folder_id"`
	}

	userID := userIDFromCtx(r.Context())

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note, err := h.Service.AddNoteToFolder(r.Context(), req.NoteID, req.FolderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	respondJSON(w, http.StatusOK, note)
}

func (h *Handler) removeNoteFromFolder(w http.ResponseWriter, r *http.Request) {
	noteID := mustID(r)
	userID := userIDFromCtx(r.Context())

	note, err := h.Service.RemoveNoteFromFolder(r.Context(), noteID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	respondJSON(w, http.StatusOK, note)
}


func (h *Handler) updateTitleFolder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NewTitle string `json:"new_title"`
	}

	userID := userIDFromCtx(r.Context())

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	folderID := mustID(r)

	folder, err := h.Service.UpdateTitleFolder(r.Context(), folderID, userID, req.NewTitle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	respondJSON(w, http.StatusOK, folder)
}


func (h *Handler) deleteFolder(w http.ResponseWriter, r *http.Request) {
	folderID := mustID(r)
	userID := userIDFromCtx(r.Context())

	if err := h.Service.DeleteFolder(r.Context(), folderID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	respondOK(w)
}

