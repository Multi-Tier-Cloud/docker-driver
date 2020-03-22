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
}

// image should be imagename:version
// TODO: find out how to grab image by hash
func PullImage(image string) {
    out, err := exec.Command("docker", "pull", image).Output()
    if err != nil {
        panic(err)
    }

    fmt.Println(string(out[:]))
}

func ListImages() []string {
    //ctx := context.Background()
    cli, err := client.NewEnvClient()
    if err != nil {
        panic(err)
    }

    images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
    if err != nil {
        panic(err)
    }

    var ilist []string
    for _, image := range images {
        ilist = append(ilist, image.ID[7:])     // skips the sha256 tag
    }

    return ilist
}

func ListRunningContainers() []string {
    //ctx := context.Background()
    cli, err := client.NewEnvClient()
    if err != nil {
        panic(err)
    }

    containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
    if err != nil {
        panic(err)
    }

    var clist []string
    for _, container := range containers {
        clist = append(clist, container.ID)
    }

    return clist
}

// stopping container
func StopContainer(cont string) {
    _, err := exec.Command("docker", "stop", cont).Output()
    if err != nil {
        panic(err)
    }
}

// deleting container
func DeleteContainer(cont string) {
    _, err := exec.Command("docker", "rm", cont).Output()
    if err != nil {
        panic(err)
    }
}

// deleting container
func RestartContainer(cont string) {
    _, err := exec.Command("docker", "restart", cont).Output()
    if err != nil {
        panic(err)
    }
}

// resizing a container instance on the fly
func ResizeContainer(cont string, size string) {
    _, err := exec.Command("docker", "container", "update", "-m", size, cont).Output()
    if err != nil {
        panic(err)
    }
}

// create and run container - interactive and detached set
// image (already pulled) should be imagename:version
// default/empty cmd is /bin/bash
func RunContainer(opt Docker_config) string {
    ctx := context.Background()
    cli, err := client.NewEnvClient()
    if err != nil {
        panic(err)
    }

    resp, err := cli.ContainerCreate(ctx, &container.Config{
        Image: opt.Image,
        Cmd: opt.Cmd,
        ExposedPorts: nat.PortSet{ nat.Port(opt.Port[0]) : struct{}{} },
        Tty: true,
    },
    &container.HostConfig{
        PortBindings: nat.PortMap{ nat.Port(opt.Port[0]) :
            []nat.PortBinding{ nat.PortBinding{ HostPort: opt.Port[1] } }, },
    },
    nil, "")
    if err != nil {
        panic(err)
    }

    // updating memory and cpu
    if (opt.Memory != "") {
        ResizeContainer(resp.ID, opt.Memory)
    }

    if (opt.Cpu != "") {
        _, err := exec.Command("docker", "container", "update", "--cpus", opt.Cpu, resp.ID).Output()
        if err != nil {
            panic(err)
        }
    }

    err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
    if err != nil {
        panic(err)
    }

    //fmt.Println(resp.ID)

    return resp.ID
}
