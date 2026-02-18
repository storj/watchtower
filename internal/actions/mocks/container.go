package mocks

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/containrrr/watchtower/pkg/container"
	wt "github.com/containrrr/watchtower/pkg/types"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"
)

// CreateMockContainer creates a container substitute valid for testing
func CreateMockContainer(id string, name string, img string, created time.Time) wt.Container {
	content := dockerContainer.InspectResponse{
		ContainerJSONBase: &dockerContainer.ContainerJSONBase{
			ID:      id,
			Image:   img,
			Name:    name,
			Created: created.String(),
			HostConfig: &dockerContainer.HostConfig{
				PortBindings: map[nat.Port][]nat.PortBinding{},
			},
		},
		Config: &dockerContainer.Config{
			Image:        img,
			Labels:       make(map[string]string),
			ExposedPorts: map[nat.Port]struct{}{},
		},
	}
	return container.NewContainer(
		&content,
		CreateMockImageInfo(img),
	)
}

// CreateMockImageInfo returns a mock image info struct based on the passed image
func CreateMockImageInfo(img string) *image.InspectResponse {
	return &image.InspectResponse{
		ID: img,
		RepoDigests: []string{
			img,
		},
	}
}

// CreateMockContainerWithImageInfo should only be used for testing
func CreateMockContainerWithImageInfo(id string, name string, img string, created time.Time, imageInfo image.InspectResponse) wt.Container {
	return CreateMockContainerWithImageInfoP(id, name, img, created, &imageInfo)
}

// CreateMockContainerWithImageInfoP should only be used for testing
func CreateMockContainerWithImageInfoP(id string, name string, img string, created time.Time, imageInfo *image.InspectResponse) wt.Container {
	content := dockerContainer.InspectResponse{
		ContainerJSONBase: &dockerContainer.ContainerJSONBase{
			ID:      id,
			Image:   img,
			Name:    name,
			Created: created.String(),
		},
		Config: &dockerContainer.Config{
			Image:  img,
			Labels: make(map[string]string),
		},
	}
	return container.NewContainer(
		&content,
		imageInfo,
	)
}

// CreateMockContainerWithDigest should only be used for testing
func CreateMockContainerWithDigest(id string, name string, img string, created time.Time, digest string) wt.Container {
	c := CreateMockContainer(id, name, img, created)
	c.ImageInfo().RepoDigests = []string{digest}
	return c
}

// CreateMockContainerWithConfig creates a container substitute valid for testing
func CreateMockContainerWithConfig(id string, name string, img string, running bool, restarting bool, created time.Time, config *dockerContainer.Config) wt.Container {
	content := dockerContainer.InspectResponse{
		ContainerJSONBase: &dockerContainer.ContainerJSONBase{
			ID:    id,
			Image: img,
			Name:  name,
			State: &dockerContainer.State{
				Running:    running,
				Restarting: restarting,
			},
			Created: created.String(),
			HostConfig: &dockerContainer.HostConfig{
				PortBindings: map[nat.Port][]nat.PortBinding{},
			},
		},
		Config: config,
	}
	return container.NewContainer(
		&content,
		CreateMockImageInfo(img),
	)
}

// CreateContainerForProgress creates a container substitute for tracking session/update progress
func CreateContainerForProgress(index int, idPrefix int, nameFormat string) (wt.Container, wt.ImageID) {
	indexStr := strconv.Itoa(idPrefix + index)
	mockID := indexStr + strings.Repeat("0", 61-len(indexStr))
	contID := "c79" + mockID
	contName := fmt.Sprintf(nameFormat, index+1)
	oldImgID := "01d" + mockID
	newImgID := "d0a" + mockID
	imageName := fmt.Sprintf("mock/%s:latest", contName)
	config := &dockerContainer.Config{
		Image: imageName,
	}
	c := CreateMockContainerWithConfig(contID, contName, oldImgID, true, false, time.Now(), config)
	return c, wt.ImageID(newImgID)
}

// CreateMockContainerWithLinks should only be used for testing
func CreateMockContainerWithLinks(id string, name string, img string, created time.Time, links []string, imageInfo *image.InspectResponse) wt.Container {
	content := dockerContainer.InspectResponse{
		ContainerJSONBase: &dockerContainer.ContainerJSONBase{
			ID:      id,
			Image:   img,
			Name:    name,
			Created: created.String(),
			HostConfig: &dockerContainer.HostConfig{
				Links: links,
			},
		},
		Config: &dockerContainer.Config{
			Image:  img,
			Labels: make(map[string]string),
		},
	}
	return container.NewContainer(
		&content,
		imageInfo,
	)
}
