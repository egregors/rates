package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/egregors/rates"
	lib "github.com/egregors/rates/pkg/http"
)

type Server struct {
	c rates.Converter
	r chi.Router
	l rates.Logger
}

func New(conv rates.Converter, l rates.Logger) *Server {
	s := &Server{
		// TODO: add conv pool
		c: conv,
		r: chi.NewRouter(),
		l: l,
	}

	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Recoverer)
	s.r.Use(middleware.StripSlashes)

	s.r.Post("/api/v0/convert", s.convert)

	return s
}

// convert converts the amount from one currency to another and writes the result to the response like JSON
func (s *Server) convert(w http.ResponseWriter, r *http.Request) {
	var req struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount,string"`
	}

	if err := lib.DecodeJSON(r.Body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.c.Conv(req.Amount, req.From, req.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lib.RespJSON(w, http.StatusOK, struct {
		Result float64 `json:"result"`
	}{
		Result: result,
	})
}

func (s *Server) Run() error {
	// TODO: get port from env
	s.l.Printf("[INFO] api server is running on :8080")
	return http.ListenAndServe(":8080", s.r)
}
