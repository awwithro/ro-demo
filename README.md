# ro-demo 

A simple golang app to launch an ec2 instance running rails. It uses go 1.10, the aws-sdk, and cobra for the cli.

## The program

The program comes with two main commands:

1. `launch` launches an instance and probes for it to be ready to serve traffic. It takes about 100s from the command being run to the application being available. It launches instances into the default VPC and uses the default security group which gives it a public IP and reachable on port 80.
2. `cleanup` terminates any running instances that were created by the tool

AWS credentials are needed and the program will pull them from the standard environment variables, AWS_ACCESS_KEY_ID and AWS_SECRET_KEY. The default region used is eu-west-1 but is configurable via the `--region` flag.

## The instance

The instance launched is a t2.micro running Ubuntu 18.04. Docker is installed via user-data and the webapp is then pulled down from docker hub and run.

## The webapp

Just a rails app. Rails 5.2.0 and Ruby 2.5.1. The dockerfile + Gemfile used to build the container is in the `app/` dir. Just rails and the basic gems to serve the welcome page. 