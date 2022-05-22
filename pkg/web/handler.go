package web

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
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
