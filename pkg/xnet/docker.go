package xnet

import (
	"net"
	"os"
	"strings"
)

// DockerNetsEnv overrides the list of Docker-like networks that go2rtc skips
// for mDNS (HomeKit) and WebRTC host candidates.
//
//   - unset: default Docker range 172.16.0.0/12 (upstream behaviour)
//   - "172.30.32.0/23,172.30.232.0/23": skip only these networks
//     (Home Assistant OS: hassio and docker0 bridges)
//   - "" (empty): disable the Docker filter completely
//
// Useful when the LAN itself uses 172.16.0.0/12 addresses (ex. 172.16.200.0/24).
const DockerNetsEnv = "GO2RTC_DOCKER_NETS"

// DockerNets is the active list of Docker-like networks.
var DockerNets = ParseDockerNets(os.LookupEnv(DockerNetsEnv))

// ParseDockerNets parses a comma separated CIDR list.
// Invalid entries are ignored. If the value is not empty but contains no valid
// entries, the default Docker range is used, so a typo can't silently disable the filter.
func ParseDockerNets(value string, ok bool) []*net.IPNet {
	if !ok {
		return []*net.IPNet{&Docker}
	}

	nets := []*net.IPNet{}
	for _, s := range strings.Split(value, ",") {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}
		if _, ipn, err := net.ParseCIDR(s); err == nil {
			nets = append(nets, ipn)
		}
	}

	if len(nets) == 0 && strings.TrimSpace(value) != "" {
		return []*net.IPNet{&Docker}
	}

	return nets
}

// IsDocker returns true if ip belongs to one of DockerNets.
func IsDocker(ip net.IP) bool {
	for _, ipn := range DockerNets {
		if ipn.Contains(ip) {
			return true
		}
	}
	return false
}
