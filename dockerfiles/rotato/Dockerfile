FROM golang:1.11

RUN wget https://github.com/cloudfoundry/bosh-bootloader/releases/download/v6.10.3/bbl-v6.10.3_linux_x86-64 -O /usr/local/bin/bbl && \
    chmod +x /usr/local/bin/bbl

RUN go get -u github.com/ishustava/rotato/rotation && mv $GOPATH/bin/rotation $GOPATH/bin/rotato
