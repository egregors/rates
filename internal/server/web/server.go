package web

import (
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	srv "github.com/egregors/rates/internal/server"
)

var tmpl = template.Must(template.ParseFiles(
	"internal/server/web/templates/base.gohtml",
	"internal/server/web/templates/index.gohtml",
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

	s.r.Handle("/static/*", http.FileServer(http.FS(static)))

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
	prompt := r.FormValue("prompt")
	from, to, amount, err := parsePrompt(prompt)
	if err != nil {
		s.l.Printf("failed to parse prompt: %v", err)
		http.Error(w, "invalid prompt", http.StatusBadRequest)
		return
	}

	rate, err := s.rp.GetRate(from, to)
	if err != nil {
		s.l.Printf("failed to get rate: %v", err)
		// TODO: return error widget instead of 500
		http.Error(w, "failed to get rate", http.StatusInternalServerError)
		return
	}

	s.histories[r.RemoteAddr] = append(
		s.histories[r.RemoteAddr],
		RateReq{
			from,
			to,
			fmt.Sprintf("%.2f %s -> %.2f %s", amount, from, amount*rate, to),
		},
	)

	if err := tmpl.Execute(w, Data{
		Currencies: s.currencies,
		History:    s.histories[r.RemoteAddr],
	}); err != nil {
		s.l.Printf("failed to render template: %v", err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
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
