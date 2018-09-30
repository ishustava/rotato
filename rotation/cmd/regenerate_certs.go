package cmd

import (
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	credhubClient "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/commands"
	"github.com/ishustava/rotato/rotation/certificates"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"log"
)

type RegenerateCertsCommand struct {
	CredHubURL          string   `short:"u" long:"credhub-server" description:"URL of the CredHub server, e.g. https://example.com:8844" env:"CREDHUB_SERVER" required:"true"`
	CredHubClient       string   `short:"c" long:"credhub-client" description:"UAA client for the CredHub Server" env:"CREDHUB_CLIENT" required:"true"`
	CredHubClientSecret string   `short:"s" long:"credhub-secret" description:"UAA client secret for the CredHub Server" env:"CREDHUB_SECRET" required:"true"`
	DeploymentName      string   `short:"d" long:"deployment-name" description:"Name of the BOSH deployment that needs certificate rotation" env:"BOSH_DEPLOYMENT" required:"true"`
	CaCerts             []string `long:"ca-cert" description:"Trusted CA for API and UAA TLS connections. Multiple flags may be provided." env:"CREDHUB_CA_CERT" required:"true"`
}

func (cmd RegenerateCertsCommand) Execute([]string) error {
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

	err = regenerateCerts(rootCAs, ch)
	if err != nil {
		return err
	}

	return nil
}

func regenerateCerts(rootCAs []credentials.Certificate, credHub *credhubClient.CredHub) error {
	for _, rootCA := range rootCAs {
		log.Printf("Regenerating certificates for %s\n", rootCA.Name)

		resuts, err := credHub.BulkRegenerate(rootCA.Name)
		if err != nil {
			return err
		}

		for _, regenerateResult := range resuts.Certificates {
			log.Printf("Regenerated %s\n", regenerateResult)
		}
	}

	return nil
}
