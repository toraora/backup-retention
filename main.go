package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	// ctx := context.Background()
	app := &cli.App{
		Name:   "retention",
		Usage:  "flexible retention policy enforcement",
		Action: Enforce,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "period",
				Value:    "",
				Usage:    "Period of backup to handle (snapshot, daily, weekly, etc.)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "mode",
				Value:    "count",
				Usage:    "TODO",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "num",
				Usage:    "Number of backups to keep for this retention period",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Enforce is the main entrypoint to the 'enforce' subcommand
func Enforce(c *cli.Context) error {
	policy, err := NewPolicy(c.String("period"), c.String("mode"), c.Int("num"))
	if err != nil {
		return err
	}

	b, err := NewS3Backend(c.Context, "tt-development-nathanw-test")
	if err != nil {
		return err
	}
	return policy.CopySnapshotAndEnforce(b, false)
}
