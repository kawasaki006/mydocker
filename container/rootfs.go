package container

import (
    "os"
    "os/exec"

    log "github.com/sirupsen/logrus"
    "github.com/kawasaki006/mydocker/utils"
)

func NewWorkspace(containerId, imageName, volume string) {
    createLower(containerId, imageName)
    createUpper(containerId)
    createWork(containerId)
    createMerged(containerId)
    mountOverlay2(containerId)
}

func DeleteWorkspace(containerId string, volume string) {
    log.Infof("cleaning up workspace...")
    // unmount overlay fs
    unmountOverlay2(containerId)
    // cleanup folders
    deleteDirs(containerId)
}

func createLower(containerId string, imageName string) {
    lowerPath := utils.GetLower(containerId)
    imagePath := utils.GetImage(imageName)

    exist, err := utils.PathExists(lowerPath)
    if err != nil {
        log.Errorf("error checking if lower dir exists")
    }

    // create lower dir if not exist and untar image file
    if !exist {
        if err = os.MkdirAll(lowerPath, 0777); err != nil {
            log.Errorf("error making lower dir at [%s]: %v", lowerPath, err)
        }
        log.Infof("untarring image...")
        if _, err = exec.Command("tar", "-xvf", imagePath, "-C", lowerPath).CombinedOutput(); err != nil {
            log.Errorf("error untarring image at %s: %v", imagePath, err)
        }
    }
}

func createUpper(containerId string) {
    upperPath := utils.GetUpper(containerId)

    if err := os.Mkdir(upperPath, 0777); err != nil {
        log.Errorf("error making upper dir at [%s]: %v", upperPath, err)
    }
}

func createWork(containerId string) {
    workPath := utils.GetWork(containerId)

    if err := os.Mkdir(workPath, 0777); err != nil {
        log.Errorf("error making work dir at [%s]: %v", workPath, err)
    }
}

func createMerged(containerId string) {
    mergedPath := utils.GetMerged(containerId)

    if err := os.Mkdir(mergedPath, 0777); err != nil {
        log.Errorf("error making merged dir at [%s]: %v", mergedPath, err)
    }
}

func mountOverlay2(containerId string) {
    dirs := utils.GetOverlay2Path(utils.GetLower(containerId), utils.GetUpper(containerId), utils.GetWork(containerId))
    mergedPath := utils.GetMerged(containerId)
    cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mergedPath)
    log.Infof("mounting overlay2: [%s]...", cmd.String())
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("error mounting overlayfs: %v", err)
    }
}

func unmountOverlay2(containerId string) {
    mountPath := utils.GetMerged(containerId)
    cmd := exec.Command("umount", mountPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    log.Infof("unmounting overlay2...")
    if err := cmd.Run(); err != nil {
        log.Errorf("error unmounting overlay2: %v", err)
    }
}

func deleteDirs(containerId string) {
    dirs := []string{
        utils.GetMerged(containerId),
        utils.GetUpper(containerId),
        utils.GetWork(containerId),
        utils.GetLower(containerId),
        utils.GetRoot(containerId),
    }

    for _, dir := range dirs {
        if err := os.RemoveAll(dir); err != nil {
            log.Errorf("error removing dir at [%s]: %v", dir, err)
        }
    }
}
