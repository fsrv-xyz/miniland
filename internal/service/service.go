package service

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	// configFilePath - Path for all service configuration files
	configFilePath string = "/etc/services/"
	LogDirectory   string = "/log"
)

type serviceConfiguration struct {
	Name      string   `json:"name"`
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}

func (config *serviceConfiguration) readConfigurationFile(filename string) *serviceConfiguration {
	serviceConfigurationFile, serviceConfigurationFileOpenError := os.Open(filename)
	if serviceConfigurationFileOpenError != nil {
		log.Printf("serviceConfigurationFile.Get #%v ", serviceConfigurationFileOpenError)
	}
	serviceConfigurationFileDecodeError := json.NewDecoder(serviceConfigurationFile).Decode(config)
	if serviceConfigurationFileDecodeError != nil {
		log.Fatalf("Unmarshal: %v", serviceConfigurationFileDecodeError)
	}
	return config
}

type service struct {
	Name          string
	Configuration serviceConfiguration
}

func (s *service) Run(wg *sync.WaitGroup) *service {
	defer wg.Done()
	log.Printf("running service: %v", s.Name)
	cmd := exec.Command(s.Configuration.Command, s.Configuration.Arguments...)
	if err := s.execute(cmd); err != nil {
		fmt.Println(err)
	}
	return s
}

func (s *service) logger(kind string, v []byte) {
	fileSingle, err := os.OpenFile(fmt.Sprintf("%s/%s.log", LogDirectory, s.Name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer fileSingle.Close()
	loggerSingle := log.New(fileSingle, "", log.LstdFlags)

	fileService, err := os.OpenFile(fmt.Sprintf("%s/services.log", LogDirectory), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer fileService.Close()
	loggerServices := log.New(fileService, "", log.LstdFlags)

	for _, msg := range strings.Split(string(v), "\n") {
		if msg != "" {
			loggerSingle.Printf("%s %s - %s", s.Name, kind, msg)
			loggerServices.Printf("%s %s - %s", s.Name, kind, msg)
		}
	}
}

func copyLogs(r io.Reader, kind string, logfn func(kind string, args []byte)) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			logfn(kind, buf[0:n])
		}
		if err != nil {
			break
		}
	}
}

func (s *service) execute(cmd *exec.Cmd) error {
	var wg sync.WaitGroup

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		copyLogs(stdout, "stdout", s.logger)
	}()

	go func() {
		defer wg.Done()
		copyLogs(stderr, "stderr", s.logger)
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func RunServices(c chan bool) {
	var wg sync.WaitGroup
	files, err := ioutil.ReadDir(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		var c serviceConfiguration
		c.readConfigurationFile(fmt.Sprintf("%s/%s", configFilePath, file.Name()))
		var s service
		s.Name = c.Name
		s.Configuration = c
		wg.Add(1)
		go s.Run(&wg)
	}
	wg.Wait()
	c <- true
}
