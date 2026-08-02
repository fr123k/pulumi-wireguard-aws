module github.com/fr123k/pulumi-wireguard-aws/packer

go 1.26.4

// wireguard-ui is fetched from the fr123k fork which keeps the upstream
// module path github.com/ngoduykhanh/wireguard-ui.
replace github.com/ngoduykhanh/wireguard-ui => github.com/fr123k/wireguard-ui v0.3.8

require (
	github.com/containifyci/secret-operator v0.6.1
	github.com/ngoduykhanh/wireguard-ui v0.3.8
)
