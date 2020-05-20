## Dependencies

### Linux

```shell
$ sudo apt-get install fuse
```

### MacOSX

```shell
$ brew cask install osxfuse

```

## Build

```shell
$ go get
$ go build -o lot-fuse main.go
```

```shell
$ export LOTPATH=/tmp/mountpoint
$ export LOTBIN=$LOTPATH/bin
$ export PATH=$PATH:$LOTBIN
$ mkdir -p $LOTBIN
```

```shell
$ umount $LOTBIN ; go build -o lot-fuse main2.go && ./lot-fuse -fuse.debug $LOTBIN
```

## References

- https://github.com/osxfuse/osxfuse/wiki/FAQ
- https://github.com/ipfs/go-ipfs/blob/master/docs/fuse.md
- https://www.cyberciti.biz/faq/reload-sysctl-conf-on-linux-using-sysctl/

```shell
# test scripts
diskutil unmount $LOTBIN ; go build -o lot-fuse main.go && ./lot-fuse -debug $LOTBIN


umount $LOTBIN ; go build -o lot-fuse main2.go && ./lot-fuse -fuse.debug $LOTBIN
```
