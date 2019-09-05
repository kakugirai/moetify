package app

import (
	"log"

	"moetify/app/handler"
	"moetify/app/middleware"
	"moetify/app/model"
	"moetify/config"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// App contains router, middleware and redis
type App struct {
	Router      *mux.Router
	Middlewares *middleware.Middleware
	Handler     *handler.Handler
}

// Initialize app
func (a *App) Initialize(e config.RedisEnv) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	addr, passwd, db := config.GetRedisEnv().Addr, config.GetRedisEnv().Password, config.GetRedisEnv().DB
	log.Printf("connect to redis (addr: %s, password: %s, db: %d)", addr, passwd, db)

	// Let the handler control Redis conn
	a.Handler = &handler.Handler{
		RS: model.NewRedisCli(addr, passwd, db),
	}
	a.Router = mux.NewRouter()
	a.Middlewares = &middleware.Middleware{}
	a.InitializeRoutes()
}

// InitializeRoutes initialize routes
func (a App) InitializeRoutes() {
	a.Router.HandleFunc("/api/shorten", a.Handler.CreateShortlink).Methods("POST")
	a.Router.HandleFunc("/api/info", a.Handler.GetShortlinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}", a.Handler.Redirect).Methods("GET")
}

// Run negroni
func (a *App) Run(addr string) {
	n := negroni.New()
	//n.Use(negroni.HandlerFunc(a.Middlewares.LoggingHandler))
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(a.Router)
	n.Run(addr)
	//log.Fatal(http.ListenAndServe(addr, a.Router))
}
