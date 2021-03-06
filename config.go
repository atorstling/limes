package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// Config hold configuration read from the configuration file
type Config struct {
	profiles Profiles
}

// NewConfig returns a new Config struct
func NewConfig() *Config {
	config := &Config{
		profiles: make(Profiles),
	}

	return config
}

// Profiles is a map for AWS profiles
type Profiles map[string]Profile

// Profile defines an AWS IAM profile
type Profile struct {
	AwsAccessKeyID     string `yaml:"aws_access_key_id"`
	AwsSecretAccessKey string `yaml:"aws_secret_access_key"`
	AwsSessionToken    string
	Region             string `yaml:"region"`
	MFASerial          string `yaml:"mfa_serial"`
	RoleARN            string `yaml:"role_arn"`
	SourceProfile      string `yaml:"source_profile"`
	RoleSessionName    string `yaml:"role_session_name"`
}

const (
	awsConfDir         = ".aws"
	awsConfigFile      = "config"
	awsCredentialsFile = "credentials"
	awsAccessKeyEnv    = "AWS_ACCESS_KEY_ID"
	awsSecretKeyEnv    = "AWS_SECRET_ACCESS_KEY"
)

var (
	errActiveAWSEnvironment     = fmt.Errorf("active AWS environment variables")
	errActiveAWSCredentialsFile = fmt.Errorf("active AWS credentials file")
	errActiveAWSConfigFile      = fmt.Errorf("active AWS config file")
)

func checkActiveAWSConfig() (err error) {
	var active bool

	if active, err = doCheck(activeAWSEnvironment, err); active {
		return errActiveAWSEnvironment
	}

	if active, err = doCheck(activeAWSCredentialsFile, err); active {
		return errActiveAWSCredentialsFile
	}

	if active, err = doCheck(activeAWSConfigFile, err); active {
		return errActiveAWSConfigFile
	}

	return fmt.Errorf("checkActiveAWSConfig: %v", err)
}

func doCheck(fn func() (bool, error), err error) (bool, error) {
	if err != nil {
		return false, err
	}

	return fn()
}

func activeAWSConfigFile() (active bool, err error) {
	home, err := homeDir()
	if err != nil {
		return false, fmt.Errorf("unable to fetch user information: %v", err)
	}

	if _, err := os.Stat(filepath.Join(home, awsConfDir, awsConfigFile)); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

func activeAWSCredentialsFile() (active bool, err error) {
	home, err := homeDir()
	if err != nil {
		return false, fmt.Errorf("unable to fetch user information: %v", err)
	}

	if _, err := os.Stat(filepath.Join(home, awsConfDir, awsCredentialsFile)); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

func activeAWSEnvironment() (active bool, err error) {
	if env := os.Getenv(awsAccessKeyEnv); env != "" {
		return true, nil
	}

	if env := os.Getenv(awsSecretKeyEnv); env != "" {
		return true, nil
	}

	return false, nil

}

func homeDir() (homeDir string, err error) {
	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir, nil
	}

	// fallback to environemt variables
	var home string
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home != "" {
			return home, nil
		}

		home = os.Getenv("USERPROFILE")
		if home != "" {
			return home, nil
		}

		return "", fmt.Errorf("fallback failed, set `USERPROFILE` environment variable")
	}

	home = os.Getenv("HOME")
	if home != "" {
		return home, nil
	}

	return "", fmt.Errorf("fallback failed, set `HOME` environment variable")
}
