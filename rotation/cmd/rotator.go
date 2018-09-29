package cmd

type RotatorCommand struct {
	AddNewCAs AddNewCAsCommand `command:"add-new-cas" alias:"a" description:"Regenerate new CAs and make them available to the next BOSH deploy"`
}

var Rotator RotatorCommand