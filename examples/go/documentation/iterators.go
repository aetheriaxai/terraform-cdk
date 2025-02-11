// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/hashicorp/aws/acmcertificate"
	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/hashicorp/aws/acmcertificatevalidation"
	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/hashicorp/aws/dataawsroute53zone"
	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/hashicorp/aws/instance"
	aws "github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/hashicorp/aws/provider"
	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/hashicorp/aws/route53record"

	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/hashicorp/aws/s3bucket"
	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/hashicorp/aws/s3bucketobject"

	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/integrations/github/datagithuborganization"
	github "github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/integrations/github/provider"
	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/integrations/github/team"
	"github.com/hashicorp/terraform-cdk/examples/go/documentation/generated/integrations/github/teammembers"
)

func NewIteratorsStack(scope constructs.Construct, name string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &name)

	aws.NewAwsProvider(stack, jsii.String("aws"), &aws.AwsProviderConfig{
		Region: jsii.String("us-west-2"),
	})

	// DOCS_BLOCK_START:iterators
	list := cdktf.NewTerraformVariable(stack, jsii.String("list"), &cdktf.TerraformVariableConfig{
		Type: cdktf.VariableType_LIST_STRING(),
	})

	simpleIterator := cdktf.TerraformIterator_FromList(list.ListValue())

	s3bucket.NewS3Bucket(stack, jsii.String("iterator-bucket"), &s3bucket.S3BucketConfig{
		ForEach: simpleIterator,
		Bucket:  cdktf.Token_AsString(simpleIterator.Value(), nil),
	})
	// DOCS_BLOCK_END:iterators

	// DOCS_BLOCK_START:iterators-complex-types
	complexList := cdktf.NewTerraformLocal(stack, jsii.String("complex-list-local"), []map[string]interface{}{
		{
			"name": "website-static-files",
			"tags": map[string]string{"app": "website"},
		},
		{
			"name": "images",
			"tags": map[string]string{"app": "image-converter"},
		},
	})

	complexIterator := cdktf.TerraformIterator_FromList(complexList.Expression())

	s3bucket.NewS3Bucket(stack, jsii.String("complex-iterator-bucket"), &s3bucket.S3BucketConfig{
		ForEach: complexIterator,
		Bucket:  complexIterator.GetString(jsii.String("name")),
		Tags:    complexIterator.GetStringMap(jsii.String("tags")),
	})
	// DOCS_BLOCK_END:iterators-complex-types

	// DOCS_BLOCK_START:iterators-list-attributes
	orgName := "my-org"

	github.NewGithubProvider(stack, jsii.String("github"), &github.GithubProviderConfig{
		Organization: &orgName,
	})

	team := team.NewTeam(stack, jsii.String("core-team"), &team.TeamConfig{
		Name: jsii.String("core"),
	})

	orgMembers := datagithuborganization.NewDataGithubOrganization(stack, jsii.String("org"), &datagithuborganization.DataGithubOrganizationConfig{
		Name: &orgName,
	})

	orgMemberIterator := cdktf.TerraformIterator_FromList(orgMembers.Members())

	teammembers.NewTeamMembers(stack, jsii.String("members"), &teammembers.TeamMembersConfig{
		TeamId: team.Id(),
		Members: orgMemberIterator.Dynamic(&map[string]interface{}{
			"username": orgMemberIterator.Value(),
			"role":     jsii.String("maintainer"),
		}),
	})
	// DOCS_BLOCK_END:iterators-list-attributes

	// DOCS_BLOCK_START:iterators-count
	servers := cdktf.NewTerraformVariable(stack, jsii.String("servers"), &cdktf.TerraformVariableConfig{
		Type: cdktf.VariableType_NUMBER(),
	})
	count := cdktf.TerraformCount_Of(servers.NumberValue())
	instance.NewInstance(stack, jsii.String("server"), &instance.InstanceConfig{
		Count:        count,
		Ami:          jsii.String("ami-a1b2c3d4"),
		InstanceType: jsii.String("t2.micro"),
		Tags: &map[string]*string{
			"Name": jsii.String("Server ${" + *cdktf.Token_AsString(count.Index(), nil) + "}"),
		},
	})
	// DOCS_BLOCK_END:iterators-count

	// DOCS_BLOCK_START:iterators-complex-lists
	cert := acmcertificate.NewAcmCertificate(stack, jsii.String("cert"), &acmcertificate.AcmCertificateConfig{
		DomainName:       jsii.String("example.com"),
		ValidationMethod: jsii.String("DNS"),
	})

	dataAwsRoute53ZoneExample := dataawsroute53zone.NewDataAwsRoute53Zone(stack, jsii.String("dns_zone"), &dataawsroute53zone.DataAwsRoute53ZoneConfig{
		Name:        jsii.String("example.com"),
		PrivateZone: jsii.Bool(false),
	})

	exampleForEachIterator := cdktf.TerraformIterator_FromComplexList(cert.DomainValidationOptions(), jsii.String("domain_name"))

	records := route53record.NewRoute53Record(stack, jsii.String("record"), &route53record.Route53RecordConfig{
		ForEach:        exampleForEachIterator,
		AllowOverwrite: jsii.Bool(true),
		Name:           exampleForEachIterator.GetString(jsii.String("name")),
		Records:        jsii.Strings(*exampleForEachIterator.GetString(jsii.String("record"))),
		Ttl:            jsii.Number(60),
		Type:           exampleForEachIterator.GetString(jsii.String("type")),
		ZoneId:         dataAwsRoute53ZoneExample.ZoneId(),
	})

	recordsIterator := cdktf.TerraformIterator_FromResources(records)

	acmcertificatevalidation.NewAcmCertificateValidation(stack, jsii.String("validation"), &acmcertificatevalidation.AcmCertificateValidationConfig{
		CertificateArn:        cert.Arn(),
		ValidationRecordFqdns: cdktf.Token_AsList(recordsIterator.PluckProperty(jsii.String("fqdn")), nil),
	})
	// DOCS_BLOCK_END:iterators-complex-lists

	// DOCS_BLOCK_START:iterators-chain
	config := cdktf.NewTerraformLocal(stack, jsii.String("config-local"), []map[string]interface{}{
		{
			"name": "website-static-files",
			"tags": map[string]string{"app": "website"},
		},
		{
			"name": "images",
			"tags": map[string]string{"app": "image-converter"},
		},
	})

	s3BucketConfigurationIterator := cdktf.TerraformIterator_FromList(config.Expression())
	s3Buckets := s3bucket.NewS3Bucket(stack, jsii.String("complex-iterator-buckets"), &s3bucket.S3BucketConfig{
		ForEach: s3BucketConfigurationIterator,
		Bucket:  s3BucketConfigurationIterator.GetString(jsii.String("name")),
		Tags:    s3BucketConfigurationIterator.GetStringMap(jsii.String("tags")),
	})

	s3BucketsIterator := cdktf.TerraformIterator_FromResources(s3Buckets)
	helpFile := cdktf.NewTerraformAsset(stack, jsii.String("help"), &cdktf.TerraformAssetConfig{
		Path: jsii.String("./help"),
	})
	s3bucketobject.NewS3BucketObject(stack, jsii.String("object"), &s3bucketobject.S3BucketObjectConfig{
		ForEach: s3BucketsIterator,
		Bucket:  s3BucketsIterator.GetString(jsii.String("id")),
		Key:     jsii.String("help"),
		Source:  helpFile.Path(),
	})
	// DOCS_BLOCK_END:iterators-chain

	// DOCS_BLOCK_START:iterators-for-expression
	values := cdktf.NewTerraformLocal(stack, jsii.String("values"), []map[string]interface{}{
		{
			"name": "website-static-files",
			"tags": map[string]string{"app": "website"},
		},
		{
			"name": "images",
			"tags": map[string]string{"app": "image-converter"},
		},
	})

	mapIterator := cdktf.TerraformIterator_FromList(values.Expression())
	cdktf.NewTerraformLocal(stack, jsii.String("list-of-keys"), mapIterator.Keys())
	cdktf.NewTerraformLocal(stack, jsii.String("list-of-values"), mapIterator.Values())
	cdktf.NewTerraformLocal(stack, jsii.String("list-of-names"), mapIterator.PluckProperty(jsii.String("name")))
	cdktf.NewTerraformLocal(stack, jsii.String("list-of-names-of-included"), mapIterator.ForExpressionForList(jsii.String("val.name if val.included")))
	cdktf.NewTerraformLocal(stack, jsii.String("map-with-names-as-key-and-tags-as-value-of-included"), mapIterator.ForExpressionForMap(jsii.String("val.name"), jsii.String("val.tags if val.included")))
	// DOCS_BLOCK_END:iterators-for-expression

	return stack
}
