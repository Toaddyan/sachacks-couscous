package main

import (
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/oauth2/google"
	cloudshell "google.golang.org/api/cloudshell/v1"
)

const (
	defaultCloudShellResource = "users/me/environments/default"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// rsa.GenerateKey() =>
// x509.MarshalPKIXPublicKey() =>
// pem.Encode()
func main() {
	// Reads in the credentials.json
	b, err := ioutil.ReadFile("./secrets/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Creates OAuth config
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, cloudshell.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Creates client based off the config
	client := getClient(config)

	// Get cloudShellService
	cloudshellService, err := cloudshell.New(client)

	// Get the Enviroments subservice
	userEnviromentsService := cloudshell.NewUsersEnvironmentsService(cloudshellService)

	// Generate a brand-new SSH key pair, and save them in secrets/rsa-id.pub and secrets/rsa-id.private
	// publicKey := generateSSHKey()

	// Start Cloudshell enviroment
	startRequest := userEnviromentsService.Start(defaultCloudShellResource, &cloudshell.StartEnvironmentRequest{
		// PublicKeys: []string{publicKey},
	})
	_, err = startRequest.Do()
	if err != nil {
		log.Panicf("Could not start cloudshell environment: %v\n", err)
	}

	call := userEnviromentsService.Get("users/me/environments/default")
	var environment *cloudshell.Environment

	for {
		time.Sleep(time.Second)

		log.Println("Checking status of Cloud Shell environment")
		environment, err = call.Do()
		if err != nil {
			log.Println("Could not get environment: %v\n", err)
		}

		if environment.State == "RUNNING" {
			log.Println("Cloud Shell environment started")
			break
		}
	}

	// Generate the ansible-playbook/inventory.yaml
	err = generateInventory(Inventory{
		UserName:   environment.SshUsername,
		Host:       environment.SshHost,
		Port:       int(environment.SshPort),
		PrivateKey: "./secrets/id_rsa",
	})
	if err != nil {
		log.Panicf("Could not generate inventory: %v\n", err)
	}

}
