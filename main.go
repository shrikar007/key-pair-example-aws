package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)
func init() {
	viper.SetConfigType("toml")
	viper.SetConfigName("config") // name of config file (without extension)

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	viper.AddConfigPath(path)
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Fatal(err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.ReadInConfig() // Find and read the config file
		if err != nil {
			log.Fatal(err)
		}
	})

	viper.WatchConfig()
}
func main(){
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials:credentials.NewStaticCredentials(viper.GetString("cred.accesskeyid"),viper.GetString("cred.secretaccesskey"),"")},
	)
	svc := ec2.New(sess)
	key:="KEY_NAME"
	result, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
				KeyName: aws.String(key),
			})
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
					exitErrorf("Keypair %q already exists.", key)
				}
				exitErrorf("Unable to create key pair: %s, %v.", key, err)
			}
			fmt.Printf("Created key pair %q %s\n%s\n",
				*result.KeyName, *result.KeyFingerprint,
				*result.KeyMaterial)

			path,err:=os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			file,err:=os.Create(path+"/"+key+".pem")
			if err != nil {
				log.Fatal(err)
			}
			file.WriteString(*result.KeyMaterial)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
