package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mgutz/logxi/v1"
)

var (
	logLevel = flag.String("loglevel", "debug", "Set the desired log level")

	listen = flag.String("listen", ":8080", "Address to bind to")
	path   = flag.String("path", "./", "Path served as document root.")
)

type testSlot struct {
	secondSlot int    // The second at which the served directory activates
	dir        string // The directory that activates
}

type testWindow struct {
	startTime time.Time
	slots     []*testSlot
	sync.Mutex
}

var (
	// create Logger interface
	logW = log.NewLogger(log.NewConcurrentWriter(os.Stdout), "HttpRoller")

	testSchedule = testWindow{
		startTime: time.Now().Round(time.Second),
		slots:     []*testSlot{},
	}
)

func main() {

	flag.Parse()

	switch strings.ToLower(*logLevel) {
	case "debug":
		logW.SetLevel(log.LevelDebug)
	case "info":
		logW.SetLevel(log.LevelInfo)
	}

	_, err := filepath.Abs(*path)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}

	// Load the inital per second slots from the
	// scenario directory
	loadTest(*path)

	// Start a service function that tracks over time the slots
	// and scenarios being used
	//
	go auditWindow()

	http.HandleFunc("/", serveHandler)

	if err = http.ListenAndServe(*listen, nil); err != nil {
		logW.Warn(err.Error())
	}
}

// loadTest examines the scenario directory for the serve directories
// that will be used and loads them into the testSchedule
//
func loadTest(scenario string) (err error) {
	testSchedule.Lock()
	defer testSchedule.Unlock()

	testSchedule.startTime = time.Now().Round(time.Second)
	testSchedule.slots = []*testSlot{}

	err = filepath.Walk(scenario,
		func(path string, f os.FileInfo, err error) error {
			if path == scenario {
				return nil
			}
			slot, err := strconv.Atoi(f.Name())
			if err == nil {
				testSchedule.slots = append(testSchedule.slots, &testSlot{
					dir:        path,
					secondSlot: slot})
			}
			if f.IsDir() {
				return filepath.SkipDir
			}
			return nil
		})
	if err != nil {
		logW.Warn(fmt.Sprintf("could not load test scenario from %s due to %s", scenario, err.Error()), "error", err)
	}

	// Sort our slots  ascending order and we are done
	sort.Slice(testSchedule.slots, func(i, j int) bool {
		return testSchedule.slots[i].secondSlot < testSchedule.slots[j].secondSlot
	})

	logW.Debug(fmt.Sprintf("loaded scenario %s", scenario))

	return err
}

func getSlotDir() (dir string) {
	testSchedule.Lock()
	defer testSchedule.Unlock()

	second := int(time.Since(testSchedule.startTime).Seconds()) - 1
	slot := sort.Search(len(testSchedule.slots), func(i int) bool { return testSchedule.slots[i].secondSlot >= second })

	if slot < len(testSchedule.slots) && testSchedule.slots[slot].secondSlot == second {
		return testSchedule.slots[slot].dir
	}

	if slot >= len(testSchedule.slots) {
		slot = len(testSchedule.slots) - 1
	} else {
		if slot < 0 {
			slot = 0
		}
	}

	return testSchedule.slots[slot].dir
}

func auditWindow() {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			logW.Debug(fmt.Sprintf("using %s", getSlotDir()))

			files, _ := ioutil.ReadDir(getSlotDir())
			for _, aFile := range files {
				if aFile.Name() == "finish" {
					loadTest(*path)
					break
				}
			}
		}
	}
}

func serveHandler(w http.ResponseWriter, r *http.Request) {

	// Locate from the current test scenario which
	// directory is the appropriate one to serve up
	//
	file := getSlotDir() + r.URL.Path
	logW.Debug(fmt.Sprintf("serving %s", file))

	http.ServeFile(w, r, file)
}
