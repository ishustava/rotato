# Rotato

A tool for CA and certificate rotation for Cloud Foundry

- Uses CredHub API
- Provides commands to implement <a href='#3-step-rotation'>3-step CA rotation</a> 

## <a name='3-step-rotation'></a>3-step CA Rotation

1. Generate new CA, configure all components to trust both old and new CAs and redeploy
2. Generate certificates signed by the new CA
3. Configure everything to only trust the new CA and redeploy

## Installation

`rotato` uses go modules to manages dependencies and requires go `v1.11`. To install, run these commands from the root project directory

```
cd rotation
go build -o rotato
```

## Usage

```
Usage:
  rotation [OPTIONS] [add-new-cas | regenerate-certs | remove-old-cas]

Help Options:
  -h, --help  Show this help message

Available commands:
  add-new-cas       Regenerate new CAs and make them available to the next BOSH deploy
  regenerate-certs  Regenerate new certificates
  remove-old-cas    Remove Old CAs
``` 

Each command requires CredHub credentials. 
In case you are using `bbl` to create your BOSH director, `eval "$(eval bbl print-env)"` will set all necessary credentials.