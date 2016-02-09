package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

//Recordset contains Route53 Recordset
type Recordset struct {
	Record       *route53.ResourceRecordSet
	Parameters   *route53.ChangeResourceRecordSetsInput
	HostZoneName string
}

//Not every ResourceRecordSet memember is supported yet
var hostedzonenameFlag string
var hostnameFlag string
var recordtypeFlag string
var ipaddressFlag string

//var regionFlag string
//var geolocationFlag string
//var healthcheckidFlag string
//var failoverFlag string
//var resourceidentifierFlag string
var ttlFlag int64

//var trafficpolicyinstanceidFlag string
//var aliastargetFlag string
//var weightFlag int64

func init() {
	flag.StringVar(&hostedzonenameFlag, "hostedzone", "", "The hosted zone name that will be updated (Required)")
	flag.StringVar(&hostnameFlag, "recordname", "", "The name of the FQDN you want to perform the action on (Required)")
	flag.StringVar(&recordtypeFlag, "type", "", "Valid values for basic resource record sets: \n A | AAAA | CNAME | MX | NS | PTR | SOA | SPF | SRV | TXT (Required)")
	flag.StringVar(&ipaddressFlag, "ip", "", "The record's IP if record type is A ")
	flag.Int64Var(&ttlFlag, "ttl", 300, "The cache time to live for the current resource record set - default 300 seconds")
}

func main() {
	flag.Parse()

	arglen := len(os.Args)
	record := &Recordset{}
	if arglen <= 1 {
		fmt.Println("No option set try --help for more information")
	}

	if (hostedzonenameFlag != "") && (hostnameFlag != "") && (recordtypeFlag != "") {
		//set the hosted zone name from the flag
		record.HostZoneName = hostedzonenameFlag
		record.Record.TTL = aws.Int64(ttlFlag)
		record.Record.Type = aws.String(recordtypeFlag)
		record.Record.ResourceRecord.Value = aws.String(ipaddressFlag)
		//get the hosted zone id from the name
		record.getzoneid()

	}

}

func (record Recordset) getzoneid() {
	svc := route53.New(session.New())
	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(record.HostZoneName),
		MaxItems: aws.String("1"),
	}
	resp, err := svc.ListHostedZonesByName(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	//if result does not match error out
	if record.HostZoneName == *resp.DNSName {
		record.Parameters.HostedZoneId = resp.HostedZoneId
	} else {
		fmt.Println("Error: Could not find Hostzone: ", record.HostZoneName)
		os.Exit(-1)
		return
	}

}

//update zone record
func (record Recordset) recordupsert() {
	svc := route53.New(session.New())

	resp, err := svc.ChangeResourceRecordSets(record.Parameters)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}
	// Pretty-print the response data.
	fmt.Println(resp)
}
