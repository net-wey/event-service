package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"event-service/internal/model"
	"event-service/internal/service"
)

// VenueHandler обрабатывает HTTP-запросы для площадок.
type VenueHandler struct {
	svc *service.VenueService
}

// NewVenueHandler создаёт новый VenueHandler.
func NewVenueHandler(svc *service.VenueService) *VenueHandler {
	return &VenueHandler{svc: svc}
}

// List возвращает список всех площадок.
// @Summary Список площадок
// @Tags venues
// @Produce json
// @Success 200 {array} model.Venue
// @Router /venues [get]
func (h *VenueHandler) List(w http.ResponseWriter, r *http.Request) {
	venues, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if venues == nil {
		venues = []model.Venue{}
	}
	writeJSON(w, http.StatusOK, venues)
}

// Get возвращает площадку по ID.
// @Summary Получить площадку
// @Tags venues
// @Produce json
// @Param venueID path string true "ID площадки"
// @Success 200 {object} model.Venue
// @Failure 404 {object} map[string]string
// @Router /venues/{venueID} [get]
func (h *VenueHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "venueID")

	venue, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrVenueNotFound) {
			writeError(w, http.StatusNotFound, "площадка не найдена")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, venue)
}

// Create обрабатывает создание площадки.
// @Summary Создать площадку
// @Tags venues
// @Accept json
// @Produce json
// @Param request body model.CreateVenueRequest true "Данные площадки"
// @Success 201 {object} model.Venue
// @Failure 400 {object} map[string]string
// @Router /venues [post]
func (h *VenueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	venue, err := h.svc.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "название площадки обязательно")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, venue)
}

// Update обрабатывает частичное обновление площадки.
// @Summary Обновить площадку
// @Tags venues
// @Accept json
// @Produce json
// @Param venueID path string true "ID площадки"
// @Param request body model.UpdateVenueRequest true "Поля для обновления"
// @Success 200 {object} model.Venue
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /venues/{venueID} [put]
func (h *VenueHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "venueID")

	var req model.UpdateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	venue, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrVenueNotFound) {
			writeError(w, http.StatusNotFound, "площадка не найдена")
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "некорректные данные")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, venue)
}

// Delete обрабатывает удаление площадки.
// @Summary Удалить площадку
// @Tags venues
// @Param venueID path string true "ID площадки"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /venues/{venueID} [delete]
func (h *VenueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "venueID")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrVenueNotFound) {
			writeError(w, http.StatusNotFound, "площадка не найдена")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
