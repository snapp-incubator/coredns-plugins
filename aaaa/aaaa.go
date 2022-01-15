package aaaa

import (
	"context"

	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

// AAAA is a plugin that returns a NXDOMAIN reply to AAAA queries.
type AAAA struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (a AAAA) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if r.Question[0].Qtype != dns.TypeAAAA {
		return plugin.NextOrFailure(a.Name(), a.Next, ctx, w, r)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	// set NXDOMAIN for AAAA record type
	m.Rcode = 3
	w.WriteMsg(m)
	return 0, nil
}

// Name implements the Handler interface.
func (a AAAA) Name() string { return "AAAA" }
