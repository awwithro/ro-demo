package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Terminates ec2 instances named 'ro-demo'",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Terminating previously created instances")
		awsSession := session.New(&aws.Config{Region: aws.String(Region)})
		ec2Svc := ec2.New(awsSession)
		reservations, err := ec2Svc.DescribeInstances(&ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: []*string{aws.String(ec2Name)},
				},
				{
					Name:   aws.String("instance-state-name"),
					Values: []*string{aws.String("running")},
				},
			},
		})
		if err != nil {
			logrus.Fatalf("Unable to get instances: %v", err)
		}
		// Get the instanceIDs from the reservations returned
		var instanceIds []*string
		for _, res := range reservations.Reservations {
			for _, ins := range res.Instances {
				instanceIds = append(instanceIds, ins.InstanceId)
			}
		}
		if len(instanceIds) < 1 {
			logrus.Info("No instances to terminate")
			return
		}
		_, err = ec2Svc.TerminateInstances(&ec2.TerminateInstancesInput{
			InstanceIds: instanceIds,
		})
		if err != nil {
			logrus.Fatalf("Unable to terminate instanced: %s", err)
		}
		logrus.Infof("Terminated %d instances", len(instanceIds))
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}
