package shortcutter

import (
	"context"
	"everstake-affiliate/api/router"
	"everstake-affiliate/conf"
	"everstake-affiliate/services"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"time"
)

// API serves the end users requests.
type API struct {
	server       *http.Server
	router       *mux.Router
	service      *services.Service
	config       conf.Config
	queryDecoder *schema.Decoder
}

func NewAPI(cfg conf.Config, s *services.Service) *API {
	queryDecoder := schema.NewDecoder()
	queryDecoder.IgnoreUnknownKeys(true)
	api := &API{
		service:      s,
		config:       cfg,
		queryDecoder: queryDecoder,
	}
	api.initialize()
	return api
}

func (api *API) Run() error {
	return api.startServe()
}

func (api *API) Title() string {
	return "Shortcutter"
}

func (api *API) Stop() error {
	cont, cancel := context.WithTimeout(context.Background(), time.Duration(api.config.WaitTimeout)*time.Second)
	defer cancel()
	return api.server.Shutdown(cont)
}

func (api *API) startServe() error {
	log.Infof("Listening on %s", api.server.Addr)
	err := api.server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Warn("API server was closed")
		return nil
	}

	if err != nil {
		log.Errorf("Cannot run API service: %s", err.Error())
	}

	return nil
}

func (api *API) initialize() {
	api.router = mux.NewRouter()

	wrapper := negroni.New()

	wrapper.Use(cors.New(cors.Options{
		AllowedOrigins:   api.config.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-User-Env", "X-Forwarded-For", "Referer"},
	}))

	publicWrapper := negroni.New()

	publicWrapper.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-User-Env", "X-Forwarded-For", "Referer"},
	}))

	//public routes (no middleware)
	router.HandleActions(api.router, publicWrapper, "", []*router.Route{
		{Path: "/", Method: http.MethodGet, Func: api.DefaultRedirect, Middleware: nil},
		{Path: "/favicon.ico", Method: http.MethodGet, Func: http.FileServer(http.Dir("./resources/static")).ServeHTTP, Middleware: nil},
		{Path: "/robots.txt", Method: http.MethodGet, Func: http.FileServer(http.Dir("./resources/static")).ServeHTTP, Middleware: nil},
		//Redirect URL
		//Always should be at the end
		{Path: "/{url_id}", Method: http.MethodGet, Func: api.GetShortCutInfo, Middleware: nil},
		{Path: "/{url_id}/{title:.*}", Method: http.MethodGet, Func: api.GetShortCutInfo, Middleware: nil},
	})

	api.server = &http.Server{Addr: fmt.Sprintf(":%d", api.config.ShortcutterListenOnPort), Handler: api.router}
}
