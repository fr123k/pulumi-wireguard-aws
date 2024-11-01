package sdk

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func StopInstance(instanceId string) string {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)
	instanceIDsToStop := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			&instanceId,
		},
	}

	// Send request to stop instances
	result, err := svc.StopInstances(instanceIDsToStop)
	if err != nil {
		panic(fmt.Errorf("Failed to stop instance '%s' : %s", instanceId, err))
	}
	// Use a waiter function to wait until the instances are stopped
	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIDsToStop.InstanceIds,
	}
	if err := svc.WaitUntilInstanceStopped(describeInstancesInput); err != nil {
		panic(fmt.Errorf("Failed to wait for instance '%s' to be stopped: %s", instanceId, err))
	}
	fmt.Printf("Instance are '%s' stopped", instanceId)
	return result.String()
}

func CopyImage(amiId string) string {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess, aws.NewConfig())

	input := &ec2.CopyImageInput{
		Name:          aws.String("wireguard-ami-new-copy"),
		SourceImageId: aws.String(amiId),
		SourceRegion:  aws.String(*sess.Config.Region),
	}

	result, err := svc.CopyImage(input)
	if err != nil {
		panic(fmt.Errorf("Failed to copy ami '%s' : %s", amiId, err))
	}

	fmt.Printf("AMI is '%s' copied", amiId)
	return result.String()
}
