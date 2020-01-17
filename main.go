package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func main() {

	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("installing docker is in your future")
	}

	fmt.Printf("docker installed in: %s\n", dockerPath)
	dockerPath, dockerExec := filepath.Split(dockerPath)

	configPath := *flag.String("config", "", "Path del file di configurazione se diverso dalla cartella di esecuzione")

	flag.Parse()

	if configPath == "" {
		configPath = "./"
	}

	configPath = path.Join(configPath, "config.yaml")

	reader, err := os.Open(configPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	config := &Config{}
	decoder := yaml.NewDecoder(reader)

	err = decoder.Decode(config)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Folders in: %s\n", config.Path)
	for k, v := range config.Images {
		fmt.Printf("Start build: %s\n", k)
		cmd := exec.Command(dockerExec, "buildx", "build", "--platform", strings.Join(v.Platforms, ","), "-t", k, path.Join(config.Path, v.Path))
		cmd.Dir = dockerPath

		stderr, _ := cmd.StderrPipe()

		err = cmd.Start()
		if err != nil {
			log.Fatalf("\t%s: %s\n", k, err)
		}

		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Printf("\t%s\n", m)
		}

		err = cmd.Wait()
		if err != nil {
			log.Fatalf("\t%s: %s\n", k, err)
		}
	}

}

type Config struct {
	Path   string           `yaml:"path"`
	Images map[string]Image `yaml:"images,omitempty"`
}

type Image struct {
	Path      string   `yaml:"path"`
	Platforms []string `yaml:"platforms,flow"`
}
