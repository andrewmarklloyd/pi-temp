package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Server struct {
		Debug bool `yaml:"debug"`
	} `yaml:"server"`
}

type AppInfo struct {
	TagName string `json:"tag_name"`
}

type HomePageData struct {
	Version       string
	LatestVersion string
	Debug         bool
}

var pin rpio.Pin
var pinNumber int
var debugMode bool
var cfg config
var cronLib *cron.Cron

func main() {
	cfg = readConfig()
	debugMode = cfg.Server.Debug
	version, err := ioutil.ReadFile("static/version")
	if err != nil {
		fmt.Println("unable to open version", err)
		os.Exit(1)
	}

	fmt.Println("Setting up http handlers")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]

		if path == "" {
			latestVersion, err := ioutil.ReadFile("static/latestVersion")
			if err != nil || len(latestVersion) == 0 {
				latestVersion = version
			}
			cfg = readConfig()
			tmpl := template.Must(template.ParseFiles("./static/index.html"))
			data := HomePageData{
				Version:       string(version),
				LatestVersion: string(latestVersion),
				Debug:         cfg.Server.Debug,
			}
			tmpl.Execute(w, data)
		} else {
			if fileExists(path) {
				d, _ := ioutil.ReadFile(string(path))
				w.Write(d)
			} else {
				// fmt.Println(path)
				http.NotFound(w, r)
			}
		}
	})
	http.HandleFunc("/system", systemHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func checkForUpdates() {
	fmt.Println("Checking for updates")
	resp, _ := http.Get("https://api.github.com/repos/andrewmarklloyd/pi-temp/releases/latest")
	var info AppInfo
	err := json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Writing latestVersion to file")
		versionInfo := []byte(info.TagName)
		err = ioutil.WriteFile("./static/latestVersion", versionInfo, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func systemHandler(w http.ResponseWriter, req *http.Request) {
	op := req.FormValue("op")
	var args []string = []string{}
	var command string = ""
	if op == "shutdown" {
		command = "sudo"
		args = []string{"shutdown", "now"}
		fmt.Fprintf(w, "shutting down")
	} else if op == "reboot" {
		command = "sudo"
		args = []string{"reboot", "now"}
		fmt.Fprintf(w, "rebooting")
	} else if op == "update" {
		command = "/home/pi/install/update.sh"
		fmt.Fprintf(w, "updating software")
	} else if op == "check-updates" {
		checkForUpdates()
		fmt.Fprintf(w, "checking for updates")
	} else {
		fmt.Fprintf(w, "command not recognized")
	}
	fmt.Printf("Running command: %s\n", command)
	fmt.Println(debugMode)
	if command != "" && !debugMode {
		cmd := exec.Command(command, args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Start()
		if err != nil {
			fmt.Println("Failed to initiate command:", err)
			os.Exit(1)
		}
		fmt.Printf("Command output: %q\n", out.String())
	}
}

func writeConfig(cfg config) {
	d, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = ioutil.WriteFile("config.yml", d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig() config {
	viper.SetConfigName("config.yml")
	viper.AddConfigPath(currentdir())
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cfg := config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Println(err)
	}
	return cfg
}

func currentdir() (cwd string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return cwd
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
