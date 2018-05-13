package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

var vpcName string = "ro-vpc"
var vpcInput = &ec2.DescribeVpcsInput{
	Filters: []*ec2.Filter{
		{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String(vpcName)},
		},
	},
}

// VpcExists checks that the VPC we expect exists
func VpcExists(ec2Svc *ec2.EC2) bool {
	logrus.Info("Checking for VPC")
	vpcResponse, err := ec2Svc.DescribeVpcs(vpcInput)
	if err != nil {
		logrus.Fatalf("Could not fetch VPCs: %s", err)
	}
	if len(vpcResponse.Vpcs) == 0 {
		logrus.Info("No VPCs have been created. Creating VPC")
		return false
	}
	return true
}

// CreateVpc will create a simple vpc for the server to live in
func CreateVpc(ec2Svc *ec2.EC2) {
	logrus.Info("Creating VPC")
	vpcResponse, err := ec2Svc.CreateDefaultVpc(&ec2.CreateDefaultVpcInput{})
	if err != nil {
		logrus.Fatalf("Couldn't create a VPC: %s", err)
	}
	logrus.Info("Waiting for VPC to be created")
	ec2Svc.WaitUntilVpcExists(vpcInput)
	logrus.Infof("VPC: %s has been created", *vpcResponse.Vpc.VpcId)
	_, err = ec2Svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{vpcResponse.Vpc.VpcId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(vpcName),
			},
		},
	})
	if err != nil {
		logrus.Fatalf("Unable to tag VPC: %s", err)
	}
}
