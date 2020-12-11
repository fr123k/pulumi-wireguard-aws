[![Build Status](https://travis-ci.com/fr123k/pulumi-wireguard-aws.svg?branch=master)](https://travis-ci.com/fr123k/pulumi-wireguard-aws)

# pulumi-wireguard-aws

# AWS Authentication

## With AWS Cli

### Prerequisites

* awscli

Installation for MacOs
```
  brew install awscli
```

### Access/Secret Keys

To setup the access keys for terraform either use the aws cli and run the following command
it will ask for the access key and secret key id and store them in a file `~/.aws/credentials`.
Those one then also picked up by the terraform aws provider.

```
  aws configure
```

## Environment Variables

The values of the following defined environment variables will work for the awscli and the terraform
aws provider and if you put leave a space before the command then they also not appear in the bash
history.

```
 export AWS_ACCESS_KEY_ID=******
 export AWS_SECRET_ACCESS_KEY=******
 export AWS_DEFAULT_REGION=eu-west-1 
```

# Terraform Wireguard

To achieve the best customization to your scenario just fork this repository and adjust the
terraform code to your needs. This whole repository is a result of one weekend work and
properly slightly away from perfection but hopefully a good starting point for running
your own wireguard VPN server in AWS.

**Contribution of course is appreciated.**

## Pulumi

Supported Version: 2.0+
Backend: local

## Prerequisites

* add your ssh rsa public key to the `keys/wireguard.pem.pub` file

## Configuration

### Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| client\_public\_key | The wireguard client public key. | `string` | `"XSGknxaW7PwqiFD061TemUozeTxxafusIRr5dz2fUhw="` | no |
| mailjet\_api\_credentials | The mailjet api credentials in the form API\_KEY:SECRET\_KEY | `string` | `""` | no |
| vpn\_enabled\_ssh | If true the ssh port restricted to the wireguard network range. Otherwise its open for public (0.0.0.0/0). | `bool` | `"true"` | no |

### Client Wireguard Key (client\_public\_key)

- Install the WireGuard tools for your OS: https://www.wireguard.com/install/
- Generate a key pair for the clients
  - `wg genkey | tee client1-privatekey | wg pubkey > client1-publickey`
- Add the desired client ip address and client public key to the variable `wg_client_public_keys` in the 
  `main.tf` file.
  ```
    module "wireguard" {
        source = "./modules/wireguard/"

        ssh_key_id            = aws_key_pair.wireguard.id
        vpc_id                = aws_vpc.wireguard.id
        subnet                = aws_subnet.wireguard
        wg_client_public_keys = [
            {"${cidrhost(aws_subnet.wireguard.cidr_block, 2)/32" = "XSGknxa................................fUhw="},
        ]
    }
  ```
  The `cidrhost` function calculate the client ip address based on the subnet cidr in this example the cidr is `10.8.0.0/24` that results then in the following ip address `10.8.0.2/32` for the client and the value `XSGknxa................................fUhw=`
  is the generated client public key from above

### VPC

- Adjust the `network.tf` file to your needs for example change the subnet cidr.
  ```
    locals {
      network_cidr = "10.61.0.0"
    }
  ```

### SSH (vpn\_enabled\_ssh)

The default is that this variable is set to `true` and therefore the ssh port is only accessible with an established
wireguard VPN connection.

For troubleshooting or debugging purpose it is helpful to access the wireguard virtual machine even without the
need to have an wireguard VPN connection in place. If the wireguard server failed to start or if you can't get the
wireguard server public key without ssh.

To open the ssh port for public access set the value of the `vpn_enabled_ssh` terraform variable to `false`.

```
 export TF_VAR_vpn_enabled_ssh=false
 make create
```

### Mailjet (mailjet\_api\_credentials)

#### Prerequisites

* mailjet account with certified sender address and api keys
* change the sender and recipient address in the `user-data.txt` file

This implementation is a proof of concept to share the wireguard server public elastic ip address and the
public key without the need to have ssh access to the server.

**It's not meant to be used because the sender/recipient address is hardcoded and eMail is not a reliable message format for this kind of information.**

The underlying idea is to push the wireguard connection credentials like
* public elastic ip address
* wireguard public key
to the client with an third party communication channel.

This connection credentials are not static and change always when the wireguard infrastructure is recreated.

```
 export TF_VAR_mailjet_api_credentials=<API_KEY:SECRET_KEY>
 make create
```

## Infrastructure

### Output Variables

| Name | Description |
|------|-------------|
| wireguard\_eips | The list of elastic ip addresses assigned to the wireguard virtual machines. |

### Build

The following command will run terraform validate and terraform plan. The plan is saved in the file
`terraform.plan` so it can applied directly with the command `terraform apply terraform.plan` then
this `terraform apply` command is without the step of creating a extra plan and asking of confirmation
of those as well. The same can be achieved with the `make create` target check the chapter Create.
```
make build
```

Example Output:
```
  terraform validate
  Success! The configuration is valid.

  terraform plan

  An execution plan has been generated and is shown below.
  Resource actions are indicated with the following symbols:
    + create
  <= read (data resources)

  Terraform will perform the following actions:

    # data.aws_route_table.wireguard will be read during apply
    # (config refers to values not yet known)
  <= data "aws_route_table" "wireguard"  {
        + associations   = (known after apply)
        ...
      }

    # aws_internet_gateway.wireguard will be created
    + resource "aws_internet_gateway" "wireguard" {
        + arn      = (known after apply)
        ...
      }

    # aws_key_pair.wireguard will be created
    + resource "aws_key_pair" "wireguard" {
        + arn         = (known after apply)
        ...
      }

    # aws_route.wireguard will be created
    + resource "aws_route" "wireguard" {
        + destination_cidr_block     = "0.0.0.0/0"
        ...
      }

    # aws_subnet.wireguard will be created
    + resource "aws_subnet" "wireguard" {
        + arn                             = (known after apply)
        ...
      }

    # aws_vpc.wireguard will be created
    + resource "aws_vpc" "wireguard" {
        + arn                              = (known after apply)
        ...
      }

    # module.wireguard.data.aws_instances.wireguards will be read during apply
    # (config refers to values not yet known)
  <= data "aws_instances" "wireguards"  {
        + id            = (known after apply)
        ...
      }

    # module.wireguard.aws_autoscaling_group.wireguard_asg will be created
    + resource "aws_autoscaling_group" "wireguard_asg" {
        + arn                       = (known after apply)
        ...
      }

    # module.wireguard.aws_launch_configuration.wireguard_launch_config will be created
    + resource "aws_launch_configuration" "wireguard_launch_config" {
        + arn                         = (known after apply)
        ...
      }

    # module.wireguard.aws_security_group.sg_wireguard_admin will be created
    + resource "aws_security_group" "sg_wireguard_admin" {
        + arn                    = (known after apply)
        ...
      }

    # module.wireguard.aws_security_group.sg_wireguard_external will be created
    + resource "aws_security_group" "sg_wireguard_external" {
        + arn                    = (known after apply)
        ...
      }

  Plan: 9 to add, 0 to change, 0 to destroy.

  Changes to Outputs:
    + wireguard_eips = (known after apply)

  ------------------------------------------------------------------------

  This plan was saved to: terraform.plan

  To perform exactly these actions, run the following command to apply:
      terraform apply "terraform.plan"
```

### Create

The following command will run the make target `build` see the chapter above and then the
`terraform apply terraform.plan` command non-interactive.

**Be aware if the build make target succeed then the infrastructure is created in the next step without asking for confirmation.**
```
  make create
```

Example Output:
```
  terraform apply terraform.plan
  aws_key_pair.wireguard: Creating...
  aws_vpc.wireguard: Creating...
  aws_key_pair.wireguard: Creation complete after 1s [id=wireguard-key]
  aws_vpc.wireguard: Creation complete after 4s [id=vpc-060a4af0251126d20]
  data.aws_route_table.wireguard: Reading...
  aws_internet_gateway.wireguard: Creating...
  aws_subnet.wireguard: Creating...
  module.wireguard.aws_security_group.sg_wireguard_external: Creating...
  data.aws_route_table.wireguard: Read complete after 0s [id=rtb-0f412214655109d05]
  aws_subnet.wireguard: Creation complete after 1s [id=subnet-07251021997cc49b5]
  aws_internet_gateway.wireguard: Creation complete after 2s [id=igw-09d85a0930bb67fe1]
  aws_route.wireguard: Creating...
  aws_route.wireguard: Creation complete after 1s [id=r-rtb-0f412214655109d051080289494]
  module.wireguard.aws_security_group.sg_wireguard_external: Creation complete after 3s [id=sg-01c85ed4aee5f7a60]
  module.wireguard.aws_launch_configuration.wireguard_launch_config: Creating...
  module.wireguard.aws_security_group.sg_wireguard_admin: Creating...
  module.wireguard.aws_launch_configuration.wireguard_launch_config: Creation complete after 1s [id=wireguard-prod-20201206125455575700000001]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Creating...
  module.wireguard.aws_security_group.sg_wireguard_admin: Creation complete after 5s [id=sg-0a7e02598b2e5a11b]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still creating... [10s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still creating... [20s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still creating... [30s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still creating... [40s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Creation complete after 40s [id=wireguard-prod-20201206125455575700000001]
  module.wireguard.data.aws_instances.wireguards: Reading...
  module.wireguard.data.aws_instances.wireguards: Read complete after 1s [id=terraform-20201206125537118500000002]
```

### Clean

The following make target will execute `terraform destroy` and list all the resource it will destroy
and it stops and waits for confirmation by of the user. If the user confirm then it starts to destroy
all the listed resources.

```
  make clean
```

Example Output:
```
  terraform destroy

  An execution plan has been generated and is shown below.
  Resource actions are indicated with the following symbols:
    - destroy

  Terraform will perform the following actions:

    # aws_internet_gateway.wireguard will be destroyed
    - resource "aws_internet_gateway" "wireguard" {
        - arn      = "arn:aws:ec2:eu-west-1:200849096175:internet-gateway/igw-09d85a0930bb67fe1" -> null
        ...
      }

    # aws_key_pair.wireguard will be destroyed
    - resource "aws_key_pair" "wireguard" {
        - arn         = "arn:aws:ec2:eu-west-1:200849096175:key-pair/wireguard-key" -> null
        ...
      }

    # aws_route.wireguard will be destroyed
    - resource "aws_route" "wireguard" {
        - destination_cidr_block = "0.0.0.0/0" -> null
        ...
      }

    # aws_subnet.wireguard will be destroyed
    - resource "aws_subnet" "wireguard" {
        - arn                             = "arn:aws:ec2:eu-west-1:200849096175:subnet/subnet-07251021997cc49b5" -> null
        ...
      }

    # aws_vpc.wireguard will be destroyed
    - resource "aws_vpc" "wireguard" {
        - arn                              = "arn:aws:ec2:eu-west-1:200849096175:vpc/vpc-060a4af0251126d20" -> null
        ...
      }

    # module.wireguard.aws_autoscaling_group.wireguard_asg will be destroyed
    - resource "aws_autoscaling_group" "wireguard_asg" {
        - arn                       = "arn:aws:autoscaling:eu-west-1:200849096175:autoScalingGroup:49375236-e563-4e98-93f3-3fdebd338985:autoScalingGroupName/wireguard-prod-20201206125455575700000001" -> null
        ...
      }

    # module.wireguard.aws_launch_configuration.wireguard_launch_config will be destroyed
    - resource "aws_launch_configuration" "wireguard_launch_config" {
        - arn                         = "arn:aws:autoscaling:eu-west-1:200849096175:launchConfiguration:0f1100bd-f3b0-4f92-b51a-c1fbf37513e4:launchConfigurationName/wireguard-prod-20201206125455575700000001" -> null
        ...
      }

    # module.wireguard.aws_security_group.sg_wireguard_admin will be destroyed
    - resource "aws_security_group" "sg_wireguard_admin" {
        - arn                    = "arn:aws:ec2:eu-west-1:200849096175:security-group/sg-0a7e02598b2e5a11b" -> null
        ...
      }

    # module.wireguard.aws_security_group.sg_wireguard_external will be destroyed
    - resource "aws_security_group" "sg_wireguard_external" {
        - arn                    = "arn:aws:ec2:eu-west-1:200849096175:security-group/sg-01c85ed4aee5f7a60" -> null
        ...
      }

  Plan: 0 to add, 0 to change, 9 to destroy.

  Changes to Outputs:
    - wireguard_eips = [
        - "34.245.64.222",
      ] -> null

  Do you really want to destroy all resources?
    Terraform will destroy all your managed infrastructure, as shown above.
    There is no undo. Only 'yes' will be accepted to confirm.

    Enter a value: yes

  aws_route.wireguard: Destroying... [id=r-rtb-0f412214655109d051080289494]
  module.wireguard.aws_security_group.sg_wireguard_admin: Destroying... [id=sg-0a7e02598b2e5a11b]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Destroying... [id=wireguard-prod-20201206125455575700000001]
  aws_route.wireguard: Destruction complete after 0s
  aws_internet_gateway.wireguard: Destroying... [id=igw-09d85a0930bb67fe1]
  module.wireguard.aws_security_group.sg_wireguard_admin: Destruction complete after 0s
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still destroying... [id=wireguard-prod-20201206125455575700000001, 10s elapsed]
  aws_internet_gateway.wireguard: Still destroying... [id=igw-09d85a0930bb67fe1, 10s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still destroying... [id=wireguard-prod-20201206125455575700000001, 20s elapsed]
  aws_internet_gateway.wireguard: Still destroying... [id=igw-09d85a0930bb67fe1, 20s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still destroying... [id=wireguard-prod-20201206125455575700000001, 30s elapsed]
  aws_internet_gateway.wireguard: Still destroying... [id=igw-09d85a0930bb67fe1, 30s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still destroying... [id=wireguard-prod-20201206125455575700000001, 40s elapsed]
  aws_internet_gateway.wireguard: Still destroying... [id=igw-09d85a0930bb67fe1, 40s elapsed]
  aws_internet_gateway.wireguard: Destruction complete after 48s
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still destroying... [id=wireguard-prod-20201206125455575700000001, 50s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still destroying... [id=wireguard-prod-20201206125455575700000001, 1m0s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still destroying... [id=wireguard-prod-20201206125455575700000001, 1m10s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Still destroying... [id=wireguard-prod-20201206125455575700000001, 1m20s elapsed]
  module.wireguard.aws_autoscaling_group.wireguard_asg: Destruction complete after 1m21s
  aws_subnet.wireguard: Destroying... [id=subnet-07251021997cc49b5]
  module.wireguard.aws_launch_configuration.wireguard_launch_config: Destroying... [id=wireguard-prod-20201206125455575700000001]
  module.wireguard.aws_launch_configuration.wireguard_launch_config: Destruction complete after 1s
  aws_key_pair.wireguard: Destroying... [id=wireguard-key]
  module.wireguard.aws_security_group.sg_wireguard_external: Destroying... [id=sg-01c85ed4aee5f7a60]
  aws_subnet.wireguard: Destruction complete after 1s
  aws_key_pair.wireguard: Destruction complete after 0s
  module.wireguard.aws_security_group.sg_wireguard_external: Destruction complete after 0s
  aws_vpc.wireguard: Destroying... [id=vpc-060a4af0251126d20]
  aws_vpc.wireguard: Destruction complete after 1s

  Destroy complete! Resources: 9 destroyed.
```

### Recreate

The following make target is an convenient target.
```
  make recreate
```

It's just combine the `clean` and `create` make target and the full form looks like this. 

```
  make clean create
```

### Shell

#### Single Wireguard Server
If you have one wireguard virtual machine created as part of the auto scaling group then just run
the following command to establish an ssh connection.
```
  make shell
```
The previous `make shell` is a shortcut make target and the full make target looks like this. 
```
  SERVER_INDEX=0 make shell
```

#### Multiple Wireguard Server's

If you have two wireguard virtual machine created as part of the auto scaling group then by specifying
the `SERVER_INDEX=1` you able to access the second virtual machine with ssh.
```
  SERVER_INDEX=1 make shell
```
# Changes

* setup travis build
* send the wireguard elastic ip address and its public key via email with mailjet (optional)
* securing the ssh port with wireguard VPN (optional)

# Todos

* support terraform workspaces to isolate travis build from local builds
* build and AWS AMI image for wireguard (use packer for this maybe as part of this repo?)
* make the email sending via Mailjet optional and pass it from outside the wireguard module
* configure the client ip addresses and public keys outside of terraform so that a change doesn't need a full recreation of the wireguard VM
  only a restart of the wireguard systemd service would be needed.
