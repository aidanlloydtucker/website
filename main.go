package main

import (
	"log"
	"net/http"

	"gopkg.in/boj/redistore.v1"

	"github.com/gorilla/mux"
	"github.com/yosssi/ace"
)

var store *redistore.RediStore

var router *mux.Router

const sessionName = "535510N"

var (
	Port string
)

func main() {
	var err error
	store, err = redistore.NewRediStore(10, "tcp", "127.0.01:6379", "", []byte("sessions"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	router = mux.NewRouter()
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/homework", HomeworkHandler).Methods("GET")
	router.HandleFunc("/homework/aidan", HomeworkAidanHandler).Methods("GET")
	router.HandleFunc("/homework/assignments", HomeworkAssignmentsHandler).Methods("GET")
	router.HandleFunc("/homework/classes", HomeworkPUTClassesHandler).Methods("PUT")
	router.HandleFunc("/homework/classes", HomeworkGETClassesHandler).Methods("GET")
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", AllHandler)

	log.Fatal(http.ListenAndServe(":"+Port, nil))
}

func AllHandler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := ace.Load("views/base", "views/index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
