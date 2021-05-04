package main

import (
	"errors"
	"fmt"
)

type period string

const (
	snapshot period = ""
	daily           = "daily"
	weekly          = "weekly"
	monthly         = "monthly"
	yearly          = "yearly"
)

type mode string

const (
	count    mode = "count"
	datetime      = "datetime"
)

// Policy represents the frequency and number of backups to keep
type Policy struct {
	period period
	mode   mode
	num    int
}

// NewPolicy validates freeform input and returns a Policy
func NewPolicy(periodStr string, modeStr string, num int) (*Policy, error) {
	p := (period)(periodStr)
	switch p {
	case snapshot, daily, weekly, monthly, yearly:
	default:
		return nil, errors.New("invalid period specified")
	}

	m := (mode)(modeStr)
	switch m {
	case count, datetime:
	default:
		return nil, errors.New("invalid mode specified")
	}

	return &Policy{
		period: p,
		mode:   m,
		num:    num,
	}, nil
}

// CopySnapshotAndEnforce copes the latest snapshot and then enforces the given Policy
func (p *Policy) CopySnapshotAndEnforce(b Backend, dryRun bool) error {
	// copy latest
	if p.period != snapshot {
		files, err := b.ListFiles("")
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return errors.New("could not find latest snapshot")
		}
		SortFilesByCreated(files, false)
		latestFile := files[0]
		destinationName := fmt.Sprintf("%s/%s", p.period, latestFile.Name)
		err = b.CopyFile(latestFile.Name, destinationName)
		if err != nil {
			return err
		}
	}

	// delete old backups
	files, err := b.ListFiles(string(p.period))
	if err != nil {
		return err
	}
	if len(files) < p.num {
		return nil
	}
	SortFilesByCreated(files, false)
	for _, file := range files[p.num:] {
		err = b.DeleteFile(file.Name)
		if err != nil {
			return err
		}
	}

	return nil
}
