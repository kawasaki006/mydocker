package container

import (
    "os"
    "fmt"
    "syscall"
    "path/filepath"

    log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {
    log.Infof("running command: %s", command)

    // mount fs
    setUpMount()

    argv := []string{command}
    if err := syscall.Exec(command, argv, os.Environ()); err != nil {
        log.Errorf("Error executing command [%s]: %v", command, err)
    }
    return nil
}

func setUpMount() {
    // get pwd
    pwd, err := os.Getwd()
    if err != nil {
        log.Errorf("Error getting current directory %v", err)
    }
    log.Infof("Current working directory: %s", pwd)
    
    // change propagation type to private (recursive) to avoid mount leakage
    syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

    // pivot root
    err = pivotRoot(pwd)
    if err != nil {
        log.Errorf("Error pivoting root: %v", err)
        return
    }

    // mount proc
    defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
    syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

    // mount dev
    syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
    // remount root so that old root and new root are not in the same fs
    if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
        return fmt.Errorf("Error mounting root: %w", err)
    }
    // create tmp folder to mount old root
    pivotDir := filepath.Join(root, ".pivot_root")
    if err := os.Mkdir(pivotDir, 0777); err != nil {
        return fmt.Errorf("Error creating tmp folder .pivotDir: %w", err)
    }
    // pivot root from host rootfs to new rootfs
    if err := syscall.PivotRoot(root, pivotDir); err != nil {
        return fmt.Errorf("Error pivoting to new rootfs: %w", err)
    }
    // change working directory to new root
    if err := syscall.Chdir("/"); err != nil {
        return fmt.Errorf("Error changing working directory to new root: %w", err)
    }
    
    // unmount old root
    pivotDir = filepath.Join("/", ".pivot_root")
    if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
        return fmt.Errorf("Error unmounting old root: %w", err)
    }
    // delete tmp folder
    return os.Remove(pivotDir)

}
