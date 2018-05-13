package cmd

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var ec2Name = "ro-demo"

// CreateEC2Instance Creates a public ec2 instance with the running web app
func CreateEC2Instance(ec2Svc *ec2.EC2) *ec2.Instance {
	ec2Response, err := ec2Svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0b91bd72"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      aws.String("default"),
		UserData:     getUserData(),
	})
	if err != nil {
		logrus.Fatalf("Couldn't provision instance: %v", err)
	}
	_, err = ec2Svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{ec2Response.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(ec2Name),
			},
		},
	})
	if err != nil {
		logrus.Errorf("Unable to tag instance: %v", err)
	}
	return ec2Response.Instances[0]
}

func getUserData() *string {
	data := `#!/bin/bash
apt update
apt install -y docker.io
docker run -d -p 80:3000 awithrow/rails:latest`
	b64 := base64.StdEncoding.EncodeToString([]byte(data))
	return &b64
}

// WaitForInstanceToBeLive will wait for a given instance ID to be ready and serving traffic on port 80
func WaitForInstanceToBeLive(ins *ec2.Instance, ec2Scv *ec2.EC2) error {
	logrus.Infof("Waiting for %v's status to be Ready", *ins.InstanceId)
	ec2Scv.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{ins.InstanceId},
	})
	logrus.Infof("%v is ready. Checking if the web application is serving traffic", *ins.InstanceId)
	// Refresh the instance data
	result, err := ec2Scv.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{ins.InstanceId},
	})
	if err != nil {
		return err
	}
	info := result.Reservations[0].Instances[0]
	ready := false
	publicIP := *info.PublicIpAddress
	url := "http://" + publicIP
	for x := 1; x <= 15; x++ {
		response, err := http.Get(url)
		if err == nil {
			defer response.Body.Close()
			if response.StatusCode == 200 {
				ready = true
				break
			}
		}
		logrus.Infof("Still waiting for %v to become available", url)
		time.Sleep(10 * time.Second)
	}
	if !ready {
		return fmt.Errorf("Instance is not yet listening on %v", url)
	}
	logrus.Infof("%v is now ready for traffic at %v", *info.InstanceId, url)
	return nil
}
