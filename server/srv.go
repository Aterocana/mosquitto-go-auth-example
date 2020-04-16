package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//go:generate stringer -type=acc
type acc uint8

const (
	_             = iota
	read      acc = 1
	write     acc = 2
	readWrite acc = 3
	subscribe acc = 4
)

type Server struct {
	http.Server
}

func New() *Server {
	srv := Server{}
	router := mux.NewRouter()
	// /user is the enpoint which has to check for user credentials
	router.Handle("/auth", srv.auth(srv.authUser)).Methods(http.MethodPost)

	// /admin is the enpoint which has to check for admin user credentials
	router.Handle("/admin_auth", srv.auth(srv.authSuperUser)).Methods(http.MethodPost)

	// /acl is the enpoint which has to check if a user can perform an operation over a topic.
	router.Handle("/acl", srv.acl()).Methods(http.MethodPost)
	srv.Server = http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	return &srv
}

func printStatus(status int, r *http.Request) {
	log.Printf("[%d] %s %s\n------------\n", status, r.Method, r.RequestURI)
}

func (srv *Server) auth(auth func(string, string) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		log.Printf("REQUEST: %s\n", r.RequestURI)
		decoder := json.NewDecoder(r.Body)
		var res map[string]string
		if err := decoder.Decode(&res); err != nil {
			printStatus(http.StatusBadRequest, r)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(res)
		username, ok := res["username"]
		if !ok {
			printStatus(http.StatusBadRequest, r)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		password, ok := res["password"] // not checking if exists for admin purposes
		log.Printf("auth with <%s, %s>: %v\n", username, password, auth(username, password))
		if !auth(username, password) {
			w.WriteHeader(http.StatusUnauthorized)
			printStatus(http.StatusUnauthorized, r)
			return
		}
		printStatus(http.StatusOK, r)
	}
}

func (srv *Server) authUser(username, password string) bool {
	return username == "test" && password == "test" || srv.authSuperUser(username, password)
}

func (srv *Server) authSuperUser(username, password string) bool {
	return username == "admin" && (password == "admin" || password == "")
}

func (srv *Server) acl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		var res struct {
			Username string `json:"username"`
			ClientID string `json:"clientid"`
			Topic    string `json:"topic"`
			ACC      acc    `json:"acc"`
		}
		if err := decoder.Decode(&res); err != nil {
			printStatus(http.StatusBadRequest, r)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(res)
		authorized := false
		switch res.ACC {
		case read:
			authorized = srv.canPub(res.Username, res.Topic)
		case write:
			authorized = srv.canPub(res.Username, res.Topic)
		case readWrite:
			authorized = srv.canSub(res.Username, res.Topic) && srv.canPub(res.Username, res.Topic)
		case subscribe:
			authorized = srv.canSub(res.Username, res.Topic)
		}
		if !authorized {
			w.WriteHeader(http.StatusUnauthorized)
			printStatus(http.StatusUnauthorized, r)
			return
		}
		printStatus(http.StatusOK, r)
	}
}

func (srv *Server) canPub(username, topic string) bool {
	return username == "test" && topic == "topic/pub"
}

func (srv *Server) canSub(username, topic string) bool {
	return username == "test" && topic == "topic/sub"
}
