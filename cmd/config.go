package cmd

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type SsmParameter struct {
	ARN              string `json:"arn"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	LastModifiedDate int64  `json:"last_modified_date"`
	Value            string `json:"value"`
	Version          int32  `json:"version"`
}

type SsmParameterGroups struct {
	String       []SsmParameter `json:"string"`
	SecureString []SsmParameter `json:"securestring"`
	StringList   []SsmParameter `json:"stringlist"`
}

func InitAwsConfig(useLocalStack bool, parametersRegion string) aws.Config {
	var awsConfig aws.Config
	if useLocalStack {
		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:4566",
				SigningRegion: "us-east-1",
			}, nil
		})
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolver(customResolver),
		)
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
		awsConfig = cfg
	} else {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(parametersRegion),
		)
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
		awsConfig = cfg
	}
	return awsConfig
}
