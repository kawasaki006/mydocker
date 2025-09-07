package main

import (
	"os"
	"context"
   
	log "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v3"
)

func main() {
    cmd := &cli.Command{
        Name: "mydocker",
		Usage: "docker toy",
		Commands: []*cli.Command{
			runCommand,
            initCommand,
            listCommand,
            logCommand,
            execCommand,
        },
        Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
            // log as json
            log.SetFormatter(&log.JSONFormatter{})
            // output to stdout
            log.SetOutput(os.Stdout)
            return nil, nil
        },
    }

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

