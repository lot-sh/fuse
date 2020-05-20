// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program is the analogon of libfuse's hello.c, a a program that
// exposes a single file "file.txt" in the root directory.
package main

import (
	"context"
	"flag"
	"log"
	"syscall"
	"fmt"
  	"os"

	"github.com/hanwen/go-fuse/fs"
	"github.com/hanwen/go-fuse/fuse"
	fuseversion "github.com/jbenet/go-fuse-version"
)

type HelloRoot struct {
	fs.Inode
}

func (r *HelloRoot) OnAdd(ctx context.Context) {
	ch := r.NewPersistentInode(
		ctx, &fs.MemRegularFile{
			Data: []byte("file.txt"),
			Attr: fuse.Attr{
				Mode: 0644,
			},
		}, fs.StableAttr{Ino: 2})
	r.AddChild("file.txt", ch, false)
}

func (r *HelloRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0755
	return 0
}

var _ = (fs.NodeGetattrer)((*HelloRoot)(nil))
var _ = (fs.NodeOnAdder)((*HelloRoot)(nil))

func main() {
	sysv, err := fuseversion.LocalFuseSystems()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	for _, s := range sysv {
		fmt.Printf("Fuse Version %s, %s, %s\n", s.FuseVersion, s.AgentVersion, s.AgentName)
	}

	debug := flag.Bool("debug", false, "print debug data")
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello MOUNTPOINT")
	}
	opts := &fs.Options{}
	opts.Debug = *debug

	fmt.Printf("Mounting directory %s\n", flag.Arg(0))
	server, err := fs.Mount(flag.Arg(0), &HelloRoot{}, opts)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Wait()
}