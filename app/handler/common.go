package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	myerror "moetify/app/error"
	"moetify/app/model"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

// Handler struct should be used as long-living struct by App.
type Handler struct {
	RS model.RedisStorage
}

// Response is HTTP response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=1"`
}

type shortlinkRes struct {
	Shortlink string `json:"shortlink"`
}

// CreateShortlink creates and responds the short link
func (h *Handler) CreateShortlink(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, myerror.StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("parse parameters failed %v", r.Body),
		}, nil)
		return
	}

	if err := validator.Validate(req); err != nil {
		respondWithError(w, myerror.StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("validate parameters failed %v", req),
		}, nil)
		return
	}

	s, err := h.RS.Shorten(req.URL, req.ExpirationInMinutes)
	if err != nil {
		respondWithError(w, nil, err)
	} else {
		respondJSON(w, Response{
			Code:    http.StatusCreated,
			Message: http.StatusText(http.StatusCreated),
			Content: shortlinkRes{Shortlink: s},
		})
	}
}

// GetShortlinkInfo responds the link info
func (h *Handler) GetShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")
	d, err := RS.ShortLinkInfo(s)
	if err != nil {
		respondWithError(w, nil, err)
	} else {
		respondJSON(w, Response{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
			Content: d,
		})
	}
}

// Redirect redirects from shortlink to full link
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u, err := RS.Unshorten(vars["shortlink"])
	if err != nil {
		respondWithError(w, nil, err)
	} else {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

// RespondWithError writes response with error
func respondWithError(w http.ResponseWriter, err myerror.StatusError, payload interface{}) {

	var resp []byte

	switch e := err.(type) {
	case myerror.StatusError:
		log.Printf("HTTP %d - %s\n", e.Status(), e)
		resp, err := json.Marshal(Response{
			Code:    e.Status(),
			Message: e.Error(),
			Content: payload,
		})
		if err != nil {
			log.Fatalln(err)
		}
		break
	default:
		resp, err := json.Marshal(Response{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Content: payload,
		})
		if err != nil {
			log.Fatalln(err)
		}
		break
	}
	respondJSON(w, resp)
}

// RespondJSON writes response with JSON
func respondJSON(w http.ResponseWriter, response Response) {
	resp, err := json.Marshal(response)
	if err != nil {
		log.Fatalln(err)
	}
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
