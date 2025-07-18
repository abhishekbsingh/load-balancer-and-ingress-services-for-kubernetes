// Copyright © 2025 Broadcom Inc. and/or its subsidiaries. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0
package models

// This file is auto-generated.

// IPAMDNSAwsProfile ipam Dns aws profile
// swagger:model IpamDnsAwsProfile
type IPAMDNSAwsProfile struct {

	// AWS access key ID. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	AccessKeyID *string `json:"access_key_id,omitempty"`

	// A list of subnets used for source IP allocation for egress services in Openshift/k8s on Aws. Field introduced in 18.2.3. Maximum of 1 items allowed. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	EgressServiceSubnets []string `json:"egress_service_subnets,omitempty"`

	// IAM assume role for cross-account access. Field introduced in 17.1.1. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	IamAssumeRole *string `json:"iam_assume_role,omitempty"`

	// If enabled and the virtual service is not floating ip capable, vip will be published to both private and public zones. This flag is applicable only for AWS DNS profile. Field introduced in 17.2.10. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	PublishVipToPublicZone *bool `json:"publish_vip_to_public_zone,omitempty"`

	// AWS region. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	Region *string `json:"region,omitempty"`

	// AWS secret access key. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	SecretAccessKey *string `json:"secret_access_key,omitempty"`

	// Default TTL for all records. Allowed values are 1-172800. Field introduced in 17.1.3. Unit is SEC. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	TTL *uint32 `json:"ttl,omitempty"`

	// Usable domains to pick from Amazon Route 53. Field introduced in 17.1.1. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	UsableDomains []string `json:"usable_domains,omitempty"`

	// Usable networks for Virtual IP. If VirtualService does not specify a network and auto_allocate_ip is set, then the first available network from this list will be chosen for IP allocation. Field introduced in 17.1.1. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	UsableNetworkUuids []string `json:"usable_network_uuids,omitempty"`

	// Use IAM roles instead of access and secret key. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	UseIamRoles *bool `json:"use_iam_roles,omitempty"`

	// VPC name. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	Vpc *string `json:"vpc,omitempty"`

	// VPC ID. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	// Required: true
	VpcID *string `json:"vpc_id"`

	// Network configuration for Virtual IP per AZ. Field introduced in 17.1.3. Allowed in Enterprise edition with any value, Essentials, Basic, Enterprise with Cloud Services edition.
	Zones []*AwsZoneNetwork `json:"zones,omitempty"`
}
