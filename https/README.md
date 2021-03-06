# https

(from https://github.com/v-byte-cpu/coredns/tree/https_plugin/plugin/https)

## Name

*https* - facilitates proxying DNS messages to upstream resolvers using DoH.

## Description

The *https* plugin performs DNS-over-HTTPS proxying. See [RFC 8484](https://tools.ietf.org/html/rfc8484).

This plugin can only be used once per Server Block.

## Syntax

In its most basic form:

~~~
https FROM TO...
~~~

* **FROM** is the base domain to match for the request to be proxied.
* **TO...** are the destination endpoints to proxy to. The number of upstreams is
  limited to 15.

Multiple upstreams are randomized (see `policy`) on first use. When a proxy returns an error
the next upstream in the list is tried.

Extra knobs are available with an expanded syntax:

~~~
https FROM TO... {
    except IGNORED_NAMES...
    tls CERT KEY CA
    tls_servername NAME
    policy random|round_robin|sequential
}
~~~

* **FROM** and **TO...** as above.
* **IGNORED_NAMES** in `except` is a space-separated list of domains to exclude from proxying.
  Requests that match none of these names will be passed through.
* `tls` **CERT** **KEY** **CA** define the TLS properties for TLS connection. From 0 to 3 arguments can be
  provided with the meaning as described below

  * `tls` - no client authentication is used, and the system CAs are used to verify the server certificate (by default)
  * `tls` **CA** - no client authentication is used, and the file CA is used to verify the server certificate
  * `tls` **CERT** **KEY** - client authentication is used with the specified cert/key pair.
    The server certificate is verified with the system CAs
  * `tls` **CERT** **KEY**  **CA** - client authentication is used with the specified cert/key pair.
    The server certificate is verified using the specified CA file

* `policy` specifies the policy to use for selecting upstream servers. The default is `random`.


## Metrics

If monitoring is enabled (via the *prometheus* plugin) then the following metric are exported:

* `coredns_https_request_duration_seconds{to}` - duration per upstream interaction.
* `coredns_https_requests_total{to}` - query count per upstream.
* `coredns_https_responses_total{to, rcode}` - count of RCODEs per upstream.
  and we are randomly (this always uses the `random` policy) spraying to an upstream.

## Examples

Proxy all requests within `example.org.` to a DoH nameserver:

~~~ corefile
example.org {
    https . cloudflare-dns.com/dns-query
}
~~~

Forward everything except requests to `example.org`

~~~ corefile
. {
    https . dns.quad9.net/dns-query {
        except example.org
    }
}
~~~

Load balance all requests between multiple upstreams

~~~ corefile
. {
    https . dns.quad9.net/dns-query cloudflare-dns.com:443/dns-query dns.google/dns-query
}
~~~

Internal DoH server:

~~~ corefile
. {
    https . 10.0.0.10:853/dns-query {
      tls ca.crt
      tls_servername internal.domain
    }
}
~~~
