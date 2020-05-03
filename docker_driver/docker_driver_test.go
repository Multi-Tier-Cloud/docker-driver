package docker_driver

import (
    "fmt"
    "testing"
)

func TestDockerDriver(test *testing.T) {
    opt := Docker_config{
        Image: "ubuntu",
        Port: [2]string{"4812", "4821"},
        Cmd: []string{},
        Memory: "500m",
        Cpu: "0.5",
    }

    fmt.Println(ListRunningContainers())
    container, _ := RunContainer(opt)
    fmt.Println(ListRunningContainers())
    StopContainer(container)
    DeleteContainer(container)
    fmt.Println(ListRunningContainers())

    PullImage("hivanco/hello-world-server@sha256:3d4002fcaaa8c2a363c94ebec142447d3b04e0dc3e77954598fe97f285e0c37e")
}
