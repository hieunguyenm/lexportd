package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lexportd/utils"
	"log"
	"net/http/httputil"
	"strings"
	"time"
)

// ListContainers lists all LXC containers.
func ListContainers(client *httputil.ClientConn) ([]string, error) {
	body, err := utils.Do(client, "GET", "/1.0/containers", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}
	var res utils.ListContainerResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ListContainers response: %v", err)
	}
	return res.Metadata, nil
}

// SnapshotContainer takes a snapshot of a container with the following
// naming scheme: containerName_YYYY-MM-DD-HH-MM-SS.
func SnapshotContainer(client *httputil.ClientConn, url string) (string, error) {
	containerName := strings.Split(url, "/")[3]
	now := time.Now().Format("2006-01-02-15-04-05")
	snapshotName := fmt.Sprintf("%s_%s", containerName, now)
	req, err := json.Marshal(&utils.SnapshotRequest{
		Name:     snapshotName,
		Stateful: false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal SnapshotRequest: %v", err)
	}
	body, err := utils.Do(client, "POST", url+"/snapshots", bytes.NewReader(req))
	if err != nil {
		return "", fmt.Errorf("failed to snapshot container %q: %v", containerName, err)
	}
	var m utils.SnapshotResponse
	if err := json.Unmarshal(body, &m); err != nil {
		return "", fmt.Errorf("failed to unmarshal SnapshotResponse: %v", err)
	}
	if m.ErrorCode != 0 {
		return "", fmt.Errorf("failed to snapshot container %q: %v", snapshotName, m)
	}
	if err := waitBackgroundOperation(client, containerName, m.Metadata.ID); err != nil {
		return "", fmt.Errorf("failed to check background snapshot operations: %v", err)
	}
	return fmt.Sprintf("%s/%s", containerName, snapshotName), nil
}

// PublishSnapshot publishes a private image of a container snapshot.
func PublishSnapshot(client *httputil.ClientConn, containerName, snapshotName string) error {
	req, err := json.Marshal(&utils.PublishRequest{
		Filename: snapshotName + ".tar.xz",
		Public:   false,
		Aliases: []utils.PublishRequestAliases{
			utils.PublishRequestAliases{
				Name:        snapshotName,
				Description: "",
			}},
		Source: utils.PublishRequestSource{
			Type: "snapshot",
			Name: fmt.Sprintf("%s/%s", containerName, snapshotName),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal PublishRequest: %v", err)
	}
	body, err := utils.Do(client, "POST", "/1.0/images", bytes.NewReader(req))
	if err != nil {
		return fmt.Errorf("failed to publish image for %s: %v", containerName, err)
	}
	var r utils.SnapshotResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return fmt.Errorf("failed to unmarshal PublishResponse for %s: %v", containerName, err)
	}
	if err := waitBackgroundOperation(client, containerName, r.Metadata.ID); err != nil {
		return fmt.Errorf("failed to check background publish operations: %v", err)
	}
	return nil
}

// ListImages lists all container images.
func ListImages(client *httputil.ClientConn) ([]string, error) {
	body, err := utils.Do(client, "GET", "/1.0/images", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %v", err)
	}
	var res utils.ListContainerResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ListImages: %v", err)
	}
	return res.Metadata, nil
}

// ExportImage exports the image as a compressed archive.
func ExportImage(client *httputil.ClientConn, fingerprintURL string) ([]byte, error) {
	body, err := utils.Do(client, "GET", fingerprintURL+"/export", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to export image: %v", err)
	}
	return body, nil
}

// GetImageFilename returns the filename for the image.
func GetImageFilename(client *httputil.ClientConn, fingerprintURL string) (string, error) {
	body, err := utils.Do(client, "GET", fingerprintURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get image properties: %v", err)
	}
	var res utils.ImagePropertiesResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("failed to unmarshal ImagePropertiesResponse: %v", err)
	}
	return res.Metadata.Filename, nil
}

func waitBackgroundOperation(client *httputil.ClientConn, name, uuid string) error {
	for {
		b, err := utils.Do(client, "GET", "/1.0/operations", nil)
		if err != nil {
			return fmt.Errorf("failed to list operations: %v", err)
		}
		var s utils.ListOperationsResponse
		if err := json.Unmarshal(b, &s); err != nil {
			return fmt.Errorf("failed to unmarshal for waitBackgroundOperation: %v", err)
		}
		if len(s.Metadata.Running) == 0 {
			return nil
		}
		for _, v := range s.Metadata.Running {
			if !strings.Contains(v, uuid) {
				return nil
			}
		}
		log.Printf("%s: Operation in progress\n", name)
		time.Sleep(10 * time.Second)
	}
}
