// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package ec2

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Get characteristics for a single EC2 Instance Type.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
//
//	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			_, err := ec2.GetInstanceType(ctx, &ec2.GetInstanceTypeArgs{
//				InstanceType: "t2.micro",
//			}, nil)
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
func GetInstanceType(ctx *pulumi.Context, args *GetInstanceTypeArgs, opts ...pulumi.InvokeOption) (*GetInstanceTypeResult, error) {
	opts = internal.PkgInvokeDefaultOpts(opts)
	var rv GetInstanceTypeResult
	err := ctx.Invoke("aws:ec2/getInstanceType:getInstanceType", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getInstanceType.
type GetInstanceTypeArgs struct {
	// Instance
	InstanceType string `pulumi:"instanceType"`
}

// A collection of values returned by getInstanceType.
type GetInstanceTypeResult struct {
	// `true` if auto recovery is supported.
	AutoRecoverySupported bool `pulumi:"autoRecoverySupported"`
	// `true` if it is a bare metal instance type.
	BareMetal bool `pulumi:"bareMetal"`
	// `true` if the instance type is a burstable performance instance type.
	BurstablePerformanceSupported bool `pulumi:"burstablePerformanceSupported"`
	// `true`  if the instance type is a current generation.
	CurrentGeneration bool `pulumi:"currentGeneration"`
	// `true` if Dedicated Hosts are supported on the instance type.
	DedicatedHostsSupported bool `pulumi:"dedicatedHostsSupported"`
	// Default number of cores for the instance type.
	DefaultCores int `pulumi:"defaultCores"`
	// The  default  number of threads per core for the instance type.
	DefaultThreadsPerCore int `pulumi:"defaultThreadsPerCore"`
	// Default number of vCPUs for the instance type.
	DefaultVcpus int `pulumi:"defaultVcpus"`
	// Indicates whether Amazon EBS encryption is supported.
	EbsEncryptionSupport string `pulumi:"ebsEncryptionSupport"`
	// Whether non-volatile memory express (NVMe) is supported.
	EbsNvmeSupport string `pulumi:"ebsNvmeSupport"`
	// Indicates that the instance type is Amazon EBS-optimized.
	EbsOptimizedSupport string `pulumi:"ebsOptimizedSupport"`
	// The baseline bandwidth performance for an EBS-optimized instance type, in Mbps.
	EbsPerformanceBaselineBandwidth int `pulumi:"ebsPerformanceBaselineBandwidth"`
	// The baseline input/output storage operations per seconds for an EBS-optimized instance type.
	EbsPerformanceBaselineIops int `pulumi:"ebsPerformanceBaselineIops"`
	// The baseline throughput performance for an EBS-optimized instance type, in MBps.
	EbsPerformanceBaselineThroughput float64 `pulumi:"ebsPerformanceBaselineThroughput"`
	// The maximum bandwidth performance for an EBS-optimized instance type, in Mbps.
	EbsPerformanceMaximumBandwidth int `pulumi:"ebsPerformanceMaximumBandwidth"`
	// The maximum input/output storage operations per second for an EBS-optimized instance type.
	EbsPerformanceMaximumIops int `pulumi:"ebsPerformanceMaximumIops"`
	// The maximum throughput performance for an EBS-optimized instance type, in MBps.
	EbsPerformanceMaximumThroughput float64 `pulumi:"ebsPerformanceMaximumThroughput"`
	// Whether Elastic Fabric Adapter (EFA) is supported.
	EfaSupported bool `pulumi:"efaSupported"`
	// Whether Elastic Network Adapter (ENA) is supported.
	EnaSupport string `pulumi:"enaSupport"`
	// Indicates whether encryption in-transit between instances is supported.
	EncryptionInTransitSupported bool `pulumi:"encryptionInTransitSupported"`
	// Describes the FPGA accelerator settings for the instance type.
	// * `fpgas.#.count` - The count of FPGA accelerators for the instance type.
	// * `fpgas.#.manufacturer` - The manufacturer of the FPGA accelerator.
	// * `fpgas.#.memory_size` - The size (in MiB) for the memory available to the FPGA accelerator.
	// * `fpgas.#.name` - The name of the FPGA accelerator.
	Fpgas []GetInstanceTypeFpga `pulumi:"fpgas"`
	// `true` if the instance type is eligible for the free tier.
	FreeTierEligible bool `pulumi:"freeTierEligible"`
	// Describes the GPU accelerators for the instance type.
	// * `gpus.#.count` - The number of GPUs for the instance type.
	// * `gpus.#.manufacturer` - The manufacturer of the GPU accelerator.
	// * `gpus.#.memory_size` - The size (in MiB) for the memory available to the GPU accelerator.
	// * `gpus.#.name` - The name of the GPU accelerator.
	Gpuses []GetInstanceTypeGpus `pulumi:"gpuses"`
	// `true` if On-Demand hibernation is supported.
	HibernationSupported bool `pulumi:"hibernationSupported"`
	// Hypervisor used for the instance type.
	Hypervisor string `pulumi:"hypervisor"`
	// The provider-assigned unique ID for this managed resource.
	Id string `pulumi:"id"`
	// Describes the Inference accelerators for the instance type.
	// * `inference_accelerators.#.count` - The number of Inference accelerators for the instance type.
	// * `inference_accelerators.#.manufacturer` - The manufacturer of the Inference accelerator.
	// * `inference_accelerators.#.name` - The name of the Inference accelerator.
	InferenceAccelerators []GetInstanceTypeInferenceAccelerator `pulumi:"inferenceAccelerators"`
	// Describes the disks for the instance type.
	// * `instance_disks.#.count` - The number of disks with this configuration.
	// * `instance_disks.#.size` - The size of the disk in GB.
	// * `instance_disks.#.type` - The type of disk.
	InstanceDisks []GetInstanceTypeInstanceDisk `pulumi:"instanceDisks"`
	// `true` if instance storage is supported.
	InstanceStorageSupported bool   `pulumi:"instanceStorageSupported"`
	InstanceType             string `pulumi:"instanceType"`
	// `true` if IPv6 is supported.
	Ipv6Supported bool `pulumi:"ipv6Supported"`
	// The maximum number of IPv4 addresses per network interface.
	MaximumIpv4AddressesPerInterface int `pulumi:"maximumIpv4AddressesPerInterface"`
	// The maximum number of IPv6 addresses per network interface.
	MaximumIpv6AddressesPerInterface int `pulumi:"maximumIpv6AddressesPerInterface"`
	// The maximum number of physical network cards that can be allocated to the instance.
	MaximumNetworkCards int `pulumi:"maximumNetworkCards"`
	// The maximum number of network interfaces for the instance type.
	MaximumNetworkInterfaces int `pulumi:"maximumNetworkInterfaces"`
	// Size of the instance memory, in MiB.
	MemorySize int `pulumi:"memorySize"`
	// Describes the network performance.
	NetworkPerformance string `pulumi:"networkPerformance"`
	// A list of architectures supported by the instance type.
	SupportedArchitectures []string `pulumi:"supportedArchitectures"`
	// A list of supported placement groups types.
	SupportedPlacementStrategies []string `pulumi:"supportedPlacementStrategies"`
	// Indicates the supported root device types.
	SupportedRootDeviceTypes []string `pulumi:"supportedRootDeviceTypes"`
	// Indicates whether the instance type is offered for spot or On-Demand.
	SupportedUsagesClasses []string `pulumi:"supportedUsagesClasses"`
	// The supported virtualization types.
	SupportedVirtualizationTypes []string `pulumi:"supportedVirtualizationTypes"`
	// The speed of the processor, in GHz.
	SustainedClockSpeed float64 `pulumi:"sustainedClockSpeed"`
	// Total memory of all FPGA accelerators for the instance type (in MiB).
	TotalFpgaMemory int `pulumi:"totalFpgaMemory"`
	// Total size of the memory for the GPU accelerators for the instance type (in MiB).
	TotalGpuMemory int `pulumi:"totalGpuMemory"`
	// The total size of the instance disks, in GB.
	TotalInstanceStorage int `pulumi:"totalInstanceStorage"`
	// List of the valid number of cores that can be configured for the instance type.
	ValidCores []int `pulumi:"validCores"`
	// List of the valid number of threads per core that can be configured for the instance type.
	ValidThreadsPerCores []int `pulumi:"validThreadsPerCores"`
}

func GetInstanceTypeOutput(ctx *pulumi.Context, args GetInstanceTypeOutputArgs, opts ...pulumi.InvokeOption) GetInstanceTypeResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (GetInstanceTypeResultOutput, error) {
			args := v.(GetInstanceTypeArgs)
			opts = internal.PkgInvokeDefaultOpts(opts)
			var rv GetInstanceTypeResult
			secret, err := ctx.InvokePackageRaw("aws:ec2/getInstanceType:getInstanceType", args, &rv, "", opts...)
			if err != nil {
				return GetInstanceTypeResultOutput{}, err
			}

			output := pulumi.ToOutput(rv).(GetInstanceTypeResultOutput)
			if secret {
				return pulumi.ToSecret(output).(GetInstanceTypeResultOutput), nil
			}
			return output, nil
		}).(GetInstanceTypeResultOutput)
}

