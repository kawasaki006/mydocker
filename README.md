# Minimal Container Runtime in Go

This project is a lightweight Docker-like container runtime implemented in Go, featuring:
- Linux namespace–based process isolation (UTS, PID, Mount, Network, IPC)
- Root filesystem isolation using `pivot_root`
- OverlayFS-based container filesystem (lower / upper / work / merged)
- Resource limitation via cgroups v2 (CPU, memory, cpuset)
- Volume mount via bind mount
- `exec` support via `setns`
- Simple metadata management and logging
- CLI interface similar to Docker (`run`, `ps`, `logs`, `exec`)

---

## Project Structure

- `cgroups`: cgroup v2 resource management
- `container`: container lifecycle, init process, rootfs, volume
- `nsenter`: exec implementation based on setns
- `utils`: filesystem and path helpers
- `run.go`: container creation flow
- `exec.go`: docker exec–like functionality
- `list.go`: container listing (ps)
- `logs.go`: container logs
- `main.go`: CLI entry

---

## 1) Run Container

### Container Creation (`run`)
- Creates a new process with isolated namespaces
- Sets up container root filesystem using OverlayFS
- Applies CPU/memory limits via cgroups
- Executes user command inside container init process
- Records container metadata under `/var/lib/mydocker/containers`

### Container Init (`init`)
- Runs as PID 1 inside the container
- Switches root filesystem using `pivot_root`
- Mounts `/proc` and `/dev`
- Executes target command via `execve`

---

## 2) Resource Management (cgroups v2)

The runtime applies resource limits through cgroups v2:

- Memory: `memory.max`
- CPU quota: `cpu.max`
- CPU affinity: `cpuset.cpus`

Each resource type implements a common interface and is managed
by a central `CgroupManager`.

---

## 3) Exec Into Running Container

`mydocker exec` uses `setns` to enter the namespaces of a running container:

- Joins IPC, UTS, NET, PID, and MNT namespaces
- Executes commands inside the container environment

This mirrors the core mechanism behind `docker exec`.

---

## CLI Usage

Run a container:
```bash
sudo mydocker run -ti busybox sh
```

Run with resource controls:
```bash
sudo mydocker run -ti -m 100m -cpu 50 busybox sh
```

Mount a volume:
```bash
sudo mydocker run -ti -v /host/data:/data busybox sh
```

List containers:
```bash
sudo mydocker ps
```

View logs:
```bash
sudo mydocker logs <container-id>
```

Exec into a container:
```bash
sudo mydocker exec <container-id> sh
```