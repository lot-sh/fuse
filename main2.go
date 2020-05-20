// Hellofs implements a simple "hello world" file system.
// https://github.com/ipfs/go-ipfs/blob/master/fuse/readonly/readonly_unix.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	// "io/ioutil"
	"io"
	"net/http"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	// fuseversion "github.com/jbenet/go-fuse-version"
	_ "bazil.org/fuse/fs/fstestutil"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	
	mountpoint := flag.Arg(0)

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("helloworld"),
		fuse.Subtype("hellofs"),
	)

	// sysv, err := fuseversion.LocalFuseSystems()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "%s\n", err)
	// 	os.Exit(1)
	// }

	// for _, s := range sysv {
	// 	fmt.Printf("Fuse Version %s, %s, %s\n", s.FuseVersion, s.AgentVersion, s.AgentName)
	// }

	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = fs.Serve(c, FS{})
	if err != nil {
		log.Fatal(err)
	}
}

// FS implements the hello world file system.
type FS struct{}

func (FS) Root() (fs.Node, error) {
	return Dir{}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct{}

func (Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0o555
	return nil
}

func (Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "hello" {
		return NewVirtualFile("https://gist.githubusercontent.com/rubeniskov/e8aeec1bebf46696a9b9557e1a1bf936/raw/3fd978bcb6e89274409e5f9eb8ad1e436939c40f/semver-patcher.py"), nil
	}
	return nil, syscall.ENOENT
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "hello", Type: fuse.DT_File},
}

func (Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.




// func (File) ReadAll(ctx context.Context) ([]byte, error) {
// 	resp, err := http.Get("https://gist.githubusercontent.com/rubeniskov/e8aeec1bebf46696a9b9557e1a1bf936/raw/3fd978bcb6e89274409e5f9eb8ad1e436939c40f/semver-patcher.py")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("%s", body)

// 	return body, nil
// }



type VirtualFile struct{
	size int64
	body io.Reader
}


func NewVirtualFile(locator string) *VirtualFile {
	hresp, err := http.Get(locator)
	if err != nil {
		log.Fatal(err)
	}
	return &VirtualFile{ 
		hresp.ContentLength,
		hresp.Body,
	}
}

func (vf *VirtualFile) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0o555
	a.Size = uint64(vf.size)
	return nil
}

func (vf *VirtualFile) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	// Data has a capacity of Size
	buf := resp.Data[:int(req.Size)]
	n, err := io.ReadFull(vf.body, buf)
	resp.Data = buf[:n]
	switch err {
	case nil, io.EOF, io.ErrUnexpectedEOF:
	default:
		return err
	}
	resp.Data = resp.Data[:n]
	return nil // may be non-nil / not succeeded
}