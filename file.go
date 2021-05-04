package main

import (
	"sort"
	"time"
)

// File contains file metadata
type File struct {
	Name    string
	Created time.Time
}

// SortFilesByCreated sorts the slice in either ascending or descending order of time
func SortFilesByCreated(files []*File, ascending bool) {
	sort.Slice(files, func(i int, j int) bool {
		return ascending == files[i].Created.Before(files[j].Created)
	})
}
