package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/docker/docker/api/types/filters"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

type manager struct {
	cli *client.Client
}

func newManager() *manager {
	os.Setenv("DOCKER_API_VERSION", "1.35")

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return &manager{cli}
}

func (m *manager) updateServices(config config) {
	services, err := m.cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		log.Println("Could not list services")
		return
	}

	_, err = m.cli.RegistryLogin(context.Background(), types.AuthConfig{
		Username:      config.Registry.User,
		Password:      config.Registry.Password,
		ServerAddress: config.Registry.URL,
	})

	if err != nil {
		log.Printf("Could not login to the registry: %v", err)
		return
	}

	for _, service := range services {
		// update if not blacklisted
		if !contains(config.Blacklist, service.Spec.Name) {

			imageNameHash := service.Spec.TaskTemplate.ContainerSpec.Image
			matches := regexp.MustCompile("(.*)@+").FindStringSubmatch(imageNameHash)

			log.Printf("Updating service %s with image %s\n",
				service.Spec.Name, matches[1])
			out, err := exec.Command(
				"docker",
				"service",
				"update",
				service.Spec.Name,
				"-d",
				"--image",
				matches[1],
				"--with-registry-auth",
			).CombinedOutput()
			if err != nil {
				log.Printf("Error executing update command: %v\n", err)
			}

			log.Println(string(out))
		}
	}
}

func (m *manager) getTask(serviceID string) ([]swarm.Task, error) {
	state := filters.KeyValuePair{Key: "desired-state", Value: "running"}
	service := filters.KeyValuePair{Key: "service", Value: serviceID}
	return m.cli.TaskList(context.Background(), types.TaskListOptions{
		Filters: filters.NewArgs(state, service),
	})
}

func (m *manager) getServices() ([]swarm.Service, error) {
	return m.cli.ServiceList(context.Background(), types.ServiceListOptions{})
}

func (m *manager) getLogs(serviceID string) (io.ReadCloser, error) {
	opt := types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
		Tail:       "30",
		Details:    false,
	}
	return m.cli.ServiceLogs(context.Background(), serviceID, opt)
}
