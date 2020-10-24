// Package server provides functionality to store tweets and retrieve emojis.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/copilot-ghostpilot-app/api/tweets"

	"github.com/gorilla/mux"
)

// Server is an API server.
type Server struct {
	Router *mux.Router
	TC     tweets.TweetsController
	EC     tweets.EmojisController
}

// ServeHTTP delegates to the mux router.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.HandleFunc("/_healthcheck", s.handleHealthCheck())
	s.Router.HandleFunc("/tweets/create", s.handleStoreTweets()).Methods(http.MethodPost)
	s.Router.HandleFunc("/tweets/emojis", s.handleGetEmojiResults()).Methods(http.MethodGet)

	s.Router.ServeHTTP(w, r)
}

func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) handleStoreTweets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data tweets.Tweet
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&data); err != nil {
			log.Printf("ERROR: server: decode payload: %v\n", err)
			http.Error(w, "decode JSON payload", http.StatusBadRequest)
			return
		}
		if err := s.TC.StoreTweet(data); err != nil {
			log.Printf("ERROR: server: store tweet %+v: %v\n", data, err)
			http.Error(w, fmt.Sprintf("store tweet for username %s and id %s", data.Username, data.ID), http.StatusInternalServerError)
			return
		}
		log.Printf("INFO: server: stored tweet for username %s and id %s\n", data.Username, data.ID)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) handleGetEmojiResults() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := s.EC.EmojiResults()
		if err != nil {
			log.Printf("ERROR: server: get all emoji results: %v\n", err)
			http.Error(w, "get emoji results", http.StatusInternalServerError)
			return
		}

		dat, err := json.Marshal(&struct {
			EmojiResults []tweets.EmojiCount `json:"emojis"`
		}{
			EmojiResults: results,
		})
		log.Printf("results=%v", results)
		if err != nil {
			log.Printf("ERROR: encode get emoji results payload: %v", err)
			http.Error(w, "encode JSON payload", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
	}
}
