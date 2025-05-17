package user

import (
	"encoding/json"
	"errors"
	"github.com/sidgim/finance_shared/meta"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sidgim/finance_shared/httphelper"
	"gorm.io/gorm"
)

// 1. El struct que guarda tu service
type Handler struct {
	svc Service
}

var validate = validator.New()

// 2. Constructor
func NewUserHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// 3. Mount: se “engancha” al router
func (h *Handler) Mount(r chi.Router) {
	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetByID)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
}

// 4. Handlers como métodos del struct

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(req); err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	usr, err := h.svc.Create(r.Context(), req)
	if err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "could not create user")
		return
	}
	httphelper.WriteSuccess(w, http.StatusCreated, usr, nil)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, "invalid UUID")
		return
	}

	usr, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if usr == nil {
		httphelper.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	httphelper.WriteSuccess(w, http.StatusOK, usr, nil)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	filters := Filters{
		FirstName: q.Get("first_name"),
		LastName:  q.Get("last_name"),
	}

	count, err := h.svc.Count(r.Context(), filters)
	if err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "count failed")
		return
	}
	m, err := meta.New(offset, limit, count, "10")
	if err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "meta error")
		return
	}

	list, err := h.svc.GetAll(r.Context(), filters, m.Offset(), m.Limit())
	if err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "fetch failed")
		return
	}
	httphelper.WriteSuccess(w, http.StatusOK, list, m)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, "invalid UUID")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.Struct(req); err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.svc.UpdateContact(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httphelper.WriteError(w, http.StatusNotFound, "user not found")
		} else {
			httphelper.WriteError(w, http.StatusInternalServerError, "update failed")
		}
		return
	}
	httphelper.WriteSuccess(w, http.StatusOK, updated, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		httphelper.WriteError(w, http.StatusBadRequest, "invalid UUID")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		httphelper.WriteError(w, http.StatusInternalServerError, "delete failed")
		return
	}
	httphelper.WriteSuccess(w, http.StatusNoContent, nil, nil)
}