// A collection of arguments for invoking getInstanceType.
type GetInstanceTypeOutputArgs struct {
	// Instance
	InstanceType pulumi.StringInput `pulumi:"instanceType"`
}

func (GetInstanceTypeOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*GetInstanceTypeArgs)(nil)).Elem()
}

// A collection of values returned by getInstanceType.
type GetInstanceTypeResultOutput struct{ *pulumi.OutputState }

func (GetInstanceTypeResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*GetInstanceTypeResult)(nil)).Elem()
}

func (o GetInstanceTypeResultOutput) ToGetInstanceTypeResultOutput() GetInstanceTypeResultOutput {
	return o
}

func (o GetInstanceTypeResultOutput) ToGetInstanceTypeResultOutputWithContext(ctx context.Context) GetInstanceTypeResultOutput {
	return o
}

// `true` if auto recovery is supported.
func (o GetInstanceTypeResultOutput) AutoRecoverySupported() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.AutoRecoverySupported }).(pulumi.BoolOutput)
}

// `true` if it is a bare metal instance type.
func (o GetInstanceTypeResultOutput) BareMetal() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.BareMetal }).(pulumi.BoolOutput)
}

// `true` if the instance type is a burstable performance instance type.
func (o GetInstanceTypeResultOutput) BurstablePerformanceSupported() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.BurstablePerformanceSupported }).(pulumi.BoolOutput)
}

