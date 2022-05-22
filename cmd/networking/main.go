package main

import (
	"encoding/json"
	"flag"
	"github.com/vishvananda/netlink"
	"log"
	"net"
	"os"
)

type Config struct {
	Routes     []Route     `json:"routes"`
	Interfaces []Interface `json:"interfaces"`
}

type Route struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
}

type Interface struct {
	Name        string `json:"name"`
	AddressIPv4 string `json:"address4"`
}

var config Config

func init() {
	flag.Parse()
	configFilePath := flag.Arg(0)

	configFile, fileOpenError := os.Open(configFilePath)
	if fileOpenError != nil {
		log.Println(fileOpenError)
		os.Exit(1)
	}

	jsonDecodeError := json.NewDecoder(configFile).Decode(&config)
	if jsonDecodeError != nil {
		log.Println(jsonDecodeError)
		os.Exit(1)
	}
}

func main() {
	for _, iface := range config.Interfaces {
		link, linkFoundError := netlink.LinkByName(iface.Name)
		if linkFoundError != nil {
			log.Println(linkFoundError)
			continue
		}
		if linkSetUpError := netlink.LinkSetUp(link); linkSetUpError != nil {
			log.Println(linkSetUpError)
			continue
		}

		address, addressParseError := netlink.ParseAddr(iface.AddressIPv4)
		if addressParseError != nil {
			log.Println(addressParseError)
			continue
		}

		if addrAddError := netlink.AddrAdd(link, address); addrAddError != nil {
			log.Println(addrAddError)
			continue
		}
	}

	for _, route := range config.Routes {
		_, destination, destinationParseError := net.ParseCIDR(route.Destination)
		if destinationParseError != nil {
			log.Println(destinationParseError)
			continue
		}
		gateway := net.ParseIP(route.Gateway)
		if gateway.String() == "" {
			log.Printf("invalid gateway address %+q", route.Gateway)
			continue
		}
		r := netlink.Route{
			Dst:   destination,
			Gw:    gateway,
			Table: 0,
		}

		routeAddError := netlink.RouteAdd(&r)
		if routeAddError != nil {
			log.Println(routeAddError)
		}
	}
}
