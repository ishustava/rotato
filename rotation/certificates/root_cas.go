package certificates

import (
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	credhubClient "code.cloudfoundry.org/credhub-cli/credhub"
	"crypto/x509"
	"encoding/pem"
)

func FindRootCAs(findResults credentials.FindResults, ch *credhubClient.CredHub) ([]credentials.Certificate, error) {
	var rootCAs []credentials.Certificate

	for _, findResult := range findResults.Credentials {
		credential, err := ch.GetLatestCertificate(findResult.Name)
		if err == nil && credential.Type == "certificate" {
			if isRootCA(credential.Value) {
				rootCAs = append(rootCAs, credential)
			}
		}
	}

	return rootCAs, nil
}

func isRootCA(certificate values.Certificate) bool {
	parsedCertificate := parsePemCertificate(certificate.Certificate)

	if parsedCertificate.IsCA && parsedCertificate.Issuer.String() == parsedCertificate.Subject.String() {
		return true
	}

	return false
}

func parsePemCertificate(pemCert string) (*x509.Certificate) {
	block, _ := pem.Decode([]byte(pemCert))
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("could not parse pem certificate")
	}
	return parsedCert
}
