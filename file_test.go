package main

import (
	"testing"
	"time"
)

func TestSortFilesByCreated(t *testing.T) {
	type args struct {
		files     []*File
		ascending bool
		expected  []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "sort files ascending",
			args: args{
				files: []*File{
					{
						Name:    "file1",
						Created: time.Now(),
					},
					{
						Name:    "file2",
						Created: time.Now().Add(time.Minute),
					},
				},
				ascending: true,
				expected:  []string{"file1", "file2"},
			},
		},
		{
			name: "sort files descending",
			args: args{
				files: []*File{
					{
						Name:    "file1",
						Created: time.Now(),
					},
					{
						Name:    "file2",
						Created: time.Now().Add(time.Minute),
					},
				},
				ascending: false,
				expected:  []string{"file2", "file1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortFilesByCreated(tt.args.files, tt.args.ascending)
			for i, file := range tt.args.files {
				if file.Name != tt.args.expected[i] {
					t.Error("files were not sorted correctly")
				}
			}
		})
	}
}
