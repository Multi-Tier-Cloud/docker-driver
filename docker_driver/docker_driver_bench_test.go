/* Copyright 2020 PhysarumSM Development Team
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
package docker_driver_test

import (
    "fmt"
    "testing"

    driver "github.com/PhysarumSM/docker-driver/docker_driver"
)

var containerIDs []string
var opt = driver.DockerConfig{
    Image: "busybox",
    Network: "host",
    Cmd: []string{"sleep", "3"},
    Memory: 10e+6,
    Cpu: 0.5,
}

func removeContainers() {
    // Use driver itself to do cleanup so that if cleanup fails, we
    // catch another potential bug.
    for _, contID := range containerIDs {
        _, err := driver.DeleteContainer(contID)
        if err != nil {
            fmt.Printf("ERROR: Unable to delete container %s\n", contID)
        }
    }
}

func BenchmarkRunContainer(bench *testing.B) {
    // Register clean-up function
    // NOTE: If using go v1.14, can use bench.Cleanup(removeContainers)
    defer removeContainers()

    // Encapsulate actual test in Run() to avoid timing clean-up operations
    bench.Run("RunContainer", func(bench *testing.B) {
        for i := 0; i < bench.N; i++ {
            contID, err := driver.RunContainer(opt)
            if err != nil || contID == "" {
                // NOTE: Yes, this error check may impact performance... but the
                //       bottleneck is likely Docker daemon and container create.
                bench.Errorf("RunContainer() returned:\n%v", err)
                continue
            }
            containerIDs = append(containerIDs, contID)
        }
    })
}
