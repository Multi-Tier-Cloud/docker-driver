/* Copyright 2020 Multi-Tier-Cloud Development Team
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package docker_driver

import (
    "math"
    "io/ioutil"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "golang.org/x/net/context"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/go-connections/nat"
)

type DockerConfig struct {
    Name string
    Image string
    Port [2]string
    Cmd []string
    Memory int64        // in bytes   min is 4M   default is inf
    Cpu float64         // between 0.00 to 1.00*cores
    Network string
    Env []string
}

// image should be imagename:version
// hash should be user/image@sha256:digest
// official images should be library/imagename
func PullImage(image string) (string, error) {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return "", err
    }

    out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
    if err != nil {
        return "", err
    }
    defer out.Close()

    // Read until EOF sent to ensure proper transfer of image
    _, err = ioutil.ReadAll(out)
    if err != nil {
        return "", err
    }

    return "success", nil
}

func ListImages() ([]string, error) {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, err
    }

    images, err := cli.ImageList(ctx, types.ImageListOptions{})
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
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, err
    }

    containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
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
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return "", err
    }

    if err := cli.ContainerStop(ctx, cont, nil); err != nil {
        return "", err
    }

    return "success", nil
}

// deleting container
func DeleteContainer(cont string) (string, error) {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return "", err
    }

    if err := cli.ContainerRemove(ctx, cont, types.ContainerRemoveOptions{}); err != nil {
        return "", err
    }

    return "success", nil
}

// restarting container
func RestartContainer(cont string) (string, error) {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return "", err
    }

    if err := cli.ContainerRestart(ctx, cont, nil); err != nil {
        return "", err
    }

    return "success", nil
}

// resizing a container instance on the fly
func ResizeContainer(cont string, mem int64, cpu float64) (string, error) {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return "", err
    }

    _, err = cli.ContainerUpdate(ctx, cont, container.UpdateConfig{
        Resources: container.Resources{
            Memory: mem,
            NanoCPUs: int64(cpu*(math.Pow(10, 9))),
        },
    });
    if err != nil {
        return "", err
    }

    return "success", nil
}

// create and run container - interactive and detached set
// image (already pulled) should be imagename:version
// default/empty cmd is /bin/bash
func RunContainer(opt DockerConfig) (string, error) {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
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
        Resources: container.Resources{
            Memory: opt.Memory,
            NanoCPUs: int64(opt.Cpu*(math.Pow(10, 9))),
        },
    },
    nil, opt.Name)
    if err != nil {
        return "", err
    }

    err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
    if err != nil {
        return "", err
    }

    return resp.ID, nil
}
