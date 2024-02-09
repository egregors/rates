package web

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"net/http"
	"sort"

	srv "github.com/egregors/rates/internal/server"
)

var tmpl = template.Must(template.ParseFiles(
	"internal/server/web/templates/base.gohtml",
	"internal/server/web/templates/index.gohtml",
	"internal/server/web/templates/rates-form-and-history.gohtml",
))

type RateReq struct {
	From, To, Rate string
}

type Currency struct {
	Code, Title string
}

type Data struct {
	Currencies []Currency
	History    []RateReq
}

type Server struct {
	currencies []Currency
	histories  map[string][]RateReq

	rp srv.RateProvider
	r  chi.Router
	l  srv.Logger
}

func New(rp srv.RateProvider, l srv.Logger) *Server {
	s := &Server{
		currencies: nil,
		histories:  make(map[string][]RateReq),
		rp:         rp,
		r:          chi.NewRouter(),
		l:          l,
	}

	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Recoverer)
	s.r.Use(middleware.StripSlashes)

	s.r.Get("/", s.index)
	s.r.Post("/rates", s.getRatesAndHistory)

	s.r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("internal/server/web/static"))))

	s.currencies = s.prepareCurrencies()

	return s
}

func (s *Server) Run() error {
	// TODO: get port from env
	s.l.Printf("[INFO] web server is running on :80")
	return http.ListenAndServe(":80", s.r)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, Data{
		Currencies: s.currencies,
		History:    s.histories[r.RemoteAddr],
	}); err != nil {
		s.l.Printf("failed to render template: %v", err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}

func (s *Server) getRatesAndHistory(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	to := r.FormValue("to")

	rate, err := s.rp.GetRate(from, to)
	if err != nil {
		s.l.Printf("failed to get rate: %v", err)
		// TODO: return error widget instead of 500
		http.Error(w, "failed to get rate", http.StatusInternalServerError)
		return
	}

	s.histories[r.RemoteAddr] = append(s.histories[r.RemoteAddr], RateReq{from, to, fmt.Sprintf("%.2f", rate)})

	if err := tmpl.Execute(w, Data{
		Currencies: s.currencies,
		History:    s.histories[r.RemoteAddr],
	}); err != nil {
		s.l.Printf("failed to render template: %v", err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}

func (s *Server) prepareCurrencies() []Currency {
	cs, err := s.rp.GetCurrencyList()
	if err != nil {
		s.l.Printf("failed to get currency list: %v", err)
		// TODO: don't panic
		panic(err)
	}

	currencies := make([]Currency, 0, len(cs))
	for c, t := range cs {
		currencies = append(currencies, Currency{c, t})
	}

	sort.SliceStable(currencies, func(i, j int) bool {
		return currencies[i].Code < currencies[j].Code
	})

	return currencies
}
