// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package hcloud

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Provides details about a specific Hetzner Cloud Certificate.
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
	Id *int `pulumi:"id"`
	// (map) User-defined labels (key-value pairs) assigned to the certificate.
	Labels map[string]interface{} `pulumi:"labels"`
	// (string) Name of the Certificate.
	Name *string `pulumi:"name"`
	// (string) Point in time when the Certificate stops being valid (in ISO-8601 format).
	NotValidAfter string `pulumi:"notValidAfter"`
	// (string) Point in time when the Certificate becomes valid (in ISO-8601 format).
	NotValidBefore string  `pulumi:"notValidBefore"`
	WithSelector   *string `pulumi:"withSelector"`
}