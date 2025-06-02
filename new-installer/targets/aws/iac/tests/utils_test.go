// SPDX-FileCopyrightText: 2025 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package aws_iac_test

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type AWSS3BackendConfig struct {
	Region string `json:"region"`
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

func CreateVPC(region string) (string, string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return "", "", fmt.Errorf("unable to load AWS SDK config: %w", err)
	}
	ec2Client := ec2.NewFromConfig(cfg)
	createVpcInput := &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.250.0.0/16"),
	}
	vpcOutput, err := ec2Client.CreateVpc(context.Background(), createVpcInput)
	if err != nil {
		return "", "", fmt.Errorf("failed to create VPC: %w", err)
	}
	vpcID := *vpcOutput.Vpc.VpcId
	log.Printf("Created VPC with ID: %s", vpcID)
	createSubnetInput := &ec2.CreateSubnetInput{
		CidrBlock: aws.String("10.250.0.0/24"),
		VpcId:     aws.String(vpcID),
	}
	subnetOutput, err := ec2Client.CreateSubnet(context.Background(), createSubnetInput)
	if err != nil {
		return "", "", fmt.Errorf("failed to create subnet: %w", err)
	}
	subnetID := *subnetOutput.Subnet.SubnetId
	log.Printf("Created Subnet with ID: %s", subnetID)
	return vpcID, subnetID, nil
}

func DeleteVPC(region string, vpcID string) error {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %w", err)
	}
	ec2Client := ec2.NewFromConfig(cfg)
	deleteVpcInput := &ec2.DeleteVpcInput{
		VpcId: aws.String(vpcID),
	}
	_, err = ec2Client.DeleteVpc(context.Background(), deleteVpcInput)
	if err != nil {
		return fmt.Errorf("failed to delete VPC: %w", err)
	}
	log.Printf("Deleted VPC with ID: %s", vpcID)
	return nil
}
