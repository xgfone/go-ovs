# OVS

A simple OVS flow executor.

## Install
```shell
$ go get -u github.com/xgfone/go-ovs
```

**Notice:** In order to log the every executed command, you maybe add the log hook for default command executor, `github.com/xgfone/go-tools/v7/execution#DefaultCmd`, such as `execution.DefaultCmd.AppendResultHooks(logHook)`. Or, you can call the convenient function `github.com/xgfone/goapp/exec#SetDefaultCmdLogHook`, such as `exec.SetDefaultCmdLogHook()`, which will add a default log hook.


## Example
```go
package main

import (
	"github.com/xgfone/go-ovs"
	"github.com/xgfone/goapp/exec"
)

func main() {
	exec.SetDefaultCmdLogHook() // Log the every executed command.

	initBridge("br-ovs")
	initFlows("br-ovs")
}

func initBridge(bridge string) {
	exec.MustExecutes("ovs-vsctl", "--may-exist", "add-br", bridge,
		"--", "set-fail-mode", bridge, "secure")
	exec.MustExecutes("ip", "link", "set", bridge, "up")
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
