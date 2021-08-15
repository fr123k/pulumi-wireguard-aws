package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
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
    userDataEnvVariables := map[string]string{
        "{{ ADMIN_PASSWORD }}":   "ADMIN_PASSWORD",
        "{{ SEED_BRANCH_JOBS }}": "SEED_BRANCH_JOBS",
    }

    userDataSetVariables := map[string]string{
        "{{ AWS_KEY_ID }}":     fmt.Sprintf("AWS_KEY_ID=%s", "undefined"),
        "{{ AWS_KEY_SECRET }}": fmt.Sprintf("AWS_KEY_SECRET=%s", "undefined"),
    }

    userData, err := model.NewUserData("cloud-init/jenkins.yaml", append(model.TemplateVariablesEnvironment(userDataEnvVariables), model.TemplateVariablesString(userDataSetVariables)...))
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

    externalSecurityGroup := model.SecurityGroup{
        Name:        "jenkins-security-group",
        Description: "Pulumi Managed.",
        Tags:        tags,
        IngressRules: []*model.SecurityRule{
            model.AllowOnePortRule("tcp", 80),
            model.AllowOnePortRule("tcp", 22).CidrBlock("95.90.244.46/32"),
            model.AllowSSHRule(security),
        },
        EgressRules: []*model.SecurityRule{
            model.AllowOnePortRule("tcp", 80),
            model.AllowOnePortRule("tcp", 443),
            model.AllowOnePortRule("tcp", 22),
            model.AllowOnePortRule("tcp", 22).CidrBlock("140.82.118.0/24"),
            model.AllowOnePortRule("tcp", 22).CidrBlock("140.82.121.4/32"),
            model.AllowOnePortRule("tcp", 22).CidrBlock("204.232.175.90/32"),
            model.AllowOnePortRule("tcp", 22).CidrBlock("207.97.227.239/32"),
        },
    }
    //The order is important the referenced security groups has to be first.
    computeArgs.SecurityGroups = []*model.SecurityGroup{
        &externalSecurityGroup,
    }

    vm, err := compute.CreateServer(ctx, computeArgs, exports)

    if err != nil {
        return err
    }

    sshConnector := actors.NewSSHConnector(
        actors.SSHConnectorArgs{
            Port:       22,
            Username:   "ubuntu",
            Timeout:    2 * time.Minute,
            SSHKeyPair: *keyPair.SSHKeyPair,
            Commands: []actors.SSHCommand{
                {
                    Command: "sudo cloud-init status --wait",
                    Output: false,
                },
            },
        },
        utility.Logger{
            Ctx: ctx,
        },
    )

    compute.ProvisionVM(ctx, &model.ProvisionArgs{
        ExportName:    "wireguard.publicKey",
        SourceCompute: &model.ComputeResult{
            Compute: vm.Server.CustomResourceState,
        },
    }, &sshConnector)

    return err
    // computeArgs.UserData = userData
}
