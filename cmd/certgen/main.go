package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/gmalbrand/http-tools/pkg/certificates"
	"github.com/gmalbrand/http-tools/pkg/utils"
	log "github.com/sirupsen/logrus"
)

var (
	output       = flag.String("output", "/tmp", "Certificate will be saved to this directory")
	timeToLive   = flag.Int("ttl", 30, "Time to live in days")
	organization = flag.String("organization", "dilizone.org", "Certificates' Organization")
	filePrefix   = flag.String("prefix", "certgen", "File prefix")
)

func main() {
	utils.InitLog()
	flag.Parse()
	log.Debugf("Certificate will be saved to %s", *output)

	if len(flag.Args()) == 0 {
		log.Fatal("You must provide DNS name list as argument")
	}

	caCert, caPrivateKey, err := certificates.CreateSelfSignedCA(*timeToLive)

	if err != nil {
		log.Fatal(err.Error())
	}

	certificates.WriteCertificate(filepath.Join(*output, fmt.Sprintf("%s-ca-crt.pem", *filePrefix)), caCert)
	certificates.WriteKey(filepath.Join(*output, fmt.Sprintf("%s-ca-key.pem", *filePrefix)), caPrivateKey)

	serverCert, serverPrivateKey, err := certificates.CreateServerCertificate(caCert, caPrivateKey, *timeToLive, flag.Args())

	if err != nil {
		log.Fatal(err.Error())
	}

	certificates.WriteCertificate(filepath.Join(*output, fmt.Sprintf("%s-server-crt.pem", *filePrefix)), serverCert)
	certificates.WriteKey(filepath.Join(*output, fmt.Sprintf("%s-server-key.pem", *filePrefix)), serverPrivateKey)

}
