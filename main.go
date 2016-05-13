package main

import (
	"log"
	"net/http"

	"gopkg.in/boj/redistore.v1"

	"github.com/gorilla/mux"
	"github.com/yosssi/ace"
	"github.com/yosssi/ace-proxy"
)

var p = proxy.New(&ace.Options{BaseDir: "views"})

var store *redistore.RediStore

var router *mux.Router

const sessionName = "535510N"

var (
	HttpPort  string
	HttpsPort string
	CertFile  string
	KeyFile   string
)

func main() {
	var err error
	store, err = redistore.NewRediStore(10, "tcp", "127.0.0.1:6379", "", []byte("sessions"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	if HttpPort == "" {
		HttpPort = "8080"
	}

	router = mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/homework", HomeworkHandler).Methods("GET")
	router.HandleFunc("/homework/aidan", HomeworkAidanHandler).Methods("GET")
	router.HandleFunc("/homework/assignments", HomeworkAssignmentsHandler).Methods("GET")
	router.HandleFunc("/homework/classes", HomeworkPUTClassesHandler).Methods("PUT")
	router.HandleFunc("/homework/classes", HomeworkGETClassesHandler).Methods("GET")
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./static/"))))
	router.PathPrefix("/keybase/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./static/"))))
	//http.HandleFunc("/", AllHandler)

	if CertFile != "" && KeyFile != "" && HttpsPort != "" {
		httpServer := http.NewServeMux()
		httpServer.HandleFunc("/", AllHandler)

		httpsServer := http.NewServeMux()
		httpsServer.HandleFunc("/", AllHandler)

		go func() {
			err := http.ListenAndServe("localhost:8081", httpServer)
			if err != nil {
				log.Fatal(err)
			}
		}()

		log.Fatal(http.ListenAndServeTLS(":"+HttpsPort, CertFile, KeyFile, httpsServer))
	} else {
		log.Fatal(http.ListenAndServe(":"+HttpPort, http.HandlerFunc(AllHandler)))
	}

}

func AllHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	router.ServeHTTP(w, r)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorHandler(w, r, "Not Found: The page you requested could not be found.", 404)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, errStr string, errNum int) {
	tpl, err := p.Load("base", "error", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Status":  errNum,
		"Message": errStr,
	}

	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := p.Load("base", "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
