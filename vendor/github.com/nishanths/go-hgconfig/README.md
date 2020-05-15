# [go-hgconfig](https://github.com/nishanths/go-hgconfig)

Read Mercurial configuration (`hg config`) values in Go.

Useful for finding out config values such as `ui.username`.

[![wercker status](https://app.wercker.com/status/7c522e9a74827fe67b2f89a8746bef0e/s "wercker status")](https://app.wercker.com/project/bykey/7c522e9a74827fe67b2f89a8746bef0e)
[![GoDoc](https://godoc.org/gopkg.in/nishanths/go-hgconfig.v1?status.svg)](https://godoc.org/gopkg.in/nishanths/go-hgconfig.v1)

# Install

````
go get gopkg.in/nishanths/go-hgconfig.v1
````

Import it as:

````
import gopkg.in/nishanths/go-hgconfig.v1
````

and refer to it as `hgconfig`.

# Test

Tests can be found in [`go_hgconfig_test.go`](https://github.com/nishanths/go-hgconfig/blob/master/go_hgconfig_test.go). To run tests, install [ginkgo](https://github.com/onsi/ginkgo) and [gomega](https://github.com/onsi/gomega), then run:

````
go test
````

# Usage

````go
package main

import (
    "fmt"
    "gopkg.in/nishanths/go-hgconfig.v1"
)

func main() {
    value, _ := hgconfig.Get("ui.username")
    fmt.Println(value)
}
````

An example can also be found at [`example/main.go`](https://github.com/nishanths/go-hgconfig/blob/master/example/main.go), and can be run using `go run example/main.go`.

# Functions

### Get

Get lets you read a `hg config` value by name.

````go
value, _, := hgconfig.Get("merge-tools.editmerge.premerge")
````

### Username

Username is a convenience function for getting the config value for `ui.username`. This is the same as calling `hgconfig.Get("ui.username")`.

````go
username, _, := hgconfig.Username()
````

# Contributing

Pull requests for new features are welcome!

# License 

go-hgconfig is licensed under the [MIT License](https://github.com/nishanths/go-hgconfig/blob/master/LICENSE)
