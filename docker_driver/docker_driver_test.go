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
package docker_driver_test

import (
    "os"
    "testing"

    driver "github.com/Multi-Tier-Cloud/docker-driver/docker_driver"
)

const (
    testImage = "busybox"

    failTestImage = "thisImageNameShouldNotExist"
    failContID = "thisIDShouldNotExist"
)

func TestBuildImage(test *testing.T) {
    buildTestTarArchive := "build-test/test-image.tar"
    buildContext, err := os.Open(buildTestTarArchive)
    if err != nil {
        test.Fatalf("Open() failed with error:\n%v", err)
    }

    buildTestImage := "test-image"
    err = driver.BuildImage(buildContext, buildTestImage)
    if err != nil {
        test.Errorf("BuildImage() returned:\n%v", err)
    }
}

func TestPullImage(test *testing.T) {
    test.Run("PullImage-success", func(test *testing.T) {
        digest, err := driver.PullImage(testImage)
        if err != nil {
            test.Errorf("PullImage() returned:\n%v", err)
        }
        if digest == "" {
            test.Errorf("PullImage() returned empty digest")
        }
        test.Logf("PullImage() returned digest: %s", digest)
    })

    test.Run("PullImage-fail", func(test *testing.T) {
        _, err := driver.PullImage(failTestImage)
        if err == nil {
            test.Errorf("PullImage() succeeded with image (%s), expected it to fail", failTestImage)
        }
    })
}

func TestSaveImage(test *testing.T) {
    test.Run("SaveImage-success", func(test *testing.T) {
        _, err := driver.SaveImage(testImage)
        if err != nil {
            test.Errorf("SaveImage() returned:\n%v", err)
        }
    })

    test.Run("SaveImage-fail", func(test *testing.T) {
        _, err := driver.SaveImage(failTestImage)
        if err == nil {
            test.Errorf("SaveImage() succeeded with image (%s), expected it to fail", failTestImage)
        }
    })
}

func TestListImages(test *testing.T) {
    _, err := driver.ListImages()
    if err != nil {
        test.Errorf("ListImages() returned:\n%v", err)
    }
}

func TestLifecycle(test *testing.T) {
    opt := driver.DockerConfig{
        Name: "lifecycle_test",
        Image: "busybox",
        Port: [2]string{"4812", "4821"},
        Cmd: []string{"sleep", "300"},
        Memory: 10e+6,
        Cpu: 0.5,
    }

    containerID := ""

    // Run separate subtests for Run, Restart, Stop, and Delete
    test.Run("RunContainer", func(test *testing.T) {
        contID, err := driver.RunContainer(opt)
        if err != nil || contID == "" {
            test.Errorf("RunContainer() returned:\n%v", err)
        }

        containerID = contID
    })


    if containerID == "" {
        test.Fatalf("Skipping remaining sub-tests (RunContainer() may have failed)")
    }

    test.Run("ResizeContainer", func(test *testing.T) {
        _, err := driver.ResizeContainer(containerID, 20e+6, 0.5)
        if err != nil {
            test.Errorf("ResizeContainer() returned:\n%v", err)
        }
    })

    test.Run("RestartContainer", func(test *testing.T) {
        _, err := driver.RestartContainer(containerID)
        if err != nil {
            test.Errorf("RestartContainer() returned:\n%v", err)
        }
    })

    test.Run("StopContainer", func(test *testing.T) {
        _, err := driver.StopContainer(containerID)
        if err != nil {
            test.Errorf("StopContainer() returned:\n%v", err)
        }
    })

    test.Run("DeleteContainer", func(test *testing.T) {
        _, err := driver.DeleteContainer(containerID)
        if err != nil {
            test.Errorf("DeleteContainer() returned:\n%v", err)
        }
    })

}

func TestListRunningContainers(test *testing.T) {
    _, err := driver.ListRunningContainers()
    if err != nil {
        test.Errorf("ListRunningContainers() returned:\n%v", err)
    }
}

func TestResizeContainer(test *testing.T) {
    // Test failure case (success case covered in lifecycle test)
    _, err := driver.ResizeContainer(failContID, 10e+6, 0.7)
    if err == nil {
        test.Errorf("ResizeContainer() succeeded with container (%s), expected it to fail", failContID)
    }
}

func TestRestartContainer(test *testing.T) {
    // Test failure case (success case covered in lifecycle test)
    _, err := driver.RestartContainer(failContID)
    if err == nil {
        test.Errorf("RestartContainer() succeeded with container (%s), expected it to fail", failContID)
    }
}

func TestStopContainer(test *testing.T) {
    // Test failure case (success case covered in lifecycle test)
    _, err := driver.StopContainer(failContID)
    if err == nil {
        test.Errorf("StopContainer() succeeded with container (%s), expected it to fail", failContID)
    }
}

func TestDeleteContainer(test *testing.T) {
    // Test failure case (success case covered in lifecycle test)
    _, err := driver.DeleteContainer(failContID)
    if err == nil {
        test.Errorf("DeleteContainer() succeeded with container (%s), expected it to fail", failContID)
    }
}

