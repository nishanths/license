# go-appdir

[![GoDoc](https://godoc.org/github.com/ProtonMail/go-appdir?status.svg)](https://godoc.org/github.com/ProtonMail/go-appdir)

Minimalistic Go package to get application directories such as config and cache.

Platform | Windows | [Linux/BSDs] | [macOS]
-------- | ------- | ------------------------------------------------------------------------------------------ | -----
User-specific config | `%APPDATA%` (`C:\Users\%USERNAME%\AppData\Roaming`) | `$XDG_CONFIG_HOME` (`$HOME/.config`) | `$HOME/Library/Application Support`
User-specific cache | `%LOCALAPPDATA%` (`C:\Users\%USERNAME%\AppData\Local`) | `$XDG_CACHE_HOME` (`$HOME/.cache`) | `$HOME/Library/Caches`
User-specific logs | `%LOCALAPPDATA%` (`C:\Users\%USERNAME%\AppData\Local`) | `$XDG_STATE_HOME` (`$HOME/.local/state`) | `$HOME/Library/Logs`
User-specific data | `%LOCALAPPDATA%` (`C:\Users\%USERNAME%\AppData\Local`) | `$XDG_DATA_HOME` (`$HOME/.local/share`) | `$HOME/Library/Application Support`

[Linux/BSDs]: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
[macOS]: https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/FileSystemProgrammingGuide/FileSystemOverview/FileSystemOverview.html#//apple_ref/doc/uid/TP40010672-CH2-SW1

Inspired by [`configdir`](https://github.com/shibukawa/configdir).

## Usage

```go
package main

import (
	"os"
	"path/filepath"

	"github.com/ProtonMail/go-appdir"
)

func main() {
	// Get directories for our app
	dirs := appdir.New("my-awesome-app")

	// Get user-specific config dir
	p := dirs.UserConfig()

	// Create our app config dir
	if err := os.MkdirAll(p, 0755); err != nil {
		panic(err)
	}

	// Now we can use it
	f, err := os.Create(filepath.Join(p, "config-file"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write([]byte("<3"))
}
```

## License

MIT
