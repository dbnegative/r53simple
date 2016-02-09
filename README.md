# r53simple

Simple utility to Upsert a Route53 record from an AWS ASG instance. Currently only A records supported.

## Usage
```
Usage of ./r53simple:
  -hostedzone string 
    	The hosted zone name that will be updated (Required)
  -ip string
    	The record's IP if record type is A
  -recordname string
    	The name of the FQDN you want to perform the action on (Required)
  -ttl int
    	The cache time to live for the current resource record set - default 300 seconds (default 300)
  -type string
    	Valid values for basic resource record sets:
 A | AAAA | CNAME | MX | NS | PTR | SOA | SPF | SRV | TXT (Required)
```
## Build

Requires the AWS Go SDK

Linux:
```
GOOS=linux GOARCH=amd64 go build r53simple.go
```



