package main

import (
	"bufio"
	"docker-builder/models"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Docker struct definition.
type Docker struct {
	config     *models.Config
	dockerExec string
	dockerPath string
}

// NewDocker constructor.
func NewDocker(config *models.Config) *Docker {
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("installing docker is in your future")
	}

	dockerPath, dockerExec := filepath.Split(dockerPath)
	docker := &Docker{dockerPath: dockerPath, dockerExec: dockerExec, config: config}

	return docker
}

func (docker *Docker) runCommand(args []string) (*exec.Cmd, error) {
	cmd := exec.Command(docker.dockerExec, args...)
	cmd.Dir = docker.dockerPath

	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		return cmd, err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		m := scanner.Text()
		log.Printf("\t%s\r\n", m)
	}

	err = cmd.Wait()

	return cmd, err
}

// BuildImage from config.
func (docker *Docker) BuildImage(tag string, image models.Image) {
	args := []string{"build"}

	if len(image.Platforms) > 0 {
		args = append([]string{"buildx"}, args...)
		args = append(args, "--platform", strings.Join(image.Platforms, ","))
	}

	args = append(args, "-t", tag, path.Join(docker.config.Path, image.Path))

	_, err := docker.runCommand(args)
	if err != nil {
		log.Printf("\t%s: %s\r\n", tag, err)
	}
}

// PushImage from config.
func (docker *Docker) PushImage(tag string) {
	args := []string{"push", tag}

	_, err := docker.runCommand(args)
	if err != nil {
		log.Printf("\t%s: %s\r\n", tag, err)
	}
}
