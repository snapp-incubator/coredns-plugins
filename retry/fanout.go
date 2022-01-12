package retry

import (
	"fmt"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin/pkg/parse"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/pkg/errors"
	"github.com/snapp-incubator/fanout"
)

func initFanout(c *caddy.Controller) (*fanout.Fanout, error) {
	f := fanout.New()
	from := "."
	f.AddFrom(from)

	if !c.Args(&from) {
		return f, c.ArgErr()
	}

	to := c.RemainingArgs()
	if len(to) == 0 {
		return f, c.ArgErr()
	}

	toHosts, err := parse.HostPortOrFile(to...)
	if err != nil {
		return f, err
	}
	//set default value
	net := "udp"
	if c.NextBlock() {
		if strings.ToLower(c.Val()) == "network" {
			net, err = parseProtocol(c)
			if err != nil {
				return f, err
			}

		} else {
			return f, fmt.Errorf("additional parameters not allowed")

		}

	}

	for c.NextBlock() {
		return f, fmt.Errorf("additional parameters not allowed")
	}
	for _, host := range toHosts {
		trans, h := parse.Transport(host)
		if trans != transport.DNS {
			return f, fmt.Errorf("only dns transport allowed")
		}
		client := fanout.NewClient(h, net)
		f.AddClient(client)

	}
	return f, nil
}

func parseProtocol(c *caddy.Controller) (string, error) {
	if !c.NextArg() {
		return "", c.ArgErr()
	}
	net := strings.ToLower(c.Val())
	if net != "tcp" && net != "udp" && net != "tcptlc" {
		return "", errors.New("unknown network protocol")
	}
	return net, nil
}
