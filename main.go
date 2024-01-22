package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var systemInfoKeys = []string{"MY_APP_VERSION", "MY_POD_NAME", "MY_POD_IP", "MY_POD_SERVICE_ACCOUNT", "MY_POD_NAMESPACE", "MY_NODE_NAME"}
var startupTime time.Time

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SystemInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var users []User = []User{
	{
		Name:  "Work Shop",
		Email: "workshop@dynatrace.com",
	},
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)

	users = append(users, newUser)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func systemInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	envs := []SystemInfo{
		{
			Name:  "STARTUP_TIME",
			Value: startupTime.Format("2006-01-02T15:04:05Z07:00"),
		},
	}
	for _, v := range systemInfoKeys {
		envs = append(envs, SystemInfo{Name: v, Value: os.Getenv(v)})
	}
	json.NewEncoder(w).Encode(envs)
}
func oom(w http.ResponseWriter, r *http.Request) {
	// Allocate a large amount of memory to trigger an OOM error
	data := make([]byte, 10000000000) // 10 GB
	_ = data                          // Avoid "unused variable" error

	// This line will never be reached due to the OOM error
	log.Println("Allocated a large amount of memory")
}

func main() {
	startupTime = time.Now()
	// Create a new router
	r := mux.NewRouter()
	// Define the endpoints
	r.HandleFunc("/health", health).Methods("GET")
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", addUser).Methods("POST")
	r.HandleFunc("/oom", oom).Methods("POST")
	r.HandleFunc("/system/info", systemInfo).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// request handlers
	r.Use(loggingHandler)
	r.Use(slowDownRequestHandler)
	r.Use(fakingHttpStatusCodeHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", r))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(200)
	w.Write([]byte("Everything ok"))
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func fakingHttpStatusCodeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentPodName := r.URL.Query().Get("pod-name")
		myPodName := getMyPodName()
		if currentPodName == "" || currentPodName == myPodName {
			httpStatusCode := r.URL.Query().Get("http-status-code")
			if httpStatusCode != "" {
				httpStatusCodeInt, err := strconv.Atoi(httpStatusCode)
				if err != nil {
					w.Header().Set("Content-Type", "application/text")
					w.WriteHeader(500)
					w.Write([]byte("could parse query param http-status-code"))
					return
				}
				w.Header().Set("Content-Type", "application/text")
				w.WriteHeader(httpStatusCodeInt)
				w.Write([]byte(fmt.Sprintf("You requested the HTTP-Status-Code %d, here you go...enjoy it :)", httpStatusCodeInt)))
				return
			}
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func slowDownRequestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentPodName := r.URL.Query().Get("pod-name")
		myPodName := getMyPodName()
		if currentPodName == "" || currentPodName == myPodName {
			sleepInterval := r.URL.Query().Get("slow-down")
			if sleepInterval != "" {
				sleepIntervalInt, err := strconv.Atoi(sleepInterval)
				if err != nil {
					w.Header().Set("Content-Type", "application/text")
					w.WriteHeader(500)
					w.Write([]byte("could parse query param slow-down"))
					return
				}
				time.Sleep(time.Duration(sleepIntervalInt) * time.Second)
			}
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func matchPodName(queryParamPodName string) bool {
	return queryParamPodName == "" || queryParamPodName == getMyPodName()
}
func getMyPodName() string {
	return os.Getenv("MY_POD_NAME")
}
