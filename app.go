package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortlinkRes struct {
	Shortlink string `json:"shortlink"`
}

func (a *App) Initialize() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
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
		})
		return
	}

	if err := validator.Validate(req); err != nil {
		respondWithError(w, StatusError{
			http.StatusBadRequest,
			fmt.Errorf("validate parameters failed %v", req),
		})
		return
	}
	defer r.Body.Close()

	fmt.Printf("%v", req)
}

func (a App) getShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")
	fmt.Printf("%s\n", s)
}

func (a App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shortlink"])
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func respondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s", e.Status(), e)
		respondWithJSON(w, e.Status(), e.Error())
	default:
		respondWithJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}