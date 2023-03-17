package service

import (
	"log"
	"os"
	"path"
	"strings"
)

const ServiceConfigDir = "/etc/services/"

func DiscoverServices() ([]Service, error) {
	files, err := os.ReadDir(ServiceConfigDir)
	if err != nil {
		return nil, err
	}

	services := []Service{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileDescriptor, fileOpenError := os.Open(path.Join(ServiceConfigDir, file.Name()))
		if fileOpenError != nil {
			log.Fatal(fileOpenError)
		}

		sev := Service{
			Identifier: strings.TrimSuffix(file.Name(), path.Ext(file.Name())),
		}

		if err := sev.ReadConfiguration(fileDescriptor); err != nil {
			log.Fatal(err)
		}

		services = append(services, sev)

		fileDescriptor.Close()
	}

	return services, nil
}
