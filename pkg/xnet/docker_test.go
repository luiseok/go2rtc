package xnet

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDockerNets(t *testing.T) {
	lan := net.ParseIP("172.16.200.10")
	hassio := net.ParseIP("172.30.32.1")
	docker0 := net.ParseIP("172.30.232.1")
	bridge := net.ParseIP("172.18.0.1")
	home := net.ParseIP("192.168.0.10")

	contains := func(nets []*net.IPNet, ip net.IP) bool {
		for _, ipn := range nets {
			if ipn.Contains(ip) {
				return true
			}
		}
		return false
	}

	// unset: upstream default 172.16.0.0/12
	nets := ParseDockerNets("", false)
	require.True(t, contains(nets, lan))
	require.True(t, contains(nets, hassio))
	require.True(t, contains(nets, bridge))
	require.False(t, contains(nets, home))

	// HAOS bridges only: LAN in 172.16.0.0/12 is kept
	nets = ParseDockerNets(" 172.30.32.0/23, 172.30.232.0/23 ", true)
	require.Len(t, nets, 2)
	require.False(t, contains(nets, lan))
	require.True(t, contains(nets, hassio))
	require.True(t, contains(nets, docker0))
	require.False(t, contains(nets, bridge))

	// empty: filter disabled
	nets = ParseDockerNets("", true)
	require.Empty(t, nets)
	require.False(t, contains(nets, hassio))

	// invalid entries are ignored
	nets = ParseDockerNets("bad,172.30.32.0/23", true)
	require.Len(t, nets, 1)

	// only invalid entries: fall back to default
	nets = ParseDockerNets("172.30.32.0/33", true)
	require.True(t, contains(nets, lan))
}

func TestIsDocker(t *testing.T) {
	saved := DockerNets
	defer func() { DockerNets = saved }()

	DockerNets = ParseDockerNets("172.30.32.0/23", true)
	require.True(t, IsDocker(net.ParseIP("172.30.33.5")))
	require.False(t, IsDocker(net.ParseIP("172.16.200.10")))

	DockerNets = ParseDockerNets("", true)
	require.False(t, IsDocker(net.ParseIP("172.30.33.5")))
}
