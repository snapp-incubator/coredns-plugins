# retry

## Name

Plugin *Retry* is able to parallel proxying the query to upstream resolvers, depending the error result provided by the initial resolver

## Description

The *retry* plugin will be used if the plugin chain returns specific error messages. The *retry* plugin will be replicated query in parallel to each listed IP (i.e. the DNS servers). The first non-negative response from any of the queried DNS Servers will be forwarded as a response to the application's DNS request.


## Syntax

```
{
    retry [original] RCODE_1[,RCODE_2,RCODE_3...] . DNS_RESOLVERS
}
```

* **original** is optional flag. If it is set then retry uses original request instead of potentially changed by other plugins
* **RCODE** is the string representation of the error response code. The complete list of valid rcode strings are defined as `RcodeToString` in <https://github.com/miekg/dns/blob/master/msg.go>, examples of which are `SERVFAIL`, `NXDOMAIN` and `REFUSED`. At least one rcode is required, but multiple rcodes may be specified, delimited by commas.
* **DNS_RESOLVERS** accepts dns resolvers list.
* **network** is a specific network protocol. Could be `tcp`, `udp`, `tcp-tls`.


## Building CoreDNS with Retry

When building CoreDNS with this plugin, _retry_ should be positioned **before** _forward_ in `/plugin.cfg`.

## Examples

### Retry to local DNS server

The following specifies that all requests are forwarded to 8.8.8.8. If the response is `NXDOMAIN`, *retry* will proxy the request to DNS_RESOLVERS, The first positive response from a proxy will be provided as the result.

```
. {
	forward . 8.8.8.8
	retry NXDOMAIN . 192.168.1.1:54 192.168.1.1:55
	log
}

```
### Retry with original request used

The following specify that `original` query will be forwarded to 192.168.1.1:53 if 8.8.8.8 response is `NXDOMAIN`. `original` means no changes from next plugins on request. With no `original` flag retry will forward request with EDNS0 option (set by rewrite).

```
. {
	forward . 8.8.8.8
	rewrite edns0 local set 0xffee 0x61626364
	retry original NXDOMAIN . 192.168.1.1:54 192.168.1.1:55
	log
}

```

### Multiple retries

Multiple retries can be specified, as long as they serve unique error responses.

```
. {
    forward . 8.8.8.8
    retry NXDOMAIN . 192.168.1.1:53
    retry original SERVFAIL,REFUSED . 192.168.1.1:54 192.168.1.1:55
    log
}

```

Sends parallel requests between three resolvers, one of which has a IPv6 address via TCP. The first response from proxy will be provided as the result.

```
. {
    forward . 8.8.8.8
    fanout . 192.168.1.1:54 192.168.1.1:55 {
        network TCP
    }
}
```