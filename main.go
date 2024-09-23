package main

import (
	"encoding/json"
	"fmt"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection"
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
		Name:  "Engineering Kiosk",
		Email: "Kiosk@Engineering.com",
	},
}

var currentLeader string

type CurrentLeader struct {
	CurrentLeader string `json:"currentLeader"`
  	IAmTheLeader string `json:"iAmTheLeader"`
}

func getCurrentLeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	
	iAmTheLeader:="currently I am not the leader :-("
	if getMyPodName() == currentLeader {
		iAmTheLeader="I am the leader :-)"
	}
	json.NewEncoder(w).Encode(CurrentLeader{
		CurrentLeader: currentLeader,
		IAmTheLeader: iAmTheLeader,
	})	
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
	r.HandleFunc("/leader", getCurrentLeader).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// request handlers
	r.Use(loggingHandler)
	r.Use(slowDownRequestHandler)
	r.Use(fakingHttpStatusCodeHandler)

	// start leader election
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal("failed to get in cluster config", err)
	}
	client := clientset.NewForConfigOrDie(kubeConfig)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		klog.Info("Received termination, signaling shutdown")
		cancel()
	}()
	leaseId := getMyPodName()
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "test-application",
			Namespace: "test",
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: leaseId,
		},
	}

	// start the leader election code loop
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   16 * time.Second,
		RenewDeadline:   8 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
				klog.Infof("started leading")
			},
			OnStoppedLeading: func() {
				// we can do cleanup here
				klog.Infof("leader lost: %s", leaseId)
				os.Exit(0)
			},
			OnNewLeader: func(identity string) {
				currentLeader = identity
				// we're notified when new leader elected
				if identity == leaseId {
					// I just got the lock
					return
				}
				klog.Infof("new leader elected: %s", identity)
			},
		},
	})

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
