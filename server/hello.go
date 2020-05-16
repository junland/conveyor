package server

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"os/exec"
	"math/rand"
	"time"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// JobRequest describes the statement of work.
type JobRequest struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`
}

// Collector is a fucntion that collects and parses incoming jobs.
func (c *Config) Collector(w http.ResponseWriter, r *http.Request) {

	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
    
    reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	var newJob JobRequest

	json.Unmarshal(reqBody, &newJob)

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if newJob.name == "" {
		http.Error(w, "You must specify a name for this job.", http.StatusBadRequest)
		return
	}

	var excmd string
	var expwd string
	var exnqdir string

	for _, cmd := range newJob.commands {
        excmd += str + ";"
	}
	
	log.Debug("Created command structure...")

	rand.Seed(time.Now().UnixNano())

	exws := strconv.Itoa(rand.Intn(c.Workers))

	expwd := "PWD=" + c.WorkspaceDir + "_" + esws

	exnqdir := "NQDIR=" + c.WorkersDir + "_" + esws

	execq := exec.Command("nq", excmd)

	execq.Env = append(os.Environ(),expwd,exnqdir,)

	log.Info("Queueing up job for worker " + esws)

	err := execq.Start()

	if err != nil {
		log.Error("Something went wrong with running nq: " + err)
	}

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)

	return
}

// helloRootHandle is a handle.
func helloRootHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(202)
	w.Write([]byte("I am root."))
}

// helloGlobalHandle is a example handler.
func helloGlobalHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello everyone!"))
}

// helloNameHandle is a example parameter handler.
func helloNameHandle(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())
	name := ps.ByName("name")
	w.WriteHeader(201)
	fmt.Fprintf(w, "Hello, %s\n", name)
}
