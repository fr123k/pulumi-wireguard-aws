// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package hcloud

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Provides details about a specific Hetzner Cloud Certificate.
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		opt0 := "sample-certificate-1"
// 		_, err := hcloud.LookupCertificate(ctx, &hcloud.LookupCertificateArgs{
// 			Name: &opt0,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		opt1 := 4711
// 		_, err = hcloud.LookupCertificate(ctx, &hcloud.LookupCertificateArgs{
// 			Id: &opt1,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func LookupCertificate(ctx *pulumi.Context, args *LookupCertificateArgs, opts ...pulumi.InvokeOption) (*LookupCertificateResult, error) {
	var rv LookupCertificateResult
	err := ctx.Invoke("hcloud:index/getCertificate:getCertificate", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getCertificate.
type LookupCertificateArgs struct {
	// ID of the certificate.
	Id *int `pulumi:"id"`
	// Name of the certificate.
	Name *string `pulumi:"name"`
	// [Label selector](https://docs.hetzner.cloud/#overview-label-selector)
	WithSelector *string `pulumi:"withSelector"`
}

// A collection of values returned by getCertificate.
type LookupCertificateResult struct {
	// (string) PEM encoded TLS certificate.
	Certificate string `pulumi:"certificate"`
	// (string) Point in time when the Certificate was created at Hetzner Cloud (in ISO-8601 format).
	Created string `pulumi:"created"`
	// (list) Domains and subdomains covered by the certificate.
	DomainNames []string `pulumi:"domainNames"`
	// (string) Fingerprint of the certificate.
	Fingerprint string `pulumi:"fingerprint"`
	// (int) Unique ID of the certificate.
	Id int `pulumi:"id"`
	// (map) User-defined labels (key-value pairs) assigned to the certificate.
	Labels map[string]interface{} `pulumi:"labels"`
	// (string) Name of the Certificate.
	Name *string `pulumi:"name"`
	// (string) Point in time when the Certificate stops being valid (in ISO-8601 format).
	NotValidAfter string `pulumi:"notValidAfter"`
	// (string) Point in time when the Certificate becomes valid (in ISO-8601 format).
	NotValidBefore string  `pulumi:"notValidBefore"`
	Type           string  `pulumi:"type"`
	WithSelector   *string `pulumi:"withSelector"`
}
