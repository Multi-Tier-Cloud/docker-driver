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
    container := driver.RunContainer(opt)
    fmt.Println(driver.ListRunningContainers())
    driver.StopContainer(container)
    driver.DeleteContainer(container)
    fmt.Println(driver.ListRunningContainers())
}
