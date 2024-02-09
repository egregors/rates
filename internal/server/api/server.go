package api

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	srv "github.com/egregors/rates/internal/server"
	lib "github.com/egregors/rates/lib/http"
)

type Server struct {
	rp srv.RateProvider
	r  chi.Router
	l  srv.Logger
}

func New(rp srv.RateProvider, l srv.Logger) *Server {
	s := &Server{
		rp: rp,
		r:  chi.NewRouter(),
		l:  l,
	}
	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Recoverer)
	s.r.Use(middleware.StripSlashes)

	s.r.Get("/currency", s.getCurrencyList)
	s.r.Get("/rate/{from}/{to}", s.getRate)

	return s
}

// getRate gets the rate between two currencies from provider and writes it to the response like JSON
func (s *Server) getRate(w http.ResponseWriter, r *http.Request) {
	from := chi.URLParam(r, "from")
	to := chi.URLParam(r, "to")

	rate, err := s.rp.GetRate(from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lib.RespJSON(w, http.StatusOK, struct {
		From string  `json:"from"`
		To   string  `json:"to"`
		Rate float64 `json:"rate"`
	}{
		From: from,
		To:   to,
		Rate: rate,
	})
}

func (s *Server) getCurrencyList(w http.ResponseWriter, _ *http.Request) {
	cs, err := s.rp.GetCurrencyList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lib.RespJSON(w, http.StatusOK, struct {
		Currencies map[string]string `json:"currencies"`
	}{
		Currencies: cs,
	})
}

func (s *Server) Run() error {
	// TODO: get port from env
	s.l.Printf("[INFO] server is running on :8080")
	return http.ListenAndServe(":8080", s.r)
}
