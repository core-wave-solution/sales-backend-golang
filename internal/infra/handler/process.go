package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process"
	processusecases "github.com/willjrcom/sales-backend-go/internal/usecases/process"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerProcessImpl struct {
	s *processusecases.Service
}

func NewHandlerProcess(processService *processusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerProcessImpl{
		s: processService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterProcess)
		c.Patch("/update/{id}", h.handlerUpdateProcess)
		c.Delete("/{id}", h.handlerDeleteProcess)
		c.Get("/all", h.handlerGetAllProcesses)
	})

	return handler.NewHandler("/process", c)
}

func (h *handlerProcessImpl) handlerRegisterProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoProcess := &processdto.CreateProcessInput{}
	if err := jsonpkg.ParseBody(r, dtoProcess); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.RegisterProcess(ctx, dtoProcess)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerProcessImpl) handlerUpdateProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.UpdateProcess(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessImpl) handlerDeleteProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteProcess(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
func (h *handlerProcessImpl) handlerGetAllProcesses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	processes, err := h.s.GetAllProcesses(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: processes})
}
