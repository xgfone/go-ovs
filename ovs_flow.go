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

// Package ovs supplies some convenient functions to execute the ovs command.
package ovs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xgfone/go-tools/v7/net2"
	"github.com/xgfone/goapp/exec"
	"github.com/xgfone/goapp/log"
	"github.com/xgfone/netaddr"
)

// PortStringToInt parses the decimal or hexadecimal string integer,
// which will panic if failing to parse the port as int.
func PortStringToInt(port string) int {
	v, err := strconv.ParseInt(port, 0, 64)
	if err != nil {
		panic(err)
	}
	return int(v)
}

// PortIntToString converts the port from integer to string as the decimal.
func PortIntToString(port int) string { return fmt.Sprintf("%d", port) }

// PortIntToHexString converts the port from integer to string as the hexdecimal
// with the prefix "0x".
func PortIntToHexString(port int) string { return fmt.Sprintf("0x%x", port) }

// GetAllFlows returns the list of all the flows of the bridge.
func GetAllFlows(bridge string, isName, isStats bool) (flows []string, err error) {
	var out string
	if isName {
		if isStats {
			out, err = exec.Outputs("ovs-ofctl", "--names", "--stats", "dump-flows", bridge)
		} else {
			out, err = exec.Outputs("ovs-ofctl", "--names", "--no-stats", "dump-flows", bridge)
		}
	} else {
		if isStats {
			out, err = exec.Outputs("ovs-ofctl", "--no-names", "--stats", "dump-flows", bridge)
		} else {
			out, err = exec.Outputs("ovs-ofctl", "--no-names", "--no-stats", "dump-flows", bridge)
		}
	}

	if err == nil {
		flows = strings.Split(out, "\n")
	}

	return
}

// AddFlows adds the flows.
func AddFlows(bridge string, flows ...string) (err error) {
	for _, flow := range flows {
		if err = exec.Executes("ovs-ofctl", "add-flow", bridge, flow); err != nil {
			return
		}
	}
	return
}

// DelFlows deletes the flows.
func DelFlows(bridge string, matches ...string) (err error) {
	for _, match := range matches {
		if err = exec.Executes("ovs-ofctl", "del-flows", bridge, match); err != nil {
			return
		}
	}
	return
}

// DelFlowsStrict deletes the flows with the option --strict.
func DelFlowsStrict(bridge string, priority int, matches ...string) (err error) {
	for _, match := range matches {
		match = fmt.Sprintf("priority=%d,%s", priority, match)
		err = exec.Executes("ovs-ofctl", "--strict", "del-flows", bridge, match)
		if err != nil {
			return
		}
	}
	return
}

// MustAddFlow is the same as AddFlows, but the program exits if there is an error.
func MustAddFlow(bridge, flow string) {
	if err := AddFlows(bridge, flow); err != nil {
		log.Fatalf("fail to add flow: %s", err)
	}
}

// MustDelFlow is the same as DelFlows, but the program exits if there is an error.
func MustDelFlow(bridge, match string) {
	if err := DelFlows(bridge, match); err != nil {
		log.Fatalf("fail to delete flows: %s", err)
	}
}

// MustDelFlowStrict is the same as DelFlowsStrict, but the program exits if there is an error.
func MustDelFlowStrict(bridge string, priority int, match string) {
	if err := DelFlowsStrict(bridge, priority, match); err != nil {
		log.Fatalf("fail to delete flows: %s", err)
	}
}

//////////////////////////////////////////////////////////////////////////////

var arpPacket = "ffffffffffff%s%s08060001080006040001%s%sffffffffffff%s"

// SendARPRequest sends the ARP request by the ovs bridge.
//
// vlanID may be 0, which won't add the VLAN header into the ARP request packet.
func SendARPRequest(bridge, output, inPort, srcMac, srcIP, dstIP string,
	vlanID ...uint16) (err error) {

	srcmac := strings.Replace(net2.NormalizeMacFu(srcMac), ":", "", -1)
	if srcmac == "" {
		return fmt.Errorf("invalid src mac '%s'", srcMac)
	}

	srcip, err := netaddr.NewIPAddress(srcIP)
	if err != nil {
		return
	}
	srcIP = srcip.Hex()

	dstip, err := netaddr.NewIPAddress(dstIP)
	if err != nil {
		return
	}
	dstIP = dstip.Hex()

	var vlan string
	if len(vlanID) != 0 && vlanID[0] != 0 {
		vlan = fmt.Sprintf("8100%04x", vlanID[0])
	}

	pkt := fmt.Sprintf(arpPacket, srcmac, vlan, srcmac, srcIP, dstIP)
	exec.Execute("ovs-ofctl", "packet-out", bridge, inPort, output, pkt)
	return
}
