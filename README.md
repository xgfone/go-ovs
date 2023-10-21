# OVS [![Build Status](https://api.travis-ci.com/xgfone/go-ovs.svg?branch=master)](https://travis-ci.com/github/xgfone/go-ovs) [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-ovs)](https://pkg.go.dev/github.com/xgfone/go-ovs) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-ovs/master/LICENSE)

A simple OVS flow executor supporting `Go1.18+`.


## Example
```go
package main

import (
	"context"

	"github.com/xgfone/go-exec"
	"github.com/xgfone/go-ovs"
)

func main() {
	initBridge("br-ovs")
	initFlows("br-ovs")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func initBridge(bridge string) {
	must(exec.Execute(context.Background(),
		"ovs-vsctl", "--may-exist", "add-br", bridge,
		"--", "set-fail-mode", bridge, "secure",
	))

	must(exec.Execute(context.Background(),
		"ip", "link", "set", bridge, "up",
	))
}

func initFlows(bridge string) {
	// Table 0
	ovs.MustAddFlow(bridge, "table=0,in_port=1,actions=goto_table:1")
	ovs.MustAddFlow(bridge, "table=0,in_port=2,actions=goto_table:2")
	ovs.MustAddFlow(bridge, "table=0,in_port=3,actions=goto_table:3")
	ovs.MustAddFlow(bridge, "table=0,in_port=4,actions=goto_table:4")
	ovs.MustAddFlow(bridge, "table=0,priority=0,actions=drop")

	// TODO ...
}
```
