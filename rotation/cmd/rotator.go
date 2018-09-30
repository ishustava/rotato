package cmd

type RotatorCommand struct {
	AddNewCAs AddNewCAsCommand `command:"add-new-cas" description:"Regenerate new CAs and make them available to the next BOSH deploy"`
	RegenerateCerts RegenerateCertsCommand `command:"regenerate-certs" description:"Regenerate new certificates"`
	RemoveOldCAs RemoveOldCAsCommand `command:"remove-old-cas" description:"Remove Old CAs"`
}

var Rotator RotatorCommand