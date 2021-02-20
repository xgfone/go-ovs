package ovs

// Some linux commands.
var (
	IPCmd    = "ip"
	OfctlCmd = "ovs-ofctl"
	VsctlCmd = "ovs-vsctl"
)

// L2 Data-Link Protocol Number
const (
	ARP  = 0x0806
	IPv4 = 0x0800
	IPv6 = 0x86dd
)

// L3 IP Protocol Number
const (
	ICMP = 1
	TCP  = 6
	UDP  = 17
	GRE  = 47
)

// OVS Actions
const (
	DROP   = "drop"
	LOCAL  = "local"
	FLOOD  = "flood"
	NORMAL = "normal"
)

// MAC address
const (
	BroadcastMac = "ff:ff:ff:ff:ff:ff"
)
