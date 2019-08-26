package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kakugirai/moetify/config"
	"github.com/urfave/negroni"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type App struct {
	Router      *mux.Router
	Middlewares *Middleware
	RS          RedisStorage
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortlinkRes struct {
	Shortlink string `json:"shortlink"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}

func (a *App) Initialize(e config.RedisEnv) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	addr, passwd, db := config.GetRedisEnv().Addr, config.GetRedisEnv().Password, config.GetRedisEnv().DB
	log.Printf("connect to redis (addr: %s, password: %s, db: %d)", addr, passwd, db)
	a.RS = NewRedisCli(addr, passwd, db)
	a.Router = mux.NewRouter()
	a.Middlewares = &Middleware{}
	a.InitializeRoutes()
}

func (a App) InitializeRoutes() {
	a.Router.HandleFunc("/api/shorten", a.createShortlink).Methods("POST")
	a.Router.HandleFunc("/api/info", a.getShortlinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("GET")
}

func (a App) createShortlink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, StatusError{
			http.StatusBadRequest,
			fmt.Errorf("parse parameters failed %v", r.Body),
		}, nil)
		return
	}

	if err := validator.Validate(req); err != nil {
		respondWithError(w, StatusError{
			http.StatusBadRequest,
			fmt.Errorf("validate parameters failed %v", req),
		}, nil)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	s, err := a.RS.Shorten(req.URL, req.ExpirationInMinutes)
	if err != nil {
		respondWithError(w, err, nil)
	} else {
		respondWithJSON(w, http.StatusCreated, shortlinkRes{Shortlink: s})
	}
}

func (a App) getShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")
	d, err := a.RS.ShortLinkInfo(s)
	if err != nil {
		respondWithError(w, err, nil)
	} else {
		respondWithJSON(w, http.StatusOK, d)
	}
}

func (a App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u, err := a.RS.Unshorten(vars["shortlink"])
	if err != nil {
		respondWithError(w, err, nil)
	} else {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

func (a *App) Run(addr string) {
	n := negroni.New()
	//n.Use(negroni.HandlerFunc(a.Middlewares.LoggingHandler))
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(a.Router)
	n.Run(addr)
	//log.Fatal(http.ListenAndServe(addr, a.Router))
}

func respondWithError(w http.ResponseWriter, err error, payload interface{}) {
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s", e.Status(), e)
		resp, _ := json.Marshal(Response{Code: e.Status(), Message: e.Error(), Content: payload})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
	default:
		respondWithJSON(w, http.StatusInternalServerError, payload)
	}
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	resp, _ := json.Marshal(Response{Code: status, Message: http.StatusText(status), Content: payload})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(resp)
}
