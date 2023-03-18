package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func apiAddressesHandler(response http.ResponseWriter, request *http.Request) {
	type smog struct {
		Interface string
		Addresses []string
	}

	var sm []smog

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {

		in := smog{
			Interface: iface.Name,
		}
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			in.Addresses = append(in.Addresses, addr.String())
		}
		sm = append(sm, in)
	}

	encoder := json.NewEncoder(response)
	encoder.SetEscapeHTML(true)
	encoder.SetIndent("", "  ")
	encoder.Encode(sm)
}

func apiFilesHandler(response http.ResponseWriter, request *http.Request) {
	var files []string
	filepath.Walk("/", func(name string, info os.FileInfo, err error) error {
		files = append(files, name)
		return nil
	})
	encoder := json.NewEncoder(response)
	encoder.SetEscapeHTML(true)
	encoder.SetIndent("", "  ")
	encoder.Encode(files)
}

func apiProcessesHandler(response http.ResponseWriter, request *http.Request) {
	plist, _ := processes("/proc")
	encoder := json.NewEncoder(response)
	encoder.SetEscapeHTML(true)
	encoder.SetIndent("", "  ")
	encoder.Encode(plist)
}

// Process - a simple process
type Process struct {
	Pid     int
	PPid    int
	Binary  string
	State   string
	Cmdline string
}

// Processes - List all running processes
func processes(procdir string) ([]Process, error) {
	var parr []Process

	files, err := os.ReadDir(procdir)

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			var p Process
			matched, _ := regexp.MatchString(`^\d+$`, file.Name())
			if matched {
				p.Pid, _ = strconv.Atoi(file.Name())
				p.refresh()
				parr = append(parr, p)
			}
		}
	}

	return parr, nil
}

// Refresh - readout data about the process
func (p *Process) refresh() error {
	statPath := fmt.Sprintf("/proc/%d/stat", p.Pid)
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", p.Pid)
	exelinkPath := fmt.Sprintf("/proc/%d/exe", p.Pid)
	dataBytes, err := os.ReadFile(statPath)
	if err != nil {
		return err
	}
	cmdlineBytes, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return err
	}

	data := string(dataBytes)
	p.Cmdline = string(cmdlineBytes)

	binStart := strings.IndexRune(data, '(') + 1
	binEnd := strings.IndexRune(data[binStart:], ')')

	p.Binary, _ = os.Readlink(exelinkPath)
	data = data[binStart+binEnd+2:]
	fmt.Sscanf(data,
		"%s %d",
		&p.State,
		&p.PPid)

	return err
}

func UsageSSEHandlerBuilder() http.HandlerFunc {
	events := make(chan Event)

	bToMb := func(b uint64) uint64 {
		return b / 1024 / 1024
	}

	go func() {
		for {
			// get load average
			loadavg, _ := os.ReadFile("/proc/loadavg")

			// get memory usage
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			// get free disk space
			var stat unix.Statfs_t
			unix.Statfs("/tmp/", &stat)

			// send load average
			events <- Event{Message: struct {
				LoadAvg   string `json:"loadavg"`
				MemUsed   uint64 `json:"memused"`
				DiskAvail uint64 `json:"diskavail"`
			}{
				LoadAvg:   string(loadavg),
				MemUsed:   bToMb(m.Sys),
				DiskAvail: bToMb(stat.Bavail * uint64(stat.Bsize)),
			}}
		}
	}()

	return ServerSendEventsHandlerBuilder(events)
}
