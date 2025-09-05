package container

import (
    "os"
    "os/exec"
    "path"
    "strings"
    "fmt"

    log "github.com/sirupsen/logrus"
)

func mountVolume(mountPath, hostPath, containerPath string) {
    if err := os.Mkdir(hostPath, 0777); err != nil {
        log.Errorf("error creating host dir at [%s]: %v", hostPath, err)
    }

    containerPathInHost := path.Join(mountPath, containerPath)
    if err := os.Mkdir(containerPathInHost, 0777); err != nil {
        log.Errorf("error creating container dir at [%s], %v", containerPathInHost, err)
    }

    // bind mount host dir to container dir
    cmd := exec.Command("mount", "-o", "bind", hostPath, containerPathInHost)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("error bind mounting volume [%s] to [%s]: %v", hostPath, containerPathInHost, err)
    }
}

func unmountVolume(mountPath, containerPath string) {
    containerPathInHost := path.Join(mountPath, containerPath)
    cmd := exec.Command("umount", containerPathInHost)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("error unmounting volume: %v", err)
    }
}

func volumeExtract(volume string) (sourcePath, destPath string, err error) {
    parts := strings.Split(volume, ":")
    if len(parts) != 2 {
        return "", "", fmt.Errorf("invalid volume path [%s]", volume)
    }

    sourcePath, destPath = parts[0], parts[1]
    if sourcePath == "" || destPath == "" {
        return "", "", fmt.Errorf("invalid volume path [%s]", volume)
    }

    return sourcePath, destPath, nil
}
