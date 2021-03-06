package cmd

import (
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	credhubClient "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/commands"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	"github.com/ishustava/rotato/rotation/certificates"
	"fmt"
	"encoding/json"
	"strings"
	"log"
)

type AddNewCAsCommand struct {
	CredHubURL          string   `short:"u" long:"credhub-server" description:"URL of the CredHub server, e.g. https://example.com:8844" env:"CREDHUB_SERVER" required:"true"`
	CredHubClient       string   `short:"c" long:"credhub-client" description:"UAA client for the CredHub Server" env:"CREDHUB_CLIENT" required:"true"`
	CredHubClientSecret string   `short:"s" long:"credhub-secret" description:"UAA client secret for the CredHub Server" env:"CREDHUB_SECRET" required:"true"`
	DeploymentName      string   `short:"d" long:"deployment-name" description:"Name of the BOSH deployment that needs certificate rotation" env:"BOSH_DEPLOYMENT" required:"true"`
	CaCerts             []string `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections. Multiple flags may be provided." env:"CREDHUB_CA_CERT" required:"true"`
}

func (cmd AddNewCAsCommand) Execute([]string) error {
	caCerts, err := commands.ReadOrGetCaCerts(cmd.CaCerts)
	if err != nil {
		return err
	}

	ch, err := credhubClient.New(
		cmd.CredHubURL,
		credhubClient.CaCerts(caCerts...),
		credhubClient.Auth(
			auth.UaaClientCredentials(cmd.CredHubClient, cmd.CredHubClientSecret)),
	)
	if err != nil {
		return err
	}

	deploymentCredentials, err := ch.FindByPartialName(cmd.DeploymentName)
	if err != nil {
		return err
	}

	rootCAs, err := certificates.FindRootCAs(deploymentCredentials, ch)
	if err != nil {
		return err
	}

	rootCAsToRotate := getCAsToRotate(rootCAs)

	err = regenerateAndSetNewCAs(rootCAsToRotate, ch)
	if err != nil {
		return err
	}

	return nil
}

func regenerateAndSetNewCAs(rootCAs []credentials.Certificate, credHub *credhubClient.CredHub) error {
	for _, rootCA := range rootCAs {
		newCredential, err := credHub.Regenerate(rootCA.Name)
		if err != nil {
			return err
		}

		caJson, _ := json.Marshal(newCredential.Value)
		var newCAValue values.Certificate
		json.Unmarshal(caJson, &newCAValue)

		newCAValue.Certificate = fmt.Sprintf("%s%s", newCAValue.Certificate, rootCA.Value.Certificate)

		_, err = credHub.SetCertificate(rootCA.Name, newCAValue)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCAsToRotate(rootCAs []credentials.Certificate) []credentials.Certificate {
	var rootCAsToRotate []credentials.Certificate

	for _, ca := range rootCAs {
		if strings.Count(ca.Value.Certificate, "-----BEGIN CERTIFICATE-----") > 1 {
			log.Printf("%s is ready for step 1 of CA rotation\n", ca.Name)
		} else {
			log.Printf("Found root CA to rotate: %s\n", ca.Name)
			rootCAsToRotate = append(rootCAsToRotate, ca)
		}
	}

	log.Printf("Will regenerate %d root CAs\n", len(rootCAsToRotate))

	return rootCAsToRotate
}
