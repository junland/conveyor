package cmd

import (
	"fmt"

	"github.com/junland/conveyor/server"
	flag "github.com/spf13/pflag"
)

// BinVersion describes built binary version.
var BinVersion string

// GoVersion Describes Go version that was used to build the binary.
var GoVersion string

// Default parameters when program starts without flags or environment variables.
const (
	defLvl          = "info"
	defAccess       = true
	defPort         = "8080"
	defPID          = "/var/run/conveyor.pid"
	defTLS          = false
	defCert         = ""
	defKey          = ""
	defWorkers      = 2
	defWorkersDir   = "/tmp/worker"
	defWorkspaceDir = "/tmp/workspace"
)

var (
	confLogLvl, confPort, confPID, confCert, confKey, confWorkersDir, confWorkspaceDir string
	enableTLS, enableAccess, version, help                                             bool
	confWorkers                                                                        int
)

// init defines configuration flags and environment variables.
func init() {
	flags := flag.CommandLine
	flags.StringVar(&confLogLvl, "log-level", GetEnvString("CONVEYOR_LOG_LEVEL", defLvl), "Specify log level for output.")
	flags.BoolVar(&enableAccess, "access-log", GetEnvBool("CONVEYOR_ACCESS_LOG", defAccess), "Specify weather to run with or without HTTP access logs.")
	flags.StringVar(&confPort, "port", GetEnvString("CONVEYOR_SERVER_PORT", defPort), "Starting server port.")
	flags.StringVar(&confPID, "pid-file", GetEnvString("CONVEYOR_SERVER_PID", defPID), "Specify server PID file path.")
	flags.BoolVar(&enableTLS, "tls", GetEnvBool("CONVEYOR_TLS", defTLS), "Specify weather to run server in secure mode.")
	flags.StringVar(&confCert, "tls-cert", GetEnvString("CONVEYOR_TLS_CERT", defCert), "Specify TLS certificate file path.")
	flags.StringVar(&confKey, "tls-key", GetEnvString("CONVEYOR_TLS_KEY", defKey), "Specify TLS key file path.")
	flags.StringVar(&confWorkspaceDir, "workspace-dir", GetEnvString("CONVEYOR_WORKSPACE_DIR", defWorkspaceDir), "Specify the working directory for builds.")
	flags.IntVar(&confWorkers, "workers", GetEnvInt("CONVEYOR_WORKERS", defWorkers), "Specify amount of executors to process requests.")
	flags.StringVar(&confWorkersDir, "workers-dir", GetEnvString("CONVEYOR_WORKERS_DIR", defWorkersDir), "Specify the working directory for builds.")
	flags.BoolVarP(&help, "help", "h", false, "Show this help")
	flags.BoolVar(&version, "version", false, "Display version information")
	flags.SortFlags = false
	flag.Parse()
}

// PrintHelp prints help text.
func PrintHelp() {
	fmt.Printf("Usage: conveyor [options] <command> [<args>]\n")
	fmt.Printf("\n")
	fmt.Printf("A simple web app template.\n")
	fmt.Printf("\n")
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

// PrintVersion prints version information about the binary.
func PrintVersion() {
	fmt.Printf("Made with love.\n")
	fmt.Printf("Version: %s\n", BinVersion)
	fmt.Printf("Go Version %s\n", GoVersion)
	fmt.Printf("License: MIT\n")
}

// Run is the entry point for starting the command line interface.
func Run() {
	config := server.Config{
		LogLvl:       confLogLvl,
		Access:       enableAccess,
		Port:         confPort,
		PID:          confPID,
		TLS:          enableTLS,
		Cert:         confCert,
		Key:          confKey,
		WorkspaceDir: confWorkspaceDir,
		Workers:      confWorkers,
		WorkersDir:   confWorkersDir,
	}

	if version {
		PrintVersion()
		return
	}

	if help {
		PrintHelp()
		return
	}

	server.Start(config)
}