// `true`  if the instance type is a current generation.
func (o GetInstanceTypeResultOutput) CurrentGeneration() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.CurrentGeneration }).(pulumi.BoolOutput)
}

// `true` if Dedicated Hosts are supported on the instance type.
func (o GetInstanceTypeResultOutput) DedicatedHostsSupported() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.DedicatedHostsSupported }).(pulumi.BoolOutput)
}

// Default number of cores for the instance type.
func (o GetInstanceTypeResultOutput) DefaultCores() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.DefaultCores }).(pulumi.IntOutput)
}

// The  default  number of threads per core for the instance type.
func (o GetInstanceTypeResultOutput) DefaultThreadsPerCore() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.DefaultThreadsPerCore }).(pulumi.IntOutput)
}

// Default number of vCPUs for the instance type.
func (o GetInstanceTypeResultOutput) DefaultVcpus() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.DefaultVcpus }).(pulumi.IntOutput)
}

// Indicates whether Amazon EBS encryption is supported.
func (o GetInstanceTypeResultOutput) EbsEncryptionSupport() pulumi.StringOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) string { return v.EbsEncryptionSupport }).(pulumi.StringOutput)
}

// Whether non-volatile memory express (NVMe) is supported.
func (o GetInstanceTypeResultOutput) EbsNvmeSupport() pulumi.StringOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) string { return v.EbsNvmeSupport }).(pulumi.StringOutput)
}

// Indicates that the instance type is Amazon EBS-optimized.
func (o GetInstanceTypeResultOutput) EbsOptimizedSupport() pulumi.StringOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) string { return v.EbsOptimizedSupport }).(pulumi.StringOutput)
}

// The baseline bandwidth performance for an EBS-optimized instance type, in Mbps.
func (o GetInstanceTypeResultOutput) EbsPerformanceBaselineBandwidth() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.EbsPerformanceBaselineBandwidth }).(pulumi.IntOutput)
}

// The baseline input/output storage operations per seconds for an EBS-optimized instance type.
func (o GetInstanceTypeResultOutput) EbsPerformanceBaselineIops() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.EbsPerformanceBaselineIops }).(pulumi.IntOutput)
}

// The baseline throughput performance for an EBS-optimized instance type, in MBps.
func (o GetInstanceTypeResultOutput) EbsPerformanceBaselineThroughput() pulumi.Float64Output {
	return o.ApplyT(func(v GetInstanceTypeResult) float64 { return v.EbsPerformanceBaselineThroughput }).(pulumi.Float64Output)
}

