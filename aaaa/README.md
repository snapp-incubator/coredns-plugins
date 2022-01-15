# aaaa

## Name

*aaaa* - gives NXDOMAIN response to AAAA queries.

## Description

*aaaa* basically blocks AAAA queries by setting RCODE to NXDOAMIN.
## Syntax

~~~ txt
aaaa
~~~

## Examples

~~~ corefile
example.org {
    whoami
    aaaa
}
~~~

A `dig +nocmd aaaa example.org +answer` now returns:

~~~ txt
->>HEADER<<- opcode: QUERY, status: NXDOMAIN, id: 35098
~~~
