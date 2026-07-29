package example

import (
	"errors"
	"net/http"

	"go-api-boilerplate/internal/shared"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/examples", h.get)
	mux.HandleFunc("POST /api/v1/examples", h.create)
	mux.HandleFunc("GET /api/v1/examples/{id}", h.getOne)
	mux.HandleFunc("PATCH /api/v1/examples/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/examples/{id}", h.delete)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := shared.DecodeJSON(r, &req); err != nil {
		shared.WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := shared.Validate(req); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.service.Create(r.Context(), req)
	if err != nil {
		handleWriteError(w, err)
		return
	}
	shared.WriteSuccess(w, http.StatusCreated, "example created successfully", item)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	limit, offset, page := shared.Pagination(r)
	items, total, err := h.service.GetAll(r.Context(), shared.Search(r), limit, offset)
	if err != nil {
		shared.WriteInternalError(w, err)
		return
	}
	shared.WritePaginated(w, "successfully retrieved examples", items, total, page, limit)
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	item, err := h.service.GetOne(r.Context(), r.PathValue("id"))
	if err != nil {
		handleWriteError(w, err)
		return
	}
	shared.WriteSuccess(w, http.StatusOK, "successfully retrieved example", item)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var req UpdateRequest
	if err := shared.DecodeJSON(r, &req); err != nil {
		shared.WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := shared.Validate(req); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.service.Update(r.Context(), r.PathValue("id"), req)
	if err != nil {
		handleWriteError(w, err)
		return
	}
	shared.WriteSuccess(w, http.StatusOK, "example updated successfully", item)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), r.PathValue("id")); err != nil {
		handleWriteError(w, err)
		return
	}
	shared.WriteSuccess(w, http.StatusOK, "example deleted successfully", nil)
}

func handleWriteError(w http.ResponseWriter, err error) {
	var validationError shared.ValidationError
	if errors.As(err, &validationError) {
		shared.WriteError(w, http.StatusBadRequest, validationError.Error())
		return
	}
	if errors.Is(err, ErrExampleNotFound) {
		shared.WriteError(w, http.StatusNotFound, ErrExampleNotFound.Error())
		return
	}
	shared.WriteInternalError(w, err)
}
