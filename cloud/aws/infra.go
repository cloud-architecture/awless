package aws

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type Infra struct {
	ec2iface.EC2API
}

func NewInfra(sess *session.Session) *Infra {
	return &Infra{ec2.New(sess)}
}

func (inf *Infra) Instances() (interface{}, error) {
	return inf.DescribeInstances(&ec2.DescribeInstancesInput{})
}

func (inf *Infra) Images() (interface{}, error) {
	return inf.DescribeImages(&ec2.DescribeImagesInput{})
}

func (inf *Infra) Vpcs() (interface{}, error) {
	return inf.DescribeVpcs(&ec2.DescribeVpcsInput{})
}

func (inf *Infra) Subnets() (interface{}, error) {
	return inf.DescribeSubnets(&ec2.DescribeSubnetsInput{})
}

func (inf *Infra) Keypairs() (interface{}, error) {
	return inf.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
}

func (inf *Infra) SecurityGroups() (interface{}, error) {
	return inf.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
}

func (inf *Infra) Volumes() (interface{}, error) {
	return inf.DescribeVolumes(&ec2.DescribeVolumesInput{})
}

func (inf *Infra) Regions() (interface{}, error) {
	return inf.DescribeRegions(&ec2.DescribeRegionsInput{})
}

func (inf *Infra) InternetGateways() (interface{}, error) {
	return inf.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{})
}

func (inf *Infra) RouteTables() (interface{}, error) {
	return inf.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
}

func (inf *Infra) Vpc(id string) (interface{}, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{awssdk.String(id)},
	}

	return inf.DescribeVpcs(input)
}

func (inf *Infra) CreateInstance(ami string) (interface{}, error) {
	params := &ec2.RunInstancesInput{
		ImageId:      awssdk.String(ami),
		MaxCount:     awssdk.Int64(1),
		MinCount:     awssdk.Int64(1),
		InstanceType: awssdk.String("t2.micro"),
	}

	return inf.RunInstances(params)
}

type AwsInfra struct {
	Vpcs             []*ec2.Vpc
	Subnets          []*ec2.Subnet
	SecurityGroups   []*ec2.SecurityGroup
	Instances        []*ec2.Instance
	Keypairs         []*ec2.KeyPairInfo
	Volumes          []*ec2.Volume
	InternetGateways []*ec2.InternetGateway
	RouteTables      []*ec2.RouteTable
}

func (inf *Infra) FetchAwsInfra() (*AwsInfra, error) {
	resultc, errc := multiFetch(inf.Instances, inf.Subnets, inf.Vpcs, inf.SecurityGroups, inf.Keypairs, inf.Volumes, inf.InternetGateways, inf.RouteTables)

	awsInfra := &AwsInfra{}

	for r := range resultc {
		switch rr := r.(type) {
		case *ec2.DescribeVpcsOutput:
			awsInfra.Vpcs = append(awsInfra.Vpcs, rr.Vpcs...)
		case *ec2.DescribeSubnetsOutput:
			awsInfra.Subnets = append(awsInfra.Subnets, rr.Subnets...)
		case *ec2.DescribeSecurityGroupsOutput:
			awsInfra.SecurityGroups = append(awsInfra.SecurityGroups, rr.SecurityGroups...)
		case *ec2.DescribeKeyPairsOutput:
			awsInfra.Keypairs = append(awsInfra.Keypairs, rr.KeyPairs...)
		case *ec2.DescribeVolumesOutput:
			awsInfra.Volumes = append(awsInfra.Volumes, rr.Volumes...)
		case *ec2.DescribeInternetGatewaysOutput:
			awsInfra.InternetGateways = append(awsInfra.InternetGateways, rr.InternetGateways...)
		case *ec2.DescribeRouteTablesOutput:
			awsInfra.RouteTables = append(awsInfra.RouteTables, rr.RouteTables...)
		case *ec2.DescribeInstancesOutput:
			for _, reservation := range rr.Reservations {
				awsInfra.Instances = append(awsInfra.Instances, reservation.Instances...)
			}
		}
	}

	return awsInfra, <-errc
}
