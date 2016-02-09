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
	HostZoneName string
	HostZoneID   string
}

//Not every ResourceRecordSet memember is supported yet
var hostedzonenameFlag string
var hostnameFlag string
var recordtypeFlag string
var ipaddressFlag string
var ttlFlag int64

//init all flag values
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
	record.Record = &route53.ResourceRecordSet{}
	record.Record.ResourceRecords = make([]*route53.ResourceRecord, 1)

	if arglen <= 1 {
		fmt.Println("No option set try --help for more information")
		os.Exit(-1)
	}

	if (hostedzonenameFlag != "") && (hostnameFlag != "") && (recordtypeFlag != "") {
		//debug
		fmt.Println("hosted name: ", hostedzonenameFlag, " recordname: ", hostnameFlag, " recordtype: ", recordtypeFlag)
		// TODO: Validate FLAGS!!!
		//set the vars from flags
		record.HostZoneName = hostedzonenameFlag
		record.Record.Name = aws.String(hostnameFlag)
		record.Record.TTL = aws.Int64(ttlFlag)
		record.Record.Type = aws.String(recordtypeFlag)
		record.Record.ResourceRecords[0] = &route53.ResourceRecord{Value: aws.String(ipaddressFlag)}
		//get the hosted zone id from the name
		record.getZoneID()
		//upsert record
		record.recordUpsert()
	} else {
		fmt.Println("Option missing or not set, try --help for more information")
		os.Exit(-1)
	}

}

func (record *Recordset) getZoneID() {
	svc := route53.New(session.New())
	//intialise search params
	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(record.HostZoneName),
		MaxItems: aws.String("1"), //limit records to 1
	}

	resp, err := svc.ListHostedZonesByName(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//if result does not match error out
	if record.HostZoneName+"." == *resp.HostedZones[0].Name {
		record.HostZoneID = *resp.HostedZones[0].Id
	} else {
		fmt.Println("Error: Could not find Hostzone: ", record.HostZoneName)
		os.Exit(-1)
		return
	}

}

//update zone record
func (record *Recordset) recordUpsert() {
	svc := route53.New(session.New())

	//initialise upsert params
	parameters := &route53.ChangeResourceRecordSetsInput{}
	parameters.ChangeBatch = &route53.ChangeBatch{}
	parameters.HostedZoneId = aws.String(record.HostZoneID)
	parameters.ChangeBatch.Changes = make([]*route53.Change, 1)
	parameters.ChangeBatch.Changes[0] = &route53.Change{
		Action: aws.String("UPSERT"),
		ResourceRecordSet: &route53.ResourceRecordSet{
			Type: aws.String(*record.Record.Type),
			TTL:  aws.Int64(*record.Record.TTL),
			Name: aws.String(*record.Record.Name),
			ResourceRecords: []*route53.ResourceRecord{&route53.ResourceRecord{
				Value: record.Record.ResourceRecords[0].Value,
			}},
		},
	}

	//fmt.Println(*record.Parameters.ChangeBatch.Changes[0].ResourceRecordSet)
	resp, err := svc.ChangeResourceRecordSets(parameters)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println("Error: Could not upsert record - ", err.Error())
		os.Exit(-1)
		return
	}
	// Pretty-print the response data.
	fmt.Println(resp)
}
