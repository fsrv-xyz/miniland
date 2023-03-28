package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"

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

	filterType := request.URL.Query().Get("type")
	typeExist, _ := regexp.MatchString(`^(user|system)$`, filterType)
	if !typeExist && filterType != "" {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Invalid type"))
		return
	}

	// filter by type
	if filterType != "" {
		var filtered []Process
		for _, p := range plist {
			if p.Type == ProcessType(filterType) {
				filtered = append(filtered, p)
			}
		}
		plist = filtered
	}

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
	Uid     int
	Gid     int
	Type    ProcessType
	Comm    string
}

type ProcessType string

const (
	ProcessTypeUser   ProcessType = "user"
	ProcessTypeSystem ProcessType = "system"
)

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
	commandPath := fmt.Sprintf("/proc/%d/comm", p.Pid)
	exelinkPath := fmt.Sprintf("/proc/%d/exe", p.Pid)
	dataBytes, err := os.ReadFile(statPath)
	if err != nil {
		return err
	}
	cmdlineBytes, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return err
	}
	p.Comm = func(n []byte, _ error) string { return strings.TrimSuffix(string(n), "\n") }(os.ReadFile(commandPath))

	fi, _ := os.Stat(path.Join("/proc", strconv.Itoa(p.Pid)))
	p.Uid = int(fi.Sys().(*syscall.Stat_t).Uid)
	p.Gid = int(fi.Sys().(*syscall.Stat_t).Gid)

	data := string(dataBytes)
	p.Cmdline = string(cmdlineBytes)
	p.Cmdline = strings.ReplaceAll(p.Cmdline, "\x00", " ")
	p.Cmdline = strings.TrimSuffix(p.Cmdline, " ")

	binStart := strings.IndexRune(data, '(') + 1
	binEnd := strings.IndexRune(data[binStart:], ')')

	p.Binary, _ = os.Readlink(exelinkPath)
	data = data[binStart+binEnd+2:]
	fmt.Sscanf(data,
		"%s %d",
		&p.State,
		&p.PPid)

	if p.Cmdline == "" {
		p.Type = ProcessTypeSystem
	} else {
		p.Type = ProcessTypeUser
	}

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

			type Filesystem struct {
				Path  string `json:"path"`
				Total uint64 `json:"total"`
				Used  uint64 `json:"used"`
			}
			var filesystems []Filesystem

			// get a list of all mounted filesystems from /proc/mounts
			mounts, _ := os.ReadFile("/proc/mounts")
			for _, mount := range strings.Split(string(mounts), "\n") {
				if !strings.Contains(mount, " /") {
					continue
				}
				parts := strings.Split(mount, " ")
				if len(parts) < 2 {
					continue
				}
				var stat unix.Statfs_t
				unix.Statfs(parts[1], &stat)

				filesystems = append(filesystems, Filesystem{
					Path:  parts[1],
					Total: bToMb(stat.Blocks * uint64(stat.Bsize)),
					Used:  bToMb((stat.Blocks - stat.Bfree) * uint64(stat.Bsize)),
				})
			}

			// send data
			events <- Event{Message: struct {
				LoadAvg     string       `json:"loadavg"`
				MemUsed     uint64       `json:"memused"`
				Filesystems []Filesystem `json:"filesystems"`
			}{
				LoadAvg:     string(loadavg),
				MemUsed:     bToMb(m.Sys),
				Filesystems: filesystems,
			}}
		}
	}()

	return ServerSendEventsHandlerBuilder(events)
}
