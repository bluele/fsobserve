# fsobserve

fsobserve is yet another observer for file system on multiplatform.

Currently supports Windows, OS X, Linux, Android, FreeBSD 4.1. Furthermore, this project uses [fsnotify](https://github.com/go-fsnotify/fsnotify) to capture notification for file system, so supporting os is depend on [fsnotify](https://github.com/go-fsnotify/fsnotify).

# Usage

```
// current directory have changed, simply execute "ls -al".
$ fsobserve -c "ls -al"

// specific observing directory (default current directory).
$ fsobserve -c "ls -al" -d "/path/to/dir"

// specific file patterns (default all files).
$ fsobserve -c "ls -al" -p "*.go *.md"

// specific observing interval (default 3 second).
$ fsobserve -c "ls -al" -i 10s
```

# Install

```
$ go install github.com/bluele/fsobserve
```

# Thanks

[fsnotify](https://github.com/go-fsnotify/fsnotify)

# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>
