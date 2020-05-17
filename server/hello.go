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
	log "github.com/sirupsen/logrus"
)

// JobRequest describes the statement of work.
type JobRequest struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`
}

// CreateJob is a function that collects and parses incoming jobs.
func (c *Config) CreateJob(w http.ResponseWriter, r *http.Request) {

	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
    
    reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Something went wrong with parsing the json request: %s", err)
		respondJSON(w, http.StatusBadRequest, "Could not parse json.")
		return
	}

	var newJob JobRequest

	json.Unmarshal(reqBody, &newJob)

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if newJob.Name == "" {
		respondJSON(w, http.StatusBadRequest, "No job name specified.")
		return
	}

    var exws, expwd, exnqdir, excmd string

	for _, cmd := range newJob.Commands {
        excmd += cmd + ";"
	}
	
	log.Debug("Created command structure...")

	exectime := time.Now().UnixNano() / 1e6

	rand.Seed(exectime)

	exws = strconv.Itoa(rand.Intn(c.Workers))

	expwd = fmt.Sprintf("PWD=%s_%s", c.WorkspaceDir, exws)

	exnqdir = fmt.Sprintf("NQDIR=%s_%s", c.WorkersDir, exws)

	execq := exec.Command("nqe", "-p " + strconv.Itoa(exectime), "-- ", excmd)

	execq.Env = append(os.Environ(),expwd,exnqdir,)

	log.Info("Queueing up job for worker " + exws)

	err = execq.Start()
	if err != nil {
		log.Error("Something went wrong with running nq: %s", err)
		respondJSON(w, http.StatusInternalServerError, "Could not parse json.")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Job Submitted"})

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