// The maximum bandwidth performance for an EBS-optimized instance type, in Mbps.
func (o GetInstanceTypeResultOutput) EbsPerformanceMaximumBandwidth() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.EbsPerformanceMaximumBandwidth }).(pulumi.IntOutput)
}

// The maximum input/output storage operations per second for an EBS-optimized instance type.
func (o GetInstanceTypeResultOutput) EbsPerformanceMaximumIops() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.EbsPerformanceMaximumIops }).(pulumi.IntOutput)
}

// The maximum throughput performance for an EBS-optimized instance type, in MBps.
func (o GetInstanceTypeResultOutput) EbsPerformanceMaximumThroughput() pulumi.Float64Output {
	return o.ApplyT(func(v GetInstanceTypeResult) float64 { return v.EbsPerformanceMaximumThroughput }).(pulumi.Float64Output)
}

// Whether Elastic Fabric Adapter (EFA) is supported.
func (o GetInstanceTypeResultOutput) EfaSupported() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.EfaSupported }).(pulumi.BoolOutput)
}

// Whether Elastic Network Adapter (ENA) is supported.
func (o GetInstanceTypeResultOutput) EnaSupport() pulumi.StringOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) string { return v.EnaSupport }).(pulumi.StringOutput)
}

// Indicates whether encryption in-transit between instances is supported.
func (o GetInstanceTypeResultOutput) EncryptionInTransitSupported() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.EncryptionInTransitSupported }).(pulumi.BoolOutput)
}

// Describes the FPGA accelerator settings for the instance type.
// * `fpgas.#.count` - The count of FPGA accelerators for the instance type.
// * `fpgas.#.manufacturer` - The manufacturer of the FPGA accelerator.
// * `fpgas.#.memory_size` - The size (in MiB) for the memory available to the FPGA accelerator.
// * `fpgas.#.name` - The name of the FPGA accelerator.
func (o GetInstanceTypeResultOutput) Fpgas() GetInstanceTypeFpgaArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []GetInstanceTypeFpga { return v.Fpgas }).(GetInstanceTypeFpgaArrayOutput)
}

// `true` if the instance type is eligible for the free tier.
func (o GetInstanceTypeResultOutput) FreeTierEligible() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.FreeTierEligible }).(pulumi.BoolOutput)
}

// Describes the GPU accelerators for the instance type.
// * `gpus.#.count` - The number of GPUs for the instance type.
// * `gpus.#.manufacturer` - The manufacturer of the GPU accelerator.
// * `gpus.#.memory_size` - The size (in MiB) for the memory available to the GPU accelerator.
// * `gpus.#.name` - The name of the GPU accelerator.
func (o GetInstanceTypeResultOutput) Gpuses() GetInstanceTypeGpusArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []GetInstanceTypeGpus { return v.Gpuses }).(GetInstanceTypeGpusArrayOutput)
}

// `true` if On-Demand hibernation is supported.
func (o GetInstanceTypeResultOutput) HibernationSupported() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.HibernationSupported }).(pulumi.BoolOutput)
}

// Hypervisor used for the instance type.
func (o GetInstanceTypeResultOutput) Hypervisor() pulumi.StringOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) string { return v.Hypervisor }).(pulumi.StringOutput)
}

// The provider-assigned unique ID for this managed resource.
func (o GetInstanceTypeResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) string { return v.Id }).(pulumi.StringOutput)
}

// Describes the Inference accelerators for the instance type.
// * `inference_accelerators.#.count` - The number of Inference accelerators for the instance type.
// * `inference_accelerators.#.manufacturer` - The manufacturer of the Inference accelerator.
// * `inference_accelerators.#.name` - The name of the Inference accelerator.
func (o GetInstanceTypeResultOutput) InferenceAccelerators() GetInstanceTypeInferenceAcceleratorArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []GetInstanceTypeInferenceAccelerator { return v.InferenceAccelerators }).(GetInstanceTypeInferenceAcceleratorArrayOutput)
}

// Describes the disks for the instance type.
// * `instance_disks.#.count` - The number of disks with this configuration.
// * `instance_disks.#.size` - The size of the disk in GB.
// * `instance_disks.#.type` - The type of disk.
func (o GetInstanceTypeResultOutput) InstanceDisks() GetInstanceTypeInstanceDiskArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []GetInstanceTypeInstanceDisk { return v.InstanceDisks }).(GetInstanceTypeInstanceDiskArrayOutput)
}

// `true` if instance storage is supported.
func (o GetInstanceTypeResultOutput) InstanceStorageSupported() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.InstanceStorageSupported }).(pulumi.BoolOutput)
}

