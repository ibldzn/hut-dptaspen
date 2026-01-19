package handler

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ibldzn/spinner-hut/internal/services"
	"github.com/ibldzn/spinner-hut/internal/templates"
)

type Config struct {
	EmpService *services.EmployeeService
}

type Handler struct {
	cfg       Config
	templates *template.Template
	staticFS  fs.FS
}

func NewHandler(cfg Config) (*Handler, error) {
	tpls, err := templates.ParseTemplates()
	if err != nil {
		return nil, err
	}

	staticFS, err := templates.StaticFS()
	if err != nil {
		return nil, err
	}

	return &Handler{
		templates: tpls,
		cfg:       cfg,
		staticFS:  staticFS,
	}, nil
}

func (h *Handler) Into() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	staticServer := http.FileServer(http.FS(h.staticFS))
	r.Handle("/styles.css", staticServer)
	r.Handle("/app.js", staticServer)
	r.Handle("/public/*", staticServer)

	r.Get("/spinner", h.RenderSpinnerPage)

	r.Get("/api/employees", h.GetPresentEmployees)
	r.Post("/api/employees/mark_present", h.MarkEmployeePresent)

	return r
}
