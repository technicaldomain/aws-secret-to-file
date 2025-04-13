package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "aws-secrets-to-file",
	Short: "A CLI tool to retrieve secrets from AWS Secrets Manager and write it to file",
	Run: func(cmd *cobra.Command, args []string) {
		retrieveSecret()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringSliceP("secret", "s", []string{}, "The ID or ARN of the secrets in Secrets Manager (multiple secrets supported)")
	rootCmd.PersistentFlags().StringSliceP("output", "o", []string{}, "The paths to the output files (multiple files supported)")
	rootCmd.PersistentFlags().Bool("binary", false, "Get binary secret instead of text")
	viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("binary", rootCmd.PersistentFlags().Lookup("binary"))
}

func initConfig() {
	viper.AutomaticEnv()
}

func retrieveSecret() {
	secretIDs := viper.GetStringSlice("secret")
	outputFilePaths := viper.GetStringSlice("output")
	binary := viper.GetBool("binary")

	if err := validateInputs(secretIDs, outputFilePaths); err != nil {
		log.Fatalf("Input validation error: %v", err)
	}

	var awsCfg aws.Config
	var err error
	awsCfg, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading AWS configuration: %v", err)
	}

	secretsManagerClient := secretsmanager.NewFromConfig(awsCfg)

	for i, secretID := range secretIDs {
		if err := processSecret(secretsManagerClient, secretID, outputFilePaths[i], binary); err != nil {
			log.Fatalf("Error processing secret %s: %v", secretID, err)
		}
	}
}

func validateInputs(secretIDs, outputFilePaths []string) error {
	if len(secretIDs) == 0 || len(outputFilePaths) == 0 {
		return fmt.Errorf("both --secret and --output flags must be provided")
	}
	if len(secretIDs) != len(outputFilePaths) {
		return fmt.Errorf("the number of secrets and output files must match")
	}
	return nil
}

func processSecret(client *secretsmanager.Client, secretID, outputFilePath string, binary bool) error {
	secretValueInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	secretValueOutput, err := client.GetSecretValue(context.TODO(), secretValueInput)
	if err != nil {
		return fmt.Errorf("error retrieving secret: %w", err)
	}

	if binary {
		return writeBinarySecret(secretValueOutput, secretID, outputFilePath)
	}
	return writeStringSecret(secretValueOutput, secretID, outputFilePath)
}

func writeBinarySecret(secretValueOutput *secretsmanager.GetSecretValueOutput, secretID, outputFilePath string) error {
	if len(secretValueOutput.SecretBinary) > 0 {
		if err := os.WriteFile(outputFilePath, secretValueOutput.SecretBinary, 0644); err != nil {
			return fmt.Errorf("error writing binary secret to file: %w", err)
		}
		fmt.Printf("Secret %s has been written to %s\n", secretID, outputFilePath)
		return nil
	}
	return fmt.Errorf("secret %s has no binary data", secretID)
}

func writeStringSecret(secretValueOutput *secretsmanager.GetSecretValueOutput, secretID, outputFilePath string) error {
	secretString := aws.ToString(secretValueOutput.SecretString)
	if len(secretString) > 0 {
		if err := os.WriteFile(outputFilePath, []byte(secretString), 0644); err != nil {
			return fmt.Errorf("error writing string secret to file: %w", err)
		}
		fmt.Printf("Secret %s has been written to %s\n", secretID, outputFilePath)
		return nil
	}
	return fmt.Errorf("secret %s has no string data", secretID)
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		retrieveSecret()
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
