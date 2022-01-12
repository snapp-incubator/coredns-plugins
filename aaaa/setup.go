package aaaa

import (
        "github.com/coredns/caddy"
        "github.com/coredns/coredns/core/dnsserver"
        "github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("aaaa", setup) }

func setup(c *caddy.Controller) error {
        a := AAAA{}

        dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
                a.Next = next
                return a
        })

        return nil
}
