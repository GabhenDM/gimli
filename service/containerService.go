package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ahmetalpbalkan/dlog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/gabhendm/gimli/utils"
)

func StartContainer(imageURL string, cmd []string, debugFlag bool) (io.Reader, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	reader, err := cli.ImagePull(ctx, imageURL, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	if debugFlag {
		io.Copy(os.Stdout, reader)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageURL,
		Cmd:   cmd,
		Tty:   false,
	}, nil, nil, nil, utils.CreateContainerName(imageURL))
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	if debugFlag {
		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			return nil, err
		}
		stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	r := dlog.NewReader(out)

	if err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
		return nil, err
	}
	fmt.Println(fmt.Sprintf("Container: %s - Finished Executing...", utils.CreateContainerName(imageURL)))

	return r, nil
}

func StartContainerDetached(imageURL string, cmd []string, portBindingHost string, portBindingContainer string, debugFlag bool) (container.ContainerCreateCreatedBody, error) {

	containerHostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port(portBindingContainer): {{HostIP: "127.0.0.1", HostPort: portBindingHost}}},
	}

	containerConfig := &container.Config{
		Image:        imageURL,
		ExposedPorts: nat.PortSet{nat.Port(portBindingContainer): struct{}{}},
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, imageURL, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	if debugFlag {
		io.Copy(os.Stdout, out)
	}
	resp, err := cli.ContainerCreate(ctx, containerConfig, containerHostConfig, nil, nil, utils.CreateContainerName(imageURL))
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	if debugFlag {
		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}
		stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	}

	return resp, err
}

func ListRunningContainers() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID, " --- ", container.Image)
	}
}

func RemoveRunningContainers() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if strings.Contains(container.Names[0], "gimli") {
			fmt.Print("Stopping container ", container.ID[:10], "... ")
			if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
				return err
			}
			if err = cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
				return err
			}
			fmt.Println("Success")
		}
	}
	return nil
}
