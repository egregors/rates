package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/egregors/rates"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const historySize = 5

var tmpl = template.Must(template.ParseFiles(
	"internal/server/web/templates/base.gohtml",
	"internal/server/web/templates/rates-form-and-history.gohtml",
))

//go:embed static
var static embed.FS

type RateReq struct {
	From, To, Rate string
}

type Currency struct {
	Code, Title string
}

type Data struct {
	History []RateReq
	Error   string
}

type Server struct {
	histories map[string][]RateReq

	c rates.Converter
	r chi.Router
	l rates.Logger
}

func New(conv rates.Converter, l rates.Logger) *Server {
	s := &Server{
		histories: make(map[string][]RateReq),
		c:         conv,
		r:         chi.NewRouter(),
		l:         l,
	}

	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Recoverer)
	s.r.Use(middleware.StripSlashes)

	s.r.Get("/", s.index)
	s.r.Post("/rates", s.getRatesAndHistory)

	s.r.Handle("/static/*", http.FileServer(http.FS(static)))

	return s
}

func (s *Server) Run() error {
	// TODO: get port from env
	s.l.Printf("[INFO] web server is running on :80")
	return http.ListenAndServe(":80", s.r)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	history := s.histories[r.RemoteAddr]
	if len(history) > historySize {
		history = history[:historySize]
	}

	if err := tmpl.Execute(w, Data{
		History: history,
	}); err != nil {
		s.l.Printf("failed to render template: %v", err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}

func (s *Server) getRatesAndHistory(w http.ResponseWriter, r *http.Request) {
	var errMsg string

	defer func() {
		history := s.histories[r.RemoteAddr]
		if len(history) > historySize {
			history = history[:historySize]
		}

		if err := tmpl.Execute(w, Data{
			History: history,
			Error:   errMsg,
		}); err != nil {
			s.l.Printf("failed to render template: %v", err)
			http.Error(w, "failed to render template", http.StatusInternalServerError)
		}
	}()

	prompt := r.FormValue("prompt")
	from, to, amount, err := parsePrompt(prompt)
	if err != nil {
		s.l.Printf("failed to parse prompt: %v", err)
		errMsg = fmt.Sprintf("invalid prompt: %s", prompt)

		return
	}

	res, err := s.c.Conv(amount, from, to)
	if err != nil {
		s.l.Printf("failed to convert: %v", err)
		errMsg = fmt.Sprintf("failed to convert: %v", err)

		return
	}

	s.pushToHistory(r.RemoteAddr, from, to, fmt.Sprintf("%.2f %s is %.2f %s", amount, from, res, to))
}

func (s *Server) pushToHistory(key, from, to, text string) {
	s.histories[key] = append(
		[]RateReq{
			{
				from,
				to,
				text,
			},
		},
		s.histories[key]...,
	)
}

func parsePrompt(prompt string) (from, to string, amount float64, err error) {
	xs := strings.Split(prompt, " ")
	if len(xs) != 4 {
		return "", "", 0, fmt.Errorf("invalid prompt: %s", prompt)
	}

	amount, err = strconv.ParseFloat(xs[0], 64)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid amount: %s", xs[0])
	}

	return strings.ToLower(xs[1]), strings.ToLower(xs[3]), amount, nil
}
