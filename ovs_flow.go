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
	"context"
	"fmt"
	"log"
	"net/netip"
	"strconv"
	"strings"

	"github.com/xgfone/go-atexit"
	"github.com/xgfone/go-exec"
)

// StringToInt parses the decimal or hexadecimal string  to the integer,
// which will panic if failing.
func StringToInt(s string) int {
	v, err := strconv.ParseInt(strings.TrimSpace(s), 0, 64)
	if err != nil {
		panic(err)
	}
	return int(v)
}

// IntToString converts the integer to string as the decimal.
func IntToString(i int) string { return fmt.Sprintf("%d", i) }

// IntToHexString converts the integer to string as the hexdecimal
// with the prefix "0x".
func IntToHexString(i int) string { return fmt.Sprintf("0x%x", i) }

// GetAllFlows returns the list of all the flows of the bridge.
func GetAllFlows(bridge string, isName, isStats bool) (flows []string, err error) {
	var out string
	if isName {
		if isStats {
			out, err = exec.Output(context.Background(), OfctlCmd, "--names", "--stats", "dump-flows", bridge)
		} else {
			out, err = exec.Output(context.Background(), OfctlCmd, "--names", "--no-stats", "dump-flows", bridge)
		}
	} else {
		if isStats {
			out, err = exec.Output(context.Background(), OfctlCmd, "--no-names", "--stats", "dump-flows", bridge)
		} else {
			out, err = exec.Output(context.Background(), OfctlCmd, "--no-names", "--no-stats", "dump-flows", bridge)
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
		err = exec.Execute(context.Background(), OfctlCmd, "add-flow", bridge, flow)
		if err != nil {
			return
		}
	}
	return
}

// DelFlows deletes the flows.
func DelFlows(bridge string, matches ...string) (err error) {
	for _, match := range matches {
		err = exec.Execute(context.Background(), OfctlCmd, "del-flows", bridge, match)
		if err != nil {
			return
		}
	}
	return
}

// DelFlowsStrict deletes the flows with the option --strict.
func DelFlowsStrict(bridge string, priority int, matches ...string) (err error) {
	for _, match := range matches {
		match = fmt.Sprintf("priority=%d,%s", priority, match)
		err = exec.Execute(context.Background(), OfctlCmd, "--strict", "del-flows", bridge, match)
		if err != nil {
			return
		}
	}
	return
}

// MustAddFlow is the same as AddFlows, but the program exits if there is an error.
func MustAddFlow(bridge, flow string) {
	if err := AddFlows(bridge, flow); err != nil {
		log.Printf("fail to add flow: bridge=%s, flow=%s, err=%v", bridge, flow, err)
		atexit.Exit(1)
	}
}

// MustDelFlow is the same as DelFlows, but the program exits if there is an error.
func MustDelFlow(bridge, match string) {
	if err := DelFlows(bridge, match); err != nil {
		log.Printf("fail to delete flows: bridge=%s, match=%s, err=%v", bridge, match, err)
		atexit.Exit(1)
	}
}

// MustDelFlowStrict is the same as DelFlowsStrict, but the program exits if there is an error.
func MustDelFlowStrict(bridge string, priority int, match string) {
	if err := DelFlowsStrict(bridge, priority, match); err != nil {
		log.Printf("fail to delete flows: bridge=%s, priority=%d, match=%s, err=%v", bridge, priority, match, err)
		atexit.Exit(1)
	}
}

//////////////////////////////////////////////////////////////////////////////

var arpPacket = "ffffffffffff%s%s08060001080006040001%s%sffffffffffff%s"

// SendARPRequest sends the ARP request by the ovs bridge.
//
// vlanID may be 0, which won't add the VLAN header into the ARP request packet.
func SendARPRequest(bridge, output, inPort, srcMac, srcIP, dstIP string,
	vlanID ...uint16) (err error) {

	srcmac := strings.Replace(normalizeMac(srcMac), ":", "", -1)
	if srcmac == "" {
		return fmt.Errorf("invalid src mac '%s'", srcMac)
	}

	srcip, err := netip.ParseAddr(srcIP)
	if err != nil {
		return
	}
	srcIP = fmt.Sprintf("%0x", srcip.As4())

	dstip, err := netip.ParseAddr(dstIP)
	if err != nil {
		return
	}
	dstIP = fmt.Sprintf("%0x", dstip.As4())

	var vlan string
	if len(vlanID) != 0 && vlanID[0] != 0 {
		vlan = fmt.Sprintf("8100%04x", vlanID[0])
	}

	pkt := fmt.Sprintf(arpPacket, srcmac, vlan, srcmac, srcIP, dstIP)
	return exec.Execute(context.Background(), OfctlCmd, "packet-out", bridge, inPort, output, pkt)
}

func normalizeMac(mac string) string {
	macs := strings.Split(mac, ":")
	if len(macs) != 6 {
		return ""
	}

	for i, m := range macs {
		v, err := strconv.ParseUint(m, 16, 8)
		if err != nil {
			return ""
		}
		macs[i] = fmt.Sprintf("%02x", v)
	}

	return strings.Join(macs, ":")
}
