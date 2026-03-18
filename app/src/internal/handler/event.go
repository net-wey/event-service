package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"event-service/internal/model"
	"event-service/internal/service"
)

// EventHandler обрабатывает HTTP-запросы для мероприятий.
type EventHandler struct {
	svc *service.EventService
}

// NewEventHandler создаёт новый EventHandler.
func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

// List возвращает список всех мероприятий.
// @Summary Список мероприятий
// @Tags events
// @Produce json
// @Success 200 {array} model.Event
// @Router /events [get]
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	events, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if events == nil {
		events = []model.Event{}
	}
	json.NewEncoder(w).Encode(events)
}

// Get возвращает мероприятие по ID.
// @Summary Получить мероприятие
// @Tags events
// @Produce json
// @Param eventID path string true "ID мероприятия"
// @Success 200 {object} model.Event
// @Failure 404 {object} map[string]string
// @Router /events/{eventID} [get]
func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "eventID")

	event, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(event)
}

// Create обрабатывает создание мероприятия.
// @Summary Создать мероприятие
// @Tags events
// @Accept json
// @Produce json
// @Param request body model.CreateEventRequest true "Данные мероприятия"
// @Success 201 {object} model.Event
// @Failure 400 {object} map[string]string
// @Router /events [post]
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.svc.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "title and event_date are required (event_date in RFC3339 format)")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// Update обрабатывает частичное обновление мероприятия.
// @Summary Обновить мероприятие
// @Tags events
// @Accept json
// @Produce json
// @Param eventID path string true "ID мероприятия"
// @Param request body model.UpdateEventRequest true "Поля для обновления"
// @Success 200 {object} model.Event
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /events/{eventID} [put]
func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "eventID")

	var req model.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "invalid input")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(event)
}

// Delete обрабатывает удаление мероприятия.
// @Summary Удалить мероприятие
// @Tags events
// @Param eventID path string true "ID мероприятия"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /events/{eventID} [delete]
func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "eventID")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
