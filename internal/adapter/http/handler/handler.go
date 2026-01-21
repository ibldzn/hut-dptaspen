package handler

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ibldzn/spinner-hut/internal/services"
	"github.com/ibldzn/spinner-hut/internal/templates"
)

type Config struct {
	EmpService    *services.EmployeeService
	WinnerService *services.WinnerService
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
	r.Use(cors.AllowAll().Handler)

	staticServer := http.FileServer(http.FS(h.staticFS))
	r.Handle("/styles.css", staticServer)
	r.Handle("/invitation.css", staticServer)
	r.Handle("/spinner.js", staticServer)
	r.Handle("/admin.js", staticServer)
	r.Handle("/public/*", staticServer)

	r.Get("/", h.RenderInvitationPage)
	r.Get("/spinner", h.RenderSpinnerPage)
	r.Get("/admin", h.RenderAdminPage)

	validateAPIKey := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey != "Dptaspen@25!" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	apiRouter := chi.NewRouter()
	apiRouter.Group(func(r chi.Router) {
		r.Use(validateAPIKey)
		r.Get("/employees/present", h.GetPresentEmployees)
		r.Get("/employees/all", h.GetAllEmployees)
		r.Get("/employees/export", h.ExportAttendance)
		r.Post("/employees/mark_present", h.MarkEmployeePresent)
		r.Get("/winners", h.GetWinners)
		r.Get("/winners/export", h.ExportWinners)
		r.Post("/winners", h.AddWinners)
	})

	r.Mount("/api", apiRouter)

	return r
}
