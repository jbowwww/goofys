// Copyright 2015 Ka-Hing Cheung
// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"time"

	"github.com/codegangsta/cli"
)

// Set up custom help text for goofys; in particular the usage section.
func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[global options]{{end}} bucket mountpoint
   {{if .Version}}
VERSION:
   {{.Version}}
   {{end}}{{if len .Authors}}
AUTHOR(S):
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}
`
}

func newApp() (app *cli.App) {
	app = &cli.App{
		Name:     "goofys",
		Version:  "0.0.1",
		Usage:    "Mount an S3 bucket locally",
		HideHelp: true,
		Writer:   os.Stderr,
		Flags: []cli.Flag{

			cli.BoolFlag{
				Name:  "help, h",
				Usage: "Print this help text and exit successfuly.",
			},

			/////////////////////////
			// File system
			/////////////////////////

			cli.StringSliceFlag{
				Name:  "o",
				Usage: "Additional system-specific mount options. Be careful!",
			},

			cli.IntFlag{
				Name:  "dir-mode",
				Value: 0755,
				Usage: "Permissions bits for directories. (default: 0755)",
			},

			cli.IntFlag{
				Name:  "file-mode",
				Value: 0644,
				Usage: "Permission bits for files (default: 0644)",
			},

			cli.IntFlag{
				Name:  "uid",
				Value: -1,
				Usage: "UID owner of all inodes.",
			},

			cli.IntFlag{
				Name:  "gid",
				Value: -1,
				Usage: "GID owner of all inodes.",
			},

			cli.BoolFlag{
				Name: "implicit-dirs",
				Usage: "Implicitly define directories based on content. See" +
					"docs/semantics.md",
			},

			cli.StringFlag{
				Name: "storage-class",
				Value: "STANDARD",
				Usage: "The type of storage to use when writing objects." +
					" Possible values: REDUCED_REDUNDANCY, STANDARD (default), STANDARD_IA.",
			},

			/////////////////////////
			// Goofys
			/////////////////////////

			cli.Float64Flag{
				Name:  "limit-bytes-per-sec",
				Value: -1,
				Usage: "Bandwidth limit for reading data, measured over a 30-second " +
					"window. (use -1 for no limit)",
			},

			cli.Float64Flag{
				Name:  "limit-ops-per-sec",
				Value: 5.0,
				Usage: "Operations per second limit, measured over a 30-second window " +
					"(use -1 for no limit)",
			},

			/////////////////////////
			// Tuning
			/////////////////////////

			cli.DurationFlag{
				Name:  "stat-cache-ttl",
				Value: time.Minute,
				Usage: "How long to cache StatObject results and inode attributes.",
			},

			cli.DurationFlag{
				Name:  "type-cache-ttl",
				Value: time.Minute,
				Usage: "How long to cache name -> file/dir mappings in directory " +
					"inodes.",
			},

			/////////////////////////
			// Debugging
			/////////////////////////

			cli.BoolFlag{
				Name:  "debug_fuse",
				Usage: "Enable fuse-related debugging output.",
			},

			cli.BoolFlag{
				Name:  "debug_invariants",
				Usage: "Panic when internal invariants are violated.",
			},

			cli.BoolFlag{
				Name:  "debug_s3",
				Usage: "Enable S3-related debugging output.",
			},
		},
	}

	return
}

type flagStorage struct {
	// File system
	MountOptions map[string]string
	DirMode      os.FileMode
	FileMode     os.FileMode
	Uid          uint32
	Gid          uint32
	ImplicitDirs bool
	StorageClass string

	// Goofys
	EgressBandwidthLimitBytesPerSecond float64
	OpRateLimitHz                      float64

	// Tuning
	StatCacheTTL time.Duration
	TypeCacheTTL time.Duration

	// Debugging
	DebugFuse       bool
	DebugInvariants bool
	DebugS3         bool
}

// Add the flags accepted by run to the supplied flag set, returning the
// variables into which the flags will parse.
func populateFlags(c *cli.Context) (flags *flagStorage) {
	flags = &flagStorage{
		// File system
		MountOptions: make(map[string]string),
		DirMode:      os.FileMode(c.Int("dir-mode")),
		FileMode:     os.FileMode(c.Int("file-mode")),
		Uid:          uint32(c.Int("uid")),
		Gid:          uint32(c.Int("gid")),

		// Goofys,
		EgressBandwidthLimitBytesPerSecond: c.Float64("limit-bytes-per-sec"),
		OpRateLimitHz:                      c.Float64("limit-ops-per-sec"),

		// Tuning,
		StatCacheTTL: c.Duration("stat-cache-ttl"),
		TypeCacheTTL: c.Duration("type-cache-ttl"),
		ImplicitDirs: c.Bool("implicit-dirs"),
		StorageClass: c.String("storage-class"),

		// Debugging,
		DebugFuse:       c.Bool("debug_fuse"),
		DebugInvariants: c.Bool("debug_invariants"),
		DebugS3:         c.Bool("debug_s3"),
	}

	/*
		// Handle the repeated "-o" flag.
		for _, o := range c.StringSlice("o") {
			mountpkg.ParseOptions(flags.MountOptions, o)
		}
	*/
	return
}
