package main

import (
	"flag"
	"fmt"
	"lexportd/operations"
	"lexportd/utils"
	"log"
	"os"
	"strings"
)

var (
	sock = flag.String("sock", "", "path to LXD socket")
	out  = flag.String("out", "", "folder to write snapshots")
)

func main() {
	flag.Parse()
	if *sock == "" {
		log.Fatalln("-sock <path to socket> argument missing")
	}
	if *out == "" {
		w, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get current working directory: %v\n", err)
		}
		out = &w
	}

	client, err := utils.NewClient(*sock)
	if err != nil {
		log.Fatalf("failed to dial socket: %v\n", err)
	}

	// Snapshot containers
	log.Println("Snapshotting containers...")
	containerURLs, err := operations.ListContainers(client)
	if err != nil {
		log.Fatalf("failed to list containers: %v\n", err)
	}
	var snapshotNames []utils.SnapshotNamingScheme
	for _, containerURL := range containerURLs {
		res, err := operations.SnapshotContainer(client, containerURL)
		name := strings.Split(res, "/")
		if err != nil {
			log.Fatalf("failed to snapshot %s: %v\n", name[1], err)
		}
		snapshotNames = append(snapshotNames, utils.SnapshotNamingScheme{
			ContainerName: name[0],
			SnapshotName:  name[1],
		})
	}

	// Publish snapshots
	log.Println("Publishing snapshots...")
	for _, snapshot := range snapshotNames {
		if err = operations.PublishSnapshot(client, snapshot.ContainerName, snapshot.SnapshotName); err != nil {
			log.Fatalf("failed to publish image: %v\n", err)
		}
	}

	// Export and write snapshots
	images, err := operations.ListImages(client)
	if err != nil {
		log.Fatalf("failed to list images: %v", err)
	}
	for _, v := range images {
		fmt.Println(v)
		data, err := operations.ExportImage(client, v)
		if err != nil {
			log.Fatalf("failed to export images: %v", err)
		}
		name, err := operations.GetImageFilename(client, v)
		if err != nil {
			log.Fatalf("failed to get filename for %q: %v", v, err)
		}
		log.Printf("Exporting %s\n", name)
		n, err := utils.WriteSnapshot(*out, name, data)
		if err != nil {
			log.Fatalf("failed to write snapshot for %q: %v", name, err)
		}
		log.Printf("%s: wrote %d bytes", name, n)
	}
}
