package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"event-service/internal/model"
	"event-service/internal/service"
)

// ParticipantHandler обрабатывает HTTP-запросы для участников.
type ParticipantHandler struct {
	svc *service.ParticipantService
}

// NewParticipantHandler создаёт новый ParticipantHandler.
func NewParticipantHandler(svc *service.ParticipantService) *ParticipantHandler {
	return &ParticipantHandler{svc: svc}
}

// ListByEvent возвращает всех участников мероприятия.
// @Summary Список участников мероприятия
// @Tags participants
// @Produce json
// @Param eventID path string true "ID мероприятия"
// @Success 200 {array} model.Participant
// @Failure 404 {object} map[string]string
// @Router /events/{eventID}/participants [get]
func (h *ParticipantHandler) ListByEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventID")

	participants, err := h.svc.ListByEvent(r.Context(), eventID)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if participants == nil {
		participants = []model.Participant{}
	}
	json.NewEncoder(w).Encode(participants)
}

// Register добавляет нового участника на мероприятие.
// @Summary Зарегистрировать участника
// @Tags participants
// @Accept json
// @Produce json
// @Param eventID path string true "ID мероприятия"
// @Param request body model.CreateParticipantRequest true "Данные участника"
// @Success 201 {object} model.Participant
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /events/{eventID}/participants [post]
func (h *ParticipantHandler) Register(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventID")

	var req model.CreateParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	participant, err := h.svc.Register(r.Context(), eventID, req)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "name and email are required")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(participant)
}

// Update обрабатывает частичное обновление участника.
// @Summary Обновить участника
// @Tags participants
// @Accept json
// @Produce json
// @Param eventID path string true "ID мероприятия"
// @Param participantID path string true "ID участника"
// @Param request body model.UpdateParticipantRequest true "Поля для обновления"
// @Success 200 {object} model.Participant
// @Failure 404 {object} map[string]string
// @Router /events/{eventID}/participants/{participantID} [put]
func (h *ParticipantHandler) Update(w http.ResponseWriter, r *http.Request) {
	participantID := chi.URLParam(r, "participantID")

	var req model.UpdateParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	participant, err := h.svc.Update(r.Context(), participantID, req)
	if err != nil {
		if errors.Is(err, service.ErrParticipantNotFound) {
			writeError(w, http.StatusNotFound, "participant not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(participant)
}

// Delete обрабатывает удаление участника.
// @Summary Удалить участника
// @Tags participants
// @Param eventID path string true "ID мероприятия"
// @Param participantID path string true "ID участника"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /events/{eventID}/participants/{participantID} [delete]
func (h *ParticipantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	participantID := chi.URLParam(r, "participantID")

	if err := h.svc.Delete(r.Context(), participantID); err != nil {
		if errors.Is(err, service.ErrParticipantNotFound) {
			writeError(w, http.StatusNotFound, "participant not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
