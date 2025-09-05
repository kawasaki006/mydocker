package container

import (
    "os"
    "os/exec"
    "fmt"
    "syscall"
    "io"
    "strings"
    "path/filepath"

    log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {
    log.Infof("running command: %s", command)

    cmdArray := readClientCommand()
    if cmdArray == nil || len(cmdArray) == 0 {
        log.Errorf("run container init proc is trying to read user command, but command is empty")
    }

    // mount fs
    setUpMount()

    // find full path of command in new container rootfs
    path, err := exec.LookPath(cmdArray[0])
    log.Infof("command is: %s", cmdArray[0])
    if err != nil {
        log.Errorf("Invalid command name, exec look path error: %v", err)
        return err
    }
    if err = syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
        log.Errorf("Error executing command [%s]: %v", command, err)
    }
    return nil
}

func readClientCommand() []string {
    // retrieve read pipe
    readPipe := os.NewFile(uintptr(3), "pipe")
    defer readPipe.Close()
    msg, err := io.ReadAll(readPipe)
    if err != nil {
        log.Errorf("Error reading client command in init: %v", err)
        return nil
    }
    msgStr := string(msg)
    return strings.Split(msgStr, " ")
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
