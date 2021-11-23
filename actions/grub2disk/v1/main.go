package main

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/tinkerbell/hub/actions/grub2disk/v1/pkg/grub"
)

func main() {
	fmt.Println("GRUB2Disk - GRUB streamer\n------------------------\n")
	grubInstallPath := os.Getenv("GRUB_INSTALL_PATH")
	grubBlockDevice := os.Getenv("GRUB_DISK")
	filesystemType := os.Getenv("FS_TYPE")

	if _, err := exec.Command("/bin/sh", "-c", "apk add grub grub-bios").Output(); err != nil {
		log.Fatal(fmt.Errorf("failed to install grub-install with error %s", err))
	}
	if err := grub.MountGrub(grubInstallPath, grubBlockDevice, filesystemType); err != nil {
		log.Fatal(err)
	}
	log.Infof("grub successfully written on [%s]", grubInstallPath)
}
