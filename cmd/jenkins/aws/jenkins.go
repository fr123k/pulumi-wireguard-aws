package main

import (
	"os"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func exports(ctx *pulumi.Context, infra *compute.Infrastructure) {
	ctx.Export("publicIp", infra.Server.PublicIp)
	ctx.Export("publicDns", infra.Server.PublicDns)
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		const username = "pulumi-automation"

		iamUser, err := iam.NewUser(ctx, username+"-user", &iam.UserArgs{
			Tags: pulumi.StringMap{
				"Creator": pulumi.String("jenkins-aws-pulumi"),
				"JobUrl":  pulumi.String(os.Getenv("TRAVIS_JOB_WEB_URL")),
			},
		})

		if err != nil {
			return err
		}

		var s3PolicyContent = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": "s3:*",
				"Resource": "*"
			}
		]
	}`

		var ec2PolicyContent = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Action": "ec2:*",
				"Effect": "Allow",
				"Resource": "*"
			}
		]
	}
	`

		var iamPolicyContent = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Action": "iam:*",
				"Effect": "Allow",
				"Resource": "*"
			}
		]
	}
	`

		s3IamPolicy, err := iam.NewPolicy(ctx, username+"-user-policy-s3", &iam.PolicyArgs{
			Policy: pulumi.String(s3PolicyContent),
		})

		if err != nil {
			return err
		}

		ec2IamPolicy, err := iam.NewPolicy(ctx, username+"-user-policy-ec2", &iam.PolicyArgs{
			Policy: pulumi.String(ec2PolicyContent),
		})

		if err != nil {
			return err
		}

		iamIamPolicy, err := iam.NewPolicy(ctx, username+"-user-policy-iam", &iam.PolicyArgs{
			Policy: pulumi.String(iamPolicyContent),
		})

		if err != nil {
			return err
		}

		iam.NewUserPolicyAttachment(ctx, username+"-user-policy-attachment-s3", &iam.UserPolicyAttachmentArgs{
			User:      iamUser.ID(),
			PolicyArn: s3IamPolicy.ID(),
		})

		iam.NewUserPolicyAttachment(ctx, username+"-user-policy-attachment-ec2", &iam.UserPolicyAttachmentArgs{
			User:      iamUser.ID(),
			PolicyArn: ec2IamPolicy.ID(),
		})

		iam.NewUserPolicyAttachment(ctx, username+"-user-policy-attachment-iam", &iam.UserPolicyAttachmentArgs{
			User:      iamUser.ID(),
			PolicyArn: iamIamPolicy.ID(),
		})

		iamKeys, err := iam.NewAccessKey(ctx, username+"-user-keys", &iam.AccessKeyArgs{
			User: iamUser.ID(),
		})

		if err != nil {
			return err
		}

		ctx.Export("AccessKeys", iamKeys.ID())
		ctx.Export("AccessKeysSecret", iamKeys.Secret)
		return createInfraStructure(ctx)
	})
}

func createInfraStructure(ctx *pulumi.Context) error {
	config := config.New(ctx, "")

	//TODO fetch new created aws key and secret
	// awsKeyID := config.Require("key")
	// awsKeySecret := config.Require("secret")

	userData, err := shared.JenkinsUserData("cloud-init/jenkins.yaml")
	if err != nil {
		return err
	}

	security := model.NewSecurityArgsForVPC(config.GetBool("vpn_enabled_ssh"), model.VPCArgsDefault)
	security.Println()

	vpc, err := network.CreateVPC(ctx, model.VPCArgsDefault)
	if err != nil {
		return err
	}

	keyPairName := "development"
	keyPair := model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
	computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
	computeArgs.Name = "jenkins-master"
	computeArgs.UserData = userData
	computeArgs.Images = []*model.ImageArgs{
		{
			Name:   "ubuntu/images/hvm-ssd/ubuntu-*-18.04-amd64-server-*",
			Owners: []string{"099720109477"},
		},
	}

	tags := map[string]string{
		"JobUrl":         os.Getenv("TRAVIS_JOB_WEB_URL"),
		"Project":        "jenkins",
		"pulumi-managed": "True",
	}

	computeArgs.SecurityGroups = shared.JenkinsSecGroup(tags, security)

	vm, err := compute.CreateServer(ctx, computeArgs, exports)

	if err != nil {
		return err
	}

	sshConnector := shared.JenkinsProvisioner(ctx, keyPair)

	compute.ProvisionVM(ctx, &model.ProvisionArgs{
		ExportName: "wireguard.publicKey",
		SourceCompute: &model.ComputeResult{
			Compute: vm.Server.CustomResourceState,
		},
	}, &sshConnector)

	return err
}
