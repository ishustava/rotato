package cmd

type RotatorCommand struct {
	AddNewCAs AddNewCAsCommand `command:"add-new-cas" alias:"a" description:"Regenerate new CAs and make them available to the next BOSH deploy"`
	RegenerateCerts RegenerateCertsCommand `command:"regenerate-certs" alias:"r" description:"Regenerate new certificates"`
}

var Rotator RotatorCommand