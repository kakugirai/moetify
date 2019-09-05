package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	myerror "github.com/kakugirai/moetify/app/error"
	"github.com/kakugirai/moetify/app/model"

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
		//log.Println("Bad decode")
		respondWithError(w, myerror.StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("parse parameters failed %v", r.Body),
		}, nil)
		return
	}

	if err := validator.Validate(req); err != nil {
		//log.Println("Bad validation")
		respondWithError(w, myerror.StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("validate parameters failed %v", req),
		}, nil)
		return
	}

	s, err := h.RS.Shorten(req.URL, req.ExpirationInMinutes)
	if err != nil {
		//log.Println("Bad shorten")
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
	defer r.Body.Close()

	vals := r.URL.Query()
	s := vals.Get("shortlink")

	//log.Println("shortlink is ", s)

	d, err := h.RS.ShortLinkInfo(s)
	if err != nil {
		respondWithError(w, err, nil)
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
	defer r.Body.Close()

	vars := mux.Vars(r)

	d, err := h.RS.ShortLinkInfo(vars["shortlink"])
	if err != nil {
		respondWithError(w, err, nil)
	} else {
		u, err := url.Parse(d.Full)
		if err != nil {
			log.Fatalln(err)
		}
		// Make sure the redirect URL is absolute instead of relative
		// FIXME: should not hard coding http
		u.Scheme = "http"
		http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	}
}

// RespondWithError writes response with error
func respondWithError(w http.ResponseWriter, err error, payload interface{}) {
	var resp Response

	switch e := err.(type) {
	case myerror.Error:
		log.Printf("HTTP %d - %s\n", e.Status(), e)
		resp = Response{
			Code:    e.Status(),
			Message: e.Error(),
			Content: payload,
		}
		break
	default:
		resp = Response{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Content: payload,
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	_, _ = w.Write(resp)
}
