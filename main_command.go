package main

import (
    //"os"
    "fmt"
    "context"

    "github.com/kawasaki006/mydocker/container"
    "github.com/kawasaki006/mydocker/cgroups/subsystems"
    "github.com/urfave/cli/v3"
    log "github.com/sirupsen/logrus"
)

var runCommand = &cli.Command{
    Name: "run",
    Usage: "run a container",
    Flags: []cli.Flag{
		&cli.BoolFlag{
            Name: "ti",
            Usage: "enable tty",
		},
		&cli.StringFlag{
            Name: "m",
            Usage: "memory limit, e.g.: -m 100m",
		},
		&cli.IntFlag{
            Name: "cpu",
            Usage: "cpu quota, e.g.: to limit cpu share to 50% -> -cpu 50",
		},
		&cli.StringFlag{
            Name: "cpuset",
            Usage: "cpuset limit, e.g.: -cpuset 2,4",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() < 1 {
            return fmt.Errorf("missing command to run")
        }
        cmdArray := cmd.Args().Slice()
        tty := cmd.Bool("ti")
        res := &subsystems.ResourceConfig{
            MemoryLimit: cmd.String("m"),
            CpuSet:      cmd.String("cpuset"),
            CpuCfsQuota: cmd.Int("cpu"),
        }
        Run(tty, cmdArray, res)
        return nil
    },
}

var initCommand = &cli.Command{
    Name: "init",
    Usage: "init container; for internal call ONLY",
    Action: func(ctx context.Context, cmd *cli.Command) error {
        log.Infof("init come on...")
        c :=  cmd.Args().Get(0)
        log.Infof("command: %s", c)
        err := container.RunContainerInitProcess(c, nil)
        return err
    },
}
