package main

import (
	"os"

	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	config := vault.DefaultConfig()
	config.Address = "http://vault.default:8200"

	client, err := vault.NewClient(config)
	if err != nil {
		log.WithError(err).Panicln("unable to initialize Vault client")
	}
	log.Infof("connected to vault service running on %s", config.Address)

	initialised, err := client.Sys().InitStatus()
	if err != nil {
		log.WithError(err).Panicln("error checking vault initialization status")
	}

	if initialised {
		log.Infoln("vault already initialized, exiting")
		os.Exit(0)
	}

	log.Infoln("initializing vault")
	initResponse, err := client.Sys().Init(&vault.InitRequest{
		SecretShares:    5,
		SecretThreshold: 3,
	})
	if err != nil {
		log.WithError(err).Panicln("error initializing vault")
	}

	// FIXME: this shouldn't be logged in the future these should be put in a secret or somewhere safe to be retrieved and stored securely
	log.WithFields(logrus.Fields{
		"keys":          initResponse.Keys,
		"recovery_keys": initResponse.RecoveryKeys,
		"root_token":    initResponse.RootToken,
	}).Infoln("vault initialized successfully")
	client.SetToken(initResponse.RootToken)

	log.Infoln("unsealing vault")
	for i, key := range initResponse.Keys {
		sealStatus, err := client.Sys().Unseal(key)
		if err != nil {
			log.WithError(err).Panicf("error unsealing vault with key [%d]", i+1)
		}

		if !sealStatus.Sealed {
			log.Infoln("vault unsealed successfully")
			break
		}
	}

	// FIXME: this is failing due to no handlers for the kv-v2/data/db route need to enable the secret engine?
	_, err = client.Logical().Write("kv-v2/data/db", map[string]interface{}{
		"username": "postgres",
		"password": "password",
	})
	if err != nil {
		log.WithError(err).Panicln("error writing db secret to vault")
	}

	secret, err := client.Logical().Read("kv-v2/data/db")
	if err != nil {
		log.WithError(err).Panicln("error reading db secrete from vault")
	}
	log.Infof("retrieved secret data: %+v", secret.Data["data"].(map[string]interface{}))
}
