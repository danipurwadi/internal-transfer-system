// Package docker provides utilities for working with Docker.
package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type Container struct {
	ID       string
	HostPort string
	Host     string
	Port     int
}

// StartContainer starts the specified container for running tests.
func StartContainer(image string, port string, dockerArgs []string, appArgs []string) (*Container, error) {
	arg := []string{"run", "-P", "-d"}
	arg = append(arg, dockerArgs...)
	arg = append(arg, image)
	arg = append(arg, appArgs...)
	var out bytes.Buffer
	cmd := exec.Command("docker", arg...)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("could not start container %s: %w", image, err)
	}

	id := out.String()[:12]
	_, hostPort, err := extractIPPort(id, port)
	if err != nil {
		StopContainer(id)
		return nil, fmt.Errorf("could not extract ip/port: %w", err)
	}

	c := Container{
		ID:       id,
		HostPort: net.JoinHostPort("0.0.0.0", hostPort),
	}


	return &c, nil
}

// StopContainer stops and removes the specified container.
func StopContainer(id string) error {
	if err := exec.Command("docker", "stop", id).Run(); err != nil {
		return fmt.Errorf("could not stop container: %w", err)
	}

	if err := exec.Command("docker", "rm", id, "-v").Run(); err != nil {
		return fmt.Errorf("could not remove container: %w", err)
	}

	return nil
}

// DumpContainerLogs outputs logs from the running docker container.
func DumpContainerLogs(id string) []byte {
	out, err := exec.Command("docker", "logs", id).CombinedOutput()
	if err != nil {
		return nil
	}

	return out
}

func extractIPPort(id string, port string) (hostIP string, hostPort string, err error) {
	tmpl := fmt.Sprintf("[{{range $k,$v := (index .NetworkSettings.Ports \"%s/tcp\")}}{{json $v}}{{end}}]", port)

	var out bytes.Buffer
	cmd := exec.Command("docker", "inspect", "-f", tmpl, id)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("could not inspect container %s: %w", id, err)
	}

	// When IPv6 is turned on with Docker.
	// Got  [{"HostIp":"0.0.0.0","HostPort":"49190"}{"HostIp":"::","HostPort":"49190"}]
	// Need [{"HostIp":"0.0.0.0","HostPort":"49190"},{"HostIp":"::","HostPort":"49190"}]
	data := strings.ReplaceAll(out.String(), "}{", "},{")

	var docs []struct {
		HostIP   string `json:"HostIp"`
		HostPort string `json:"HostPort"`
	}
	if err := json.Unmarshal([]byte(data), &docs); err != nil {
		return "", "", fmt.Errorf("could not decode json: %w", err)
	}

	for _, doc := range docs {
		if doc.HostIP != "::" {
			// Podman keeps HostIP empty instead of using 0.0.0.0.
			// - https://github.com/containers/podman/issues/17780
			if doc.HostIP == "" {
				return "localhost", doc.HostPort, nil
			}

			return doc.HostIP, doc.HostPort, nil
		}
	}

	return "", "", fmt.Errorf("could not locate ip/port")
}
