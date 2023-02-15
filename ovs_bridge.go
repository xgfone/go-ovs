// Copyright 2020 xgfone
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

package ovs

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xgfone/go-atexit"
	"github.com/xgfone/go-exec"
)

// ListAllOFPorts returns all the port names with its number on the bridge.
func ListAllOFPorts(bridge string) (map[string]int, error) {
	out, err := exec.Output(context.Background(), OfctlCmd, "show", bridge)
	if err != nil {
		return nil, err
	}

	if out == "" {
		return map[string]int{}, nil
	}

	lines := strings.Split(out, "\n")
	ports := make(map[string]int, len(lines))
	for _, line := range lines {
		if line = strings.TrimSpace(line); line == "" || !strings.Contains(line, " addr:") {
			continue
		} else if strings.HasPrefix(line, "LOCAL") {
			continue
		}

		if items := strings.Split(strings.SplitN(line, "):", 2)[0], "("); len(items) == 2 {
			v, err := strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				return nil, err
			}
			ports[items[1]] = int(v)
		}
	}

	return ports, nil
}

// SetInterfaceUp sets up the interface.
func SetInterfaceUp(iface string) (err error) {
	return exec.Execute(context.Background(), IPCmd, "link", "set", iface, "up")
}

// CreateBridge creates a new bridge named name if not exist.
//
// If secureFailMode is true, set the fail mode of the bridge to "secure".
func CreateBridge(name string, secureFailMode ...bool) (err error) {
	if len(secureFailMode) > 0 && secureFailMode[0] {
		err = exec.Execute(context.Background(), VsctlCmd,
			"--may-exist", "add-br", name,
			"--", "set-fail-mode", name, "secure")
	} else {
		err = exec.Execute(context.Background(), VsctlCmd, "--may-exist", "add-br", name)
	}

	if err == nil {
		err = exec.Execute(context.Background(), IPCmd, "link", "set", name, "up")
	}

	return
}

// DeleteBridge deletes the bridge named name.
func DeleteBridge(name string) (err error) {
	return exec.Execute(context.Background(), VsctlCmd, "--if-exists", "del-br", name)
}

// AddPort adds the interface to the bridge.
func AddPort(bridge, iface string, ofport int) (err error) {
	if ofport == 0 {
		err = exec.Execute(context.Background(), VsctlCmd,
			"--may-exist", "add-port", bridge, iface)
	} else {
		err = exec.Execute(context.Background(), VsctlCmd,
			"--may-exist", "add-port", bridge, iface,
			"--", "set", "interface", iface, fmt.Sprintf("ofport_request=%d", ofport))
	}

	return
}

// DelPort deletes the port from the bridge.
func DelPort(bridge, port string) (err error) {
	return exec.Execute(context.Background(), VsctlCmd, "--if-exists", "del-port", bridge, port)
}

// AddPatchPort adds a patch port for the bridge, the peer patch of which
// is peerPatch.
func AddPatchPort(bridge, patch, peerPatch string, ofport int) (err error) {
	args := []string{
		"--may-exist", "add-port", bridge, patch,
		"--", "set", "interface", patch, "type=patch",
		"--", "set", "interface", patch, "options:peer=" + peerPatch,
	}
	if ofport > 0 {
		args = append(args, "--", "set", "interface", patch, fmt.Sprintf("ofport_request=%d", ofport))
	}

	return exec.Execute(context.Background(), VsctlCmd, args...)
}

// AddVxLANPort add an VxLAN port into the bridge.
func AddVxLANPort(bridge, port, localIP, remoteIP string, ofport int) (err error) {
	args := []string{
		"-may-exist", "add-port", bridge, port,
		"--", "set", "interface", port, "type=vxlan",
		fmt.Sprintf("options:local_ip=%s", localIP),
		fmt.Sprintf("options:remote_ip=%s", remoteIP),
		"options:in_key=flow", "options:out_key=flow",
		"options:df_default=true",
	}
	if ofport > 0 {
		args = append(args, "--", "set", "interface", port, fmt.Sprintf("ofport_request=%d", ofport))
	}

	return exec.Execute(context.Background(), VsctlCmd, args...)
}

// MustSetInterfaceUp is the same as SetInterfaceUp, but exit the program if failing.
func MustSetInterfaceUp(iface string) {
	if err := SetInterfaceUp(iface); err != nil {
		log.Printf("fail to set up the interface: interface=%s, err=%v", iface, err)
		atexit.Exit(1)
	}
}

// MustCreateBridge is the same as CreateBridge, but exit the program if failing.
func MustCreateBridge(name string, secureFailMode ...bool) {
	if err := CreateBridge(name, secureFailMode...); err != nil {
		log.Printf("failed to create bridge: bridge=%s, err=%v", name, err)
		atexit.Exit(1)
	}
}

// MustDeleteBridge is the same as DeleteBridge, but exit the program if failing.
func MustDeleteBridge(name string) {
	if err := DeleteBridge(name); err != nil {
		log.Printf("failed to delete bridge: bridge=%s, err=%v", name, err)
		atexit.Exit(1)
	}
}

// MustAddPort is the same as AddPort, but exit the program if failing.
func MustAddPort(bridge, iface string, ofport int) {
	if err := AddPort(bridge, iface, ofport); err != nil {
		log.Printf("fail to add the port to the bridge: bridge=%s, interface=%s, ofport=%d, err=%v", bridge, iface, ofport, err)
		atexit.Exit(1)
	}
}

// MustDelPort is the same as DelPort, but exit the program if failing.
func MustDelPort(bridge, iface string) {
	if err := DelPort(bridge, iface); err != nil {
		log.Printf("fail to delete the port from the bridge: bridge=%s, interface=%s, err=%v", bridge, iface, err)
		atexit.Exit(1)
	}
}

// MustAddPatchPort is the same as AddPatchPort, but exit the program if failing.
func MustAddPatchPort(bridge, patch, peerPatch string, ofport int) {
	if err := AddPatchPort(bridge, patch, peerPatch, ofport); err != nil {
		log.Printf("fail to add the patch port to the bridge: bridge=%s, patch=%s, peer=%s, ofport=%d, err=%v",
			bridge, patch, peerPatch, ofport, err)
		atexit.Exit(1)
	}
}

// MustAddVxLANPort is the same as AddVxLANPort, but exit the program if failing.
func MustAddVxLANPort(bridge, port, localIP, remoteIP string, ofport int) {
	if err := AddVxLANPort(bridge, port, localIP, remoteIP, ofport); err != nil {
		log.Printf("fail to add the vxlan port to the bridge: bridge=%s, port=%s, localip=%s, remoteip=%s, ofport=%d, err=%v",
			bridge, port, localIP, remoteIP, ofport, err)
		atexit.Exit(1)
	}
}
