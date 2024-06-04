package main

import (
	"github.com/senyu-up/toolbox/tool/file"
	"github.com/senyu-up/toolbox/tool/geoip"
	"log"
	"os"
	"path"
	"regexp"
)

func main() {
	ExampleCity()

	ExampleIpRegin()
}

func ExampleCity() {
	var re = geoip.CountryConfigsByName("中国")
	log.Printf("get country  by name: %+v", re)

	re = geoip.CountryConfigsByName("USA") // ❌
	log.Printf("get country  by name: %+v", re)

	re = geoip.CountryConfigsByCode("US")
	log.Printf("get country  by name: %+v", re)
}

func getDbPath() string {
	home, err := os.UserHomeDir()
	log.Printf("user home dir is: %s, err %v", home, err)
	var gitPath = "go/pkg/mod/github.com/lionsoul2014"
	pat, err := regexp.Compile(`^ip2region`)
	if err != nil {
		log.Printf("compile regexp err: %v", err)
		return ""
	}
	if files, err := file.ScanDir(path.Join(home, gitPath), 4, &file.Pattern{Pattern: pat}); err != nil {
		log.Printf("scan dir err: %v", err)
		return ""
	} else if 0 < len(files) {
		return path.Join(files[0], "data/ip2region.db")
	}
	return ""
}

func ExampleIpRegin() {
	// init ip2region
	var fullDbPath = getDbPath()
	log.Printf("ip2region db path is: %s", fullDbPath)
	geoip.InitIP2Region(fullDbPath)
	b := geoip.IP2Area("85.185.111.113")
	log.Printf("get ip location: %+v", b)
}