func (o GetInstanceTypeResultOutput) InstanceType() pulumi.StringOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) string { return v.InstanceType }).(pulumi.StringOutput)
}

// `true` if IPv6 is supported.
func (o GetInstanceTypeResultOutput) Ipv6Supported() pulumi.BoolOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) bool { return v.Ipv6Supported }).(pulumi.BoolOutput)
}

// The maximum number of IPv4 addresses per network interface.
func (o GetInstanceTypeResultOutput) MaximumIpv4AddressesPerInterface() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.MaximumIpv4AddressesPerInterface }).(pulumi.IntOutput)
}

// The maximum number of IPv6 addresses per network interface.
func (o GetInstanceTypeResultOutput) MaximumIpv6AddressesPerInterface() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.MaximumIpv6AddressesPerInterface }).(pulumi.IntOutput)
}

// The maximum number of physical network cards that can be allocated to the instance.
func (o GetInstanceTypeResultOutput) MaximumNetworkCards() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.MaximumNetworkCards }).(pulumi.IntOutput)
}

// The maximum number of network interfaces for the instance type.
func (o GetInstanceTypeResultOutput) MaximumNetworkInterfaces() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.MaximumNetworkInterfaces }).(pulumi.IntOutput)
}

// Size of the instance memory, in MiB.
func (o GetInstanceTypeResultOutput) MemorySize() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.MemorySize }).(pulumi.IntOutput)
}

// Describes the network performance.
func (o GetInstanceTypeResultOutput) NetworkPerformance() pulumi.StringOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) string { return v.NetworkPerformance }).(pulumi.StringOutput)
}

// A list of architectures supported by the instance type.
func (o GetInstanceTypeResultOutput) SupportedArchitectures() pulumi.StringArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []string { return v.SupportedArchitectures }).(pulumi.StringArrayOutput)
}

// A list of supported placement groups types.
func (o GetInstanceTypeResultOutput) SupportedPlacementStrategies() pulumi.StringArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []string { return v.SupportedPlacementStrategies }).(pulumi.StringArrayOutput)
}

// Indicates the supported root device types.
func (o GetInstanceTypeResultOutput) SupportedRootDeviceTypes() pulumi.StringArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []string { return v.SupportedRootDeviceTypes }).(pulumi.StringArrayOutput)
}

// Indicates whether the instance type is offered for spot or On-Demand.
func (o GetInstanceTypeResultOutput) SupportedUsagesClasses() pulumi.StringArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []string { return v.SupportedUsagesClasses }).(pulumi.StringArrayOutput)
}

// The supported virtualization types.
func (o GetInstanceTypeResultOutput) SupportedVirtualizationTypes() pulumi.StringArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []string { return v.SupportedVirtualizationTypes }).(pulumi.StringArrayOutput)
}

// The speed of the processor, in GHz.
func (o GetInstanceTypeResultOutput) SustainedClockSpeed() pulumi.Float64Output {
	return o.ApplyT(func(v GetInstanceTypeResult) float64 { return v.SustainedClockSpeed }).(pulumi.Float64Output)
}

// Total memory of all FPGA accelerators for the instance type (in MiB).
func (o GetInstanceTypeResultOutput) TotalFpgaMemory() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.TotalFpgaMemory }).(pulumi.IntOutput)
}

// Total size of the memory for the GPU accelerators for the instance type (in MiB).
func (o GetInstanceTypeResultOutput) TotalGpuMemory() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.TotalGpuMemory }).(pulumi.IntOutput)
}

// The total size of the instance disks, in GB.
func (o GetInstanceTypeResultOutput) TotalInstanceStorage() pulumi.IntOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) int { return v.TotalInstanceStorage }).(pulumi.IntOutput)
}

// List of the valid number of cores that can be configured for the instance type.
func (o GetInstanceTypeResultOutput) ValidCores() pulumi.IntArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []int { return v.ValidCores }).(pulumi.IntArrayOutput)
}

// List of the valid number of threads per core that can be configured for the instance type.
func (o GetInstanceTypeResultOutput) ValidThreadsPerCores() pulumi.IntArrayOutput {
	return o.ApplyT(func(v GetInstanceTypeResult) []int { return v.ValidThreadsPerCores }).(pulumi.IntArrayOutput)
}

func init() {
	pulumi.RegisterOutputType(GetInstanceTypeResultOutput{})
}