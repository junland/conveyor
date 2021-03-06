package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

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
		respondError(w, http.StatusBadRequest, "Could not parse json.")
		return
	}

	var newJob JobRequest

	json.Unmarshal(reqBody, &newJob)

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if newJob.Name == "" {
		respondError(w, http.StatusBadRequest, "No job name specified.")
		return
	}

	var exws, expwd, exnqdir, exscript, extime string

	exectime := time.Now().UnixNano() / 1e6

	rand.Seed(exectime)

	extime = strconv.FormatInt(exectime, 10)

	randws := rand.Intn(c.Workers)

	if randws == 0 {
		randws = 1
	}

	exws = strconv.Itoa(randws)

	expwd = fmt.Sprintf("PWD=%s_%s", c.WorkspaceDir, exws)

	exnqdir = fmt.Sprintf("NQDIR=%s_%s", c.WorkersDir, exws)

	exscript = c.WorkersDir + "_" + exws + "/job-scripts.d/" + extime + ".qscript"

	log.Debug("What: " + exscript)

	file, err := os.OpenFile(exscript, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice := []byte("Bytes!\n")
	_, err = file.Write(byteSlice)
	if err != nil {
		log.Fatal(err)
	}

	for _, cmd := range newJob.Commands {
		byteSlice = []byte(cmd)
		_, err = file.Write(byteSlice)
		if err != nil {
			log.Error("Could not write to script file: %s", err)
			respondError(w, http.StatusInternalServerError, "Could not submit job.")
			return
		}
	}

	if err := os.Chmod(exscript, 0777); err != nil {
		log.Fatal(err)
	}

	execq := exec.Command("nqe", "-p", extime, exscript)

	execq.Env = append(os.Environ(), expwd, exnqdir)

	log.Info("Queueing up job for worker " + exws)

	err = execq.Start()
	if err != nil {
		log.Error("Something went wrong with running nq: %s", err)
		respondError(w, http.StatusInternalServerError, "Something went wrong with queue worker.")
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
