vagrant-build-base-image:
	packer init .
	packer build -force .

## Vagrant test targets (local testing of cloud-init scripts)
#
# Unified Vagrantfile supports two profiles via TEST_PROFILE env var:
#   minipc             - Fresh Ubuntu box/Docker, runs cloud-init/minipc.txt
#   wireguard-prebaked - Pre-built OVA, runs cloud-init/wireguard-prebaked.txt
#
# Usage:
#   make vagrant-wireguard-test          # boot OVA + run prebaked cloud-init
#   make vagrant-wireguard-destroy       # tear down the wireguard test VM
#   make vagrant-wireguard-ssh           # SSH into the wireguard test VM
#   make vagrant-minipc-test             # boot fresh box + run minipc cloud-init
#   make vagrant-minipc-destroy          # tear down the minipc test VM
#   make vagrant-minipc-ssh              # SSH into the minipc test VM
vagrant-wireguard-test:
	WIREGUARD_BOX=${HOME}/.virtualbox/packer/packer_base_ubuntu_26.box TEST_PROFILE=wireguard-prebaked vagrant up --provider=virtualbox

vagrant-wireguard-destroy:
	TEST_PROFILE=wireguard-prebaked vagrant destroy -f

vagrant-wireguard-ssh:
	TEST_PROFILE=wireguard-prebaked vagrant ssh

vagrant-minipc-test:
	TEST_PROFILE=minipc vagrant up

vagrant-minipc-destroy:
	TEST_PROFILE=minipc vagrant destroy -f

vagrant-minipc-ssh:
	TEST_PROFILE=minipc vagrant ssh