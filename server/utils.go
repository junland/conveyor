package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Pidfile is a struct that describes a PID file.
type Pidfile struct {
	Name string
}

// CreatePID creates a new PID file.
func CreatePID(name string) *Pidfile {
	log.Debug("Creating and opening PID file...")

	if _, err := os.Stat(name); !os.IsNotExist(err) {
		// file exists
		value, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal("pidfile: failed to read pid ", err)
		}

		pid, err := strconv.Atoi(string(value))
		if err != nil {
			log.Fatal("pidfile: failed to convert string to int ", err)
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			log.Info("Existing PID file does not have a running process, attempting to remove.")
			err := os.Remove(name)
			if err != nil {
				log.Error("pidfile: could not remove existing pidfile ", err)
				os.Exit(1)
			}
			log.Info("Removal complete...")
		} else {
			if err := process.Signal(syscall.Signal(0)); err == nil {
				log.Fatalf("Process %d is already running.", pid)
			}
		}
	}

	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error("pidfile: failed to open pid ", err)
		os.Exit(1)
	}

	defer file.Close()

	log.Debug("Writing PID to PID file...")

	pid := fmt.Sprintf("%d", os.Getpid())
	file.Write([]byte(pid))

	log.Debug("PID creation has been completed...")

	return &Pidfile{name}
}

// RemovePID removes the PID file.
func (pf *Pidfile) RemovePID() {
	log.Debug("Removing PID file...")

	err := os.Remove(pf.Name)
	if err != nil {
		log.Error("pidfile: failed to remove ", err)
	}
	log.Debug("PID file removed...")
}

func WriteScript(f string, i string) error {
	_, err := os.Stat(f)
	if err != nil {
		return err
	}

	script, err := os.OpenFile(f, os.O_WRONLY|os.O_APPEND, 0770)
	if err != nil {
		return err
	}
	defer script.Close()

	_, err = script.WriteString(i)
	if err != nil {
		return err
	}

	log.Debug("Script wrote: " + i)

	return nil
}

func CreateScript(f string) error {
	_, err := os.Stat(f)

	// create file if not exists
	if os.IsNotExist(err) {
		var script, err = os.Create(f)
		if err != nil {
			return err
		}
		defer script.Close()
	}

	script, err := os.OpenFile(f, os.O_RDWR, 0774)
	if err != nil {
		return err
	}
	defer script.Close()

	// Write some text line-by-line to file.
	_, err = script.WriteString("#!/bin/bash \nset -x\n")
	if err != nil {
		return err
	}

	// Save file changes.
	err = script.Sync()
	if err != nil {
		return err
	}

	err = os.Chmod(f, 0744)
	if err != nil {
		return err
	}

	log.Debug("script file created...")

	return nil
}

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
