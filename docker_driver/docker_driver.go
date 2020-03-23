package docker_driver

import (
    "fmt"
    "os/exec"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "golang.org/x/net/context"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/go-connections/nat"
)

type Docker_config struct {
    Image string
    Port [2]string
    Cmd []string
    Memory string       // b k m or g   min is 4M   default is inf
    Cpu string          // between 0.00 to 1.00*cores
    Network string
    Env []string
}

// image should be imagename:version
// hash should be user/image@sha256:digest
func PullImage(image string) (string, error) {
    out, err := exec.Command("docker", "pull", image).Output()
    if err != nil {
        return "", err
    }

    fmt.Println(string(out[:]))

    return "success", nil
}

func ListImages() ([]string, error) {
    //ctx := context.Background()
    cli, err := client.NewEnvClient()
    if err != nil {
        return nil, err
    }

    images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
    if err != nil {
        return nil, err
    }

    var ilist []string
    for _, image := range images {
        ilist = append(ilist, image.ID[7:])     // skips the sha256 tag
    }

    return ilist, nil
}

func ListRunningContainers() ([]string, error) {
    //ctx := context.Background()
    cli, err := client.NewEnvClient()
    if err != nil {
        return nil, err
    }

    containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
    if err != nil {
        return nil, err
    }

    var clist []string
    for _, container := range containers {
        clist = append(clist, container.ID)
    }

    return clist, nil
}

// stopping container
func StopContainer(cont string) (string, error) {
    _, err := exec.Command("docker", "stop", cont).Output()
    if err != nil {
        return "", err
    }

    return "success", nil
}

// deleting container
func DeleteContainer(cont string) (string, error) {
    _, err := exec.Command("docker", "rm", cont).Output()
    if err != nil {
        return "", err
    }

    return "success", nil
}

// restarting container
func RestartContainer(cont string) (string, error) {
    _, err := exec.Command("docker", "restart", cont).Output()
    if err != nil {
        return "", err
    }

    return "success", nil
}

// resizing a container instance on the fly
func ResizeContainer(cont string, size string) (string, error) {
    _, err := exec.Command("docker", "container", "update", "-m", size, cont).Output()
    if err != nil {
        return "", err
    }

    return "success", nil
}

// create and run container - interactive and detached set
// image (already pulled) should be imagename:version
// default/empty cmd is /bin/bash
func RunContainer(opt Docker_config) (string, error) {
    ctx := context.Background()
    cli, err := client.NewEnvClient()
    if err != nil {
        return "", err
    }

    resp, err := cli.ContainerCreate(ctx, &container.Config{
        Image: opt.Image,
        Cmd: opt.Cmd,
        ExposedPorts: nat.PortSet{ nat.Port(opt.Port[0]) : struct{}{} },
        Tty: true,
        Env: opt.Env,
    },
    &container.HostConfig{
        NetworkMode: container.NetworkMode(opt.Network),
        PortBindings: nat.PortMap{ nat.Port(opt.Port[0]) :
            []nat.PortBinding{ nat.PortBinding{ HostPort: opt.Port[1] } }, },
    },
    nil, "")
    if err != nil {
        return "", err
    }

    // updating memory and cpu
    if (opt.Memory != "") {
        ResizeContainer(resp.ID, opt.Memory)
    }

    if (opt.Cpu != "") {
        _, err := exec.Command("docker", "container", "update", "--cpus", opt.Cpu, resp.ID).Output()
        if err != nil {
            return "", err
        }
    }

    err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
    if err != nil {
        return "", err
    }

    //fmt.Println(resp.ID)

    return resp.ID, nil
}
