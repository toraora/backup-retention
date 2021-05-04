package main

import (
	"errors"
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
				Usage:    "Period of backup to handle (daily, weekly, etc). Leave blank to apply policy to snapshots",
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
			&cli.StringFlag{
				Name:     "backend",
				Usage:    "Which storage backend to use (local or s3)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "bucket",
				Value:    "",
				Usage:    "S3 bucket name (required when using s3 backend)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "dir",
				Value:    "",
				Usage:    "Local directory name (required when using local backend)",
				Required: false,
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

	var backend Backend
	backendStr := c.String("backend")
	switch backendStr {
	case "s3":
		bucket := c.String("bucket")
		if bucket == "" {
			return errors.New("bucket cannot be empty when s3 backend is selected")
		}
		backend, err = NewS3Backend(c.Context, bucket)
		if err != nil {
			return err
		}
	case "local":
		dir := c.String("dir")
		backend, err = NewLocalBackend(dir)
		if err != nil {
			return err
		}
	default:
		return errors.New("backend must be one of s3 or local")
	}

	return policy.CopySnapshotAndEnforce(backend, false)
}
