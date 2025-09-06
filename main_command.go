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
		&cli.BoolFlag{
            Name: "d",
            Usage: "detach container",
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
        &cli.StringFlag{
            Name: "v",
            Usage: "volume, e.g.: -v /etc/conf:/etc/conf",
        },
        &cli.StringFlag{
            Name: "name",
            Usage: "container name, e.g: -name mycontainer",
        },
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() < 1 {
            return fmt.Errorf("missing command to run")
        }

        cmdArray := cmd.Args().Slice()
        // image name should be the first arg
        imageName := cmdArray[0]
        // extra args
        cmdArray = cmdArray[1:]

        // flags
        tty := cmd.Bool("ti")
        detach := cmd.Bool("d")
        if tty && detach {
            return fmt.Errorf("ti and d cannot be provided simultaneously")
        }
        res := &subsystems.ResourceConfig{
            MemoryLimit: cmd.String("m"),
            CpuSet:      cmd.String("cpuset"),
            CpuCfsQuota: cmd.Int("cpu"),
        }
        volume := cmd.String("v")
        containerName := cmd.String("name")

        Run(tty, cmdArray, res, volume, containerName, imageName)
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

var listCommand = &cli.Command{
    Name: "ps",
    Usage: "list all containers",
    Action: func(ctx context.Context, cmd *cli.Command) error {
        log.Infof("listing all containers...")
        listContainers()
        return nil
    },
}
