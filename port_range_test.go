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

import "fmt"

func ExamplePortRuleMasking() {
	ports := PortRuleMasking(1000, 1999)
	fmt.Println(len(ports))
	fmt.Println(ports[0])
	fmt.Println(ports[1])
	fmt.Println(ports[2])
	fmt.Println(ports[3])
	fmt.Println(ports[4])
	fmt.Println(ports[5])
	fmt.Println(ports[6])

	// Unordered output:
	// 7
	// 0x03e8/0xfff8
	// 0x03f0/0xfff0
	// 0x0400/0xfe00
	// 0x0600/0xff00
	// 0x0700/0xff80
	// 0x0780/0xffc0
	// 0x07c0/0xfff0
}
