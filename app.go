package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type app struct {
	zone    string
	spanner *spanner.Client
}

func (ap *app) handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(logRequest())

	r.Get("/", ap.root)
	r.Get("/ping", ap.ping)
	r.Get("/liveness", ap.liveness)
	r.Get("/readiness", ap.readiness)
	r.Route("/api", func(r chi.Router) {
		r.Get("/messages", ap.getMessages)
		r.Post("/messages", ap.createMessage)
	})

	return r
}

func (ap *app) root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (ap *app) ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func (ap *app) liveness(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (ap *app) readiness(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

type getMessagesResponse struct {
	ServerZone string    `json:"server_zone"`
	Messages   []message `json:"messages"`
}

var readStmt = spanner.NewStatement("SELECT * FROM Messages ORDER BY CreatedAt DESC LIMIT 100")

func (ap *app) getMessages(w http.ResponseWriter, r *http.Request) {
	iter := ap.spanner.Single().Query(r.Context(), readStmt)

	resp := &getMessagesResponse{
		ServerZone: ap.zone,
		Messages:   []message{},
	}

	err := iter.Do(func(row *spanner.Row) error {
		m := message{}

		if err := row.Columns(&m.ID, &m.CreatedAt, &m.Name, &m.Body, &m.WrittenAt); err != nil {
			return err
		}

		resp.Messages = append(resp.Messages, m)

		return nil
	})

	if err != nil {
		log.Err(err).Msg("failed to read messages")
		respondError(w, http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

type createMessageRequest struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type createMessageResponse struct {
	ServerZone string   `json:"server_zone"`
	Message    *message `json:"message"`
}

func (ap *app) createMessage(w http.ResponseWriter, r *http.Request) {
	req := &createMessageRequest{}

	if err := decodeJSONBody(r, req); err != nil {
		respondError(w, http.StatusBadRequest)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		log.Err(err).Msg("failed uuid.NewRandom()")
		respondError(w, http.StatusInternalServerError)
		return
	}

	msg := &message{
		ID:        id.String(),
		CreatedAt: time.Now(),
		Name:      req.Name,
		Body:      req.Body,
		WrittenAt: ap.zone,
	}

	m := []*spanner.Mutation{
		spanner.Insert("Messages",
			[]string{"MessageId", "CreatedAt", "Name", "Body", "WrittenAt"},
			[]interface{}{msg.ID, msg.CreatedAt, msg.Name, msg.Body, msg.WrittenAt},
		),
	}

	_, err = ap.spanner.Apply(r.Context(), m)
	if err != nil {
		log.Err(err).Msg("failed to apply")
		respondError(w, http.StatusInternalServerError)
		return
	}

	resp := &createMessageResponse{
		ServerZone: ap.zone,
		Message:    msg,
	}

	respondJSON(w, http.StatusCreated, resp)
}

func decodeJSONBody(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}

func respondJSON(w http.ResponseWriter, status int, body interface{}) {
	bytes, err := json.Marshal(body)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(bytes); err != nil {
		log.Error().Err(err).Msg(err.Error())
	}
}

type errorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func respondError(w http.ResponseWriter, status int) {
	er := &errorResponse{
		Status:  status,
		Message: http.StatusText(status),
	}

	respondJSON(w, status, er)
}
