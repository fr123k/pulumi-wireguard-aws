[![Build Status](https://travis-ci.com/fr123k/pulumi-wireguard-aws.svg?branch=main)](https://travis-ci.com/fr123k/pulumi-wireguard-aws)

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

To setup the access keys for pulumi either use the aws cli and run the following command
it will ask for the access key and secret key id and store them in a file `~/.aws/credentials`.
Those one then also picked up by the pulumi aws provider.

```
  aws configure
```

## Environment Variables

The values of the following defined environment variables will work for the awscli and the pulumi
aws provider and if you put leave a space before the command then they also not appear in the bash
history.

```
 export AWS_ACCESS_KEY_ID=******
 export AWS_SECRET_ACCESS_KEY=******
 export AWS_DEFAULT_REGION=eu-west-1 
```

# Pulumi Wireguard

To achieve the best customization to your scenario just fork this repository and adjust the
pulumi code to your needs. This whole repository is a result of one weekend work and
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
  - `wg genkey | tee client-privatekey | wg pubkey > client-publickey`
- Specify the client public key and the desired client ip address as environment variables
- Add the desired client ip address and client public key to the variable `wg_client_public_keys` in the 
  `main.tf` file.
  ```
    export CLIENT_ID_ADDRESS=10.8.0.2
    export  CLIENT_PUBLICKEY="XSGknxa................................fUhw="
  ```
  The CLIENT_PUBLICKEY value `XSGknxa................................fUhw=` is the generated client public key from above

### SSH (vpn\_enabled\_ssh)

The default is that this variable is set to `true` and therefore the ssh port is only accessible with an established
wireguard VPN connection.

For troubleshooting or debugging purpose it is helpful to access the wireguard virtual machine even without the
need to have an wireguard VPN connection in place. If the wireguard server failed to start or if you can't get the
wireguard server public key without ssh.

To open the ssh port for public access set the value of the `vpn_enabled_ssh` pulumi variable to `false`.

```
 export VPN_ENABLED_SSH=false
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
 export MAILJET_API_CREDENTIALS=<API_KEY:SECRET_KEY>
 make create
```

## Infrastructure

### Output Variables

| Name | Description | Cloud Provider |
|------|-------------|----------------|
| publicIp | The elastic ip address assigned to the wireguard virtual machine. | AWS, Hetzner
| publicDns | The public FQDN of the wireguard virtual machine. | AWS, Hetzner
| cloud-init | The rendered cloud-init template provided as userdata to wireguard virtual machine. | AWS, Hetzner
| vpcId | The id of the generated VPC resource. | AWS
| subnetId | The id of the generated subnet resource. | AWS, Hetzner


### Build

The following command will build the pulumi wireguard binary and run the unit tests.
The result of this build step is the standalone pulumi binary in the `./build/` folder. 

The same can be achieved with the `make create` target check the chapter Create.
```
make build
```

Example Output:
```
  echo -e 'n\n' | ssh-keygen -t rsa -b 4096 -q -N "" -f ./keys/wireguard.pem || true
  ./keys/wireguard.pem already exists.
  Overwrite (y/n)? echo "No"
  No
  go build -o build/wireguard-aws cmd/wireguard/aws/wireguard.go
  go test -v --cover ./...
  ?   	github.com/fr123k/pulumi-wireguard-aws/cmd/wireguard/aws	[no test files]
  === RUN   TestVpcArg
  --- PASS: TestVpcArg (0.00s)
  PASS
  coverage: 100.0% of statements
  ok  	github.com/fr123k/pulumi-wireguard-aws/cmd/wireguard/config	(cached)
  ...
  === RUN   TestTemplateVariablesString
  Key: {{ TEST_CLIENT_PUBLICKEY }} Value: TEST_CLIENT_PUBLICKEY
  Key: {{ TEST_METADATA_URL }} Value: TEST_METADATA_URL
  --- PASS: TestTemplateVariablesString (0.00s)
  ...
  ok  	github.com/fr123k/pulumi-wireguard-aws/pkg/model	(cached)	coverage: 86.5% of statements
  === RUN   TestReadFileWithNonExistingFile
  --- PASS: TestReadFileWithNonExistingFile (0.00s)
  === RUN   TestReadFileFromMemory
  --- PASS: TestReadFileFromMemory (0.00s)
  PASS
  coverage: 66.7% of statements
  ok  	github.com/fr123k/pulumi-wireguard-aws/pkg/utility	(cached)	coverage: 66.7% of statements
  ln -fs wireguard-aws ./build/wireguard

```

### Create

The following command will run the make target `build` see the chapter above and then the
`pulumi up` command non-interactive.

**Be aware if the build make target succeed then the infrastructure is created in the next step without asking for confirmation.**
```
  make create
```

Example Output:
```
  echo -e 'n\n' | ssh-keygen -t rsa -b 4096 -q -N "" -f ./keys/wireguard.pem || true
  ./keys/wireguard.pem already exists.
  Overwrite (y/n)? echo "No"
  No
  go build -o build/wireguard-aws cmd/wireguard/aws/wireguard.go
  go test -v --cover ./...
  ...
  coverage: 66.7% of statements
  ok  	github.com/fr123k/pulumi-wireguard-aws/pkg/utility	(cached)	coverage: 66.7% of statements
  ln -fs wireguard-aws ./build/wireguard
  pulumi plugin install resource aws 3.21.0
  pulumi plugin ls
  NAME    KIND      VERSION  SIZE    INSTALLED    LAST USED
  aws     resource  3.21.0   253 MB  1 week ago   1 week ago
  aws     resource  3.19.3   253 MB  2 weeks ago  2 weeks ago
  hcloud  resource  0.4.0    44 MB   1 week ago   1 week ago

  TOTAL plugin cache size: 550 MB
  pulumi login --local
  Logged in to local as (file://~)
  # pulumi login --cloud-url s3://s3-pulumi-state-d12f2f1
  # pulumi stack rm -f aws
  # pulumi stack select aws
  pulumi stack select -c aws
  pulumi config set aws:region eu-west-1
  pulumi config set vpn_enabled_ssh true
  pulumi up --yes
  Previewing update (aws):
      Type                        Name                      Plan       Info
  +   pulumi:pulumi:Stack         wireguard-aws-pulumi-aws  create     9 messages
  +   ├─ aws:ec2:Vpc              wireguard                 create
  +   ├─ aws:ec2:InternetGateway  wireguard                 create
  +   ├─ aws:ec2:Subnet           wireguard                 create
  +   ├─ aws:ec2:SecurityGroup    wireguard-external        create
  +   ├─ aws:ec2:Route            wireguard                 create
  +   ├─ aws:ec2:SecurityGroup    wireguard-admin           create
  +   ├─ aws:ec2:KeyPair          wireguard                 create
  +   └─ aws:ec2:Instance         wireguard                 create
  
  Diagnostics:
    pulumi:pulumi:Stack (wireguard-aws-pulumi-aws):
      {"VPNEnabledSSH":true,"VPNCidr":"10.8.0.0/24"}
      Key: {{ CLIENT_IP_ADDRESS }} Value: CLIENT_IP_ADDRESS
      Key: {{ MAILJET_API_CREDENTIALS }} Value: MAILJET_API_CREDENTIALS
      Key: {{ METADATA_URL }} Value: METADATA_URL
      Key: {{ CLIENT_PUBLICKEY }} Value: CLIENT_PUBLICKEY
      Key: {{ CLIENT_IP_ADDRESS }} Value: CLIENT_IP_ADDRESS
      Key: {{ MAILJET_API_CREDENTIALS }} Value: MAILJET_API_CREDENTIALS
      Key: {{ METADATA_URL }} Value: METADATA_URL
      Key: {{ CLIENT_PUBLICKEY }} Value: CLIENT_PUBLICKEY
  

  Permalink: file:///Users/franki/.pulumi/stacks/aws.json
  Updating (aws):
      Type                        Name                      Status      Info
  +   pulumi:pulumi:Stack         wireguard-aws-pulumi-aws  created     9 messages
  +   ├─ aws:ec2:Vpc              wireguard                 created
  +   ├─ aws:ec2:KeyPair          wireguard                 created
  +   ├─ aws:ec2:InternetGateway  wireguard                 created
  +   ├─ aws:ec2:Subnet           wireguard                 created
  +   ├─ aws:ec2:SecurityGroup    wireguard-external        created
  +   ├─ aws:ec2:Route            wireguard                 created
  +   ├─ aws:ec2:SecurityGroup    wireguard-admin           created
  +   └─ aws:ec2:Instance         wireguard                 created
  
  Diagnostics:
    pulumi:pulumi:Stack (wireguard-aws-pulumi-aws):
      {"VPNEnabledSSH":true,"VPNCidr":"10.8.0.0/24"}
      Key: {{ CLIENT_PUBLICKEY }} Value: CLIENT_PUBLICKEY
      Key: {{ CLIENT_IP_ADDRESS }} Value: CLIENT_IP_ADDRESS
      Key: {{ MAILJET_API_CREDENTIALS }} Value: MAILJET_API_CREDENTIALS
      Key: {{ METADATA_URL }} Value: METADATA_URL
      Key: {{ CLIENT_PUBLICKEY }} Value: CLIENT_PUBLICKEY
      Key: {{ CLIENT_IP_ADDRESS }} Value: CLIENT_IP_ADDRESS
      Key: {{ MAILJET_API_CREDENTIALS }} Value: MAILJET_API_CREDENTIALS
      Key: {{ METADATA_URL }} Value: METADATA_URL
  
  Outputs:
      cloud-init: "#!/bin/bash -v\n\napt-get update -y\napt-get upgrade -y\napt-get install -y wireguard-dkms wireguard-tools \n\numask 077\n#TODO make server public key available outside the vm instance\nwg genkey | tee /tmp/server_privatekey | wg pubkey > /tmp/server_publickey\n\nMYV4IP=$(curl )\n\ncat > /etc/wireguard/wg0.conf <<- EOF\n[Interface]\nAddress = $MYV4IP/24\nPrivateKey = $(cat /tmp/server_privatekey)\nListenPort = 51820\nPostUp   = iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE\nPostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE\n\n[Peer]\nPublicKey = \"XSG................................Uhw=\"\nAllowedIPs = 10.8.0.2/32\nPersistentKeepalive = 25\nEOF\n\nchown -R root:root /etc/wireguard/\nchmod -R og-rwx /etc/wireguard/*\nsed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/' /etc/sysctl.conf\nsysctl -p\nufw allow ssh\nufw allow 51820/udp\nufw --force enable\nsystemctl enable wg-quick@wg0.service\nsystemctl start wg-quick@wg0.service\n\nMAILJET_AUTH=\"\"\n\nif [ \"$MAILJET_AUTH\" != \"\" ]; then\n\n    # TODO make the list of emails configurable per client ip\n    cat > /tmp/wireguard.email <<- EOF\n    {\n    \"Messages\":[\n        {\n        \"From\": {\n            \"Email\": \"wireguard@fr123k.uk\",\n            \"Name\": \"Wireguard $MYV4IP\"\n        },\n        \"To\": [\n            {\n            \"Email\": \"fr12_k@yahoo.com\",\n            \"Name\": \"Frank\"\n            }\n        ],\n        \"Subject\": \"Wireguard publickey\",\n        \"TextPart\": \"The wireguard public key is $(cat /tmp/server_publickey) and the ip address $MYV4IP\",\n        \"CustomID\": \"Wireguard Publickey\"\n        }\n    ]\n    }\nEOF\n\n    curl -s -X POST \\\n    --user \"${mailjet_api_credentials}\" \\\n    https://api.mailjet.com/v3.1/send \\\n    -H 'Content-Type: application/json' \\\n    --data \"@/tmp/wireguard.email\"\nfi\n"
      publicDns : "ec2-3-249-99-113.eu-west-1.compute.amazonaws.com"
      publicIp  : "3.249.99.113"
      subnetId  : "subnet-0957988f3c093db5b"
      vpcId     : "vpc-0be1174cd9bafc205"

  Resources:
      + 9 created

  Duration: 1m37s
```

### Clean

The following make target will execute `pulumi destroy` and list all the resource it will destroy and start immediately to destroy them. 
**No confirmation needed**
```
  make clean
```

Example Output:
```
  pulumi destroy --yes -s aws || true
  Previewing destroy (aws):
      Type                        Name                      Plan
  -   pulumi:pulumi:Stack         wireguard-aws-pulumi-aws  delete
  -   ├─ aws:ec2:Instance         wireguard                 delete
  -   ├─ aws:ec2:Route            wireguard                 delete
  -   ├─ aws:ec2:SecurityGroup    wireguard-admin           delete
  -   ├─ aws:ec2:Subnet           wireguard                 delete
  -   ├─ aws:ec2:InternetGateway  wireguard                 delete
  -   ├─ aws:ec2:SecurityGroup    wireguard-external        delete
  -   ├─ aws:ec2:KeyPair          wireguard                 delete
  -   └─ aws:ec2:Vpc              wireguard                 delete
  
  Outputs:
    - cloud-init: "#!/bin/bash -v\n\napt-get update -y\napt-get upgrade -y\napt-get install -y wireguard-dkms wireguard-tools \n\numask 077\n#TODO make server public key available outside the vm instance\nwg genkey | tee /tmp/server_privatekey | wg pubkey > /tmp/server_publickey\n\nMYV4IP=$(curl )\n\ncat > /etc/wireguard/wg0.conf <<- EOF\n[Interface]\nAddress = $MYV4IP/24\nPrivateKey = $(cat /tmp/server_privatekey)\nListenPort = 51820\nPostUp   = iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE\nPostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE\n\n[Peer]\nPublicKey = \"XS................................hw=\"\nAllowedIPs = 10.8.0.2/32\nPersistentKeepalive = 25\nEOF\n\nchown -R root:root /etc/wireguard/\nchmod -R og-rwx /etc/wireguard/*\nsed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/' /etc/sysctl.conf\nsysctl -p\nufw allow ssh\nufw allow 51820/udp\nufw --force enable\nsystemctl enable wg-quick@wg0.service\nsystemctl start wg-quick@wg0.service\n\nMAILJET_AUTH=\"\"\n\nif [ \"$MAILJET_AUTH\" != \"\" ]; then\n\n    # TODO make the list of emails configurable per client ip\n    cat > /tmp/wireguard.email <<- EOF\n    {\n    \"Messages\":[\n        {\n        \"From\": {\n            \"Email\": \"wireguard@fr123k.uk\",\n            \"Name\": \"Wireguard $MYV4IP\"\n        },\n        \"To\": [\n            {\n            \"Email\": \"fr12_k@yahoo.com\",\n            \"Name\": \"Frank\"\n            }\n        ],\n        \"Subject\": \"Wireguard publickey\",\n        \"TextPart\": \"The wireguard public key is $(cat /tmp/server_publickey) and the ip address $MYV4IP\",\n        \"CustomID\": \"Wireguard Publickey\"\n        }\n    ]\n    }\nEOF\n\n    curl -s -X POST \\\n    --user \"${mailjet_api_credentials}\" \\\n    https://api.mailjet.com/v3.1/send \\\n    -H 'Content-Type: application/json' \\\n    --data \"@/tmp/wireguard.email\"\nfi\n"
    - publicDns : "ec2-3-248-222-248.eu-west-1.compute.amazonaws.com"
    - publicIp  : "3.248.222.248"
    - subnetId  : "subnet-04c22ab4638cab72c"
    - vpcId     : "vpc-0f3c1a362758cb265"

  Resources:
      - 9 to delete

  Permalink: file:///Users/franki/.pulumi/stacks/aws.json
  Destroying (aws):
      Type                        Name                      Status
  -   pulumi:pulumi:Stack         wireguard-aws-pulumi-aws  deleted
  -   ├─ aws:ec2:Instance         wireguard                 deleted
  -   ├─ aws:ec2:Route            wireguard                 deleted
  -   ├─ aws:ec2:SecurityGroup    wireguard-admin           deleted
  -   ├─ aws:ec2:InternetGateway  wireguard                 deleted
  -   ├─ aws:ec2:Subnet           wireguard                 deleted
  -   ├─ aws:ec2:SecurityGroup    wireguard-external        deleted
  -   ├─ aws:ec2:KeyPair          wireguard                 deleted
  -   └─ aws:ec2:Vpc              wireguard                 deleted
  
  Outputs:
    - cloud-init: "#!/bin/bash -v\n\napt-get update -y\napt-get upgrade -y\napt-get install -y wireguard-dkms wireguard-tools \n\numask 077\n#TODO make server public key available outside the vm instance\nwg genkey | tee /tmp/server_privatekey | wg pubkey > /tmp/server_publickey\n\nMYV4IP=$(curl )\n\ncat > /etc/wireguard/wg0.conf <<- EOF\n[Interface]\nAddress = $MYV4IP/24\nPrivateKey = $(cat /tmp/server_privatekey)\nListenPort = 51820\nPostUp   = iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE\nPostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE\n\n[Peer]\nPublicKey = \"XSGknxaW7PwqiFD061TemUozeTxxafusIRr5dz2fUhw=\"\nAllowedIPs = {{ CLIENT_IP_ADDRESS }}/32\nPersistentKeepalive = 25\nEOF\n\nchown -R root:root /etc/wireguard/\nchmod -R og-rwx /etc/wireguard/*\nsed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/' /etc/sysctl.conf\nsysctl -p\nufw allow ssh\nufw allow 51820/udp\nufw --force enable\nsystemctl enable wg-quick@wg0.service\nsystemctl start wg-quick@wg0.service\n\nMAILJET_AUTH=\"\"\n\nif [ \"$MAILJET_AUTH\" != \"\" ]; then\n\n    # TODO make the list of emails configurable per client ip\n    cat > /tmp/wireguard.email <<- EOF\n    {\n    \"Messages\":[\n        {\n        \"From\": {\n            \"Email\": \"wireguard@fr123k.uk\",\n            \"Name\": \"Wireguard $MYV4IP\"\n        },\n        \"To\": [\n            {\n            \"Email\": \"fr12_k@yahoo.com\",\n            \"Name\": \"Frank\"\n            }\n        ],\n        \"Subject\": \"Wireguard publickey\",\n        \"TextPart\": \"The wireguard public key is $(cat /tmp/server_publickey) and the ip address $MYV4IP\",\n        \"CustomID\": \"Wireguard Publickey\"\n        }\n    ]\n    }\nEOF\n\n    curl -s -X POST \\\n    --user \"${mailjet_api_credentials}\" \\\n    https://api.mailjet.com/v3.1/send \\\n    -H 'Content-Type: application/json' \\\n    --data \"@/tmp/wireguard.email\"\nfi\n"
    - publicDns : "ec2-3-248-222-248.eu-west-1.compute.amazonaws.com"
    - publicIp  : "3.248.222.248"
    - subnetId  : "subnet-04c22ab4638cab72c"
    - vpcId     : "vpc-0f3c1a362758cb265"

  Resources:
      - 9 deleted

  Duration: 1m5s
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

To open a SSH shell just run the following command.
```
  make shell
```

# Changes

* setup travis build
* send the wireguard elastic ip address and its public key via email with mailjet (optional)
* securing the ssh port with wireguard VPN (optional default is true)
* make the email sending via Mailjet optional and pass it from outside the wireguard module

# Todos

* build and AWS AMI image for wireguard (use packer for this maybe as part of this repo? or use pulumi to build an AMI)
* configure the client ip addresses and public keys outside of pulumi so that a change doesn't need a full recreation of the wireguard VM
  only a restart of the wireguard systemd service would be needed.
