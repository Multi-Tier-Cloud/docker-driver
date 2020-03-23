package main

import (
    "fmt"

    driver "github.com/Multi-Tier-Cloud/docker-driver/docker_driver"
)

func main() {
    opt := driver.Docker_config{
        Image: "ubuntu",
        Port: [2]string{"4812", "4821"},
        Cmd: []string{},
        Memory: "500m",
        Cpu: "0.5",
    }

    fmt.Println(driver.ListRunningContainers())
    container, _ := driver.RunContainer(opt)
    fmt.Println(driver.ListRunningContainers())
    driver.StopContainer(container)
    driver.DeleteContainer(container)
    fmt.Println(driver.ListRunningContainers())

    driver.PullImage("hivanco/hello-world-server@sha256:3d4002fcaaa8c2a363c94ebec142447d3b04e0dc3e77954598fe97f285e0c37e")
}
