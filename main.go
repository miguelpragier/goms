package main

import (
	_ "database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/miguelpragier/microconfig"
	"github.com/prometheus/common/log"
	"gitlab.com/openvistacode/mysqlkebab"

	"os"
	"path/filepath"
	"runtime"
	"time"
)

const defaultListeningPort = 50000

var (
	webserverListeningPort = defaultListeningPort
	bootTime               = time.Now()
)

var ubi struct {
	rt *mux.Router
	db *mysql.MySQLDriver
}

func uptime() string {
	return time.Now().Sub(bootTime).String()
}

func summary(appName string) {
	fmt.Printf("<‹—-—-—-—-—-—-—-—¤ %s ¤—-—-—-—-—-—-—-—›>\n", appName)

	info := map[string]string{
		"Current executable":        filepath.Base(os.Args[0]),
		"Number of processor cores": fmt.Sprint(runtime.NumCPU()),
		"Operating System:":         runtime.GOOS,
		"Runtime version":           runtime.Version(),
		"Operating system's time":   time.Now().Format(time.RFC3339),
		"gitRevisionHash":           gitRevisionHash,
		"compilationTimestamp":      compilationTimestamp,
	}

	format := "%-30s: %s\n"
	var hist string

	for k, v := range info {
		fmt.Printf(format, k, v)
		hist += fmt.Sprintf("%s - %s\n", k, v)
	}
}

func main() {
	if err := microconfig.Load("conf.json"); err != nil {
		log.Fatal(err)
	}

	summary("THIS SERVICE NAME")

	fmt.Println("Connecting relational database")

	ko := mysqlkebab.KebabOptions{ConnectionAttemptsBeforeError: -1, MaxOpenConnections: 10}

	if d, err := mysqlkebab.NewWithOptions(true, ko); err == nil {
		ubi.db = d
	}

	webserverStart()
}
