package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launches an ec2 instance",
	Long: `Launches and ec2 instance into the default VPC with the default security-group.
Will poll for the instance to become ready and will check that the endpoint is able to serve traffic.`,
	Run: func(cmd *cobra.Command, args []string) {
		awsSession := session.New(&aws.Config{Region: aws.String(Region)})
		ec2Svc := ec2.New(awsSession)

		if !VpcExists(ec2Svc) {
			CreateVpc(ec2Svc)
		}

		logrus.Info("Creating new instance")
		instance := CreateEC2Instance(ec2Svc)
		logrus.Infof("Created Instance: %v", *instance.InstanceId)
		err := WaitForInstanceToBeLive(instance, ec2Svc)
		if err != nil {
			logrus.Fatal("The instance took to long to be ready")
		}
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)
}
