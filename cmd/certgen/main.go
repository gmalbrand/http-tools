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
	certInfo   = certificates.CertificateInfo{}
	output     = flag.String("output", "/tmp", "Certificate will be saved to this directory")
	filePrefix = flag.String("prefix", "certgen", "File prefix")
)

func initFlag() {
	flag.StringVar(&certInfo.Organization, "organization", "", "Name of your organization (Mandatory)")
	flag.StringVar(&certInfo.Country, "country", "", "Country of your organization")
	flag.StringVar(&certInfo.Province, "province", "", "Province of your organization")
	flag.StringVar(&certInfo.Locality, "locality", "", "Locality of your organization")
	flag.StringVar(&certInfo.Province, "street-address", "", "Street address of your organization")
	flag.StringVar(&certInfo.Province, "postal-code", "", "Postal code of your organization")
	flag.IntVar(&certInfo.Validity, "validity", 30, "Validity period in days (default 30")
}

func main() {
	utils.InitLog()
	initFlag()
	flag.Parse()
	log.Debugf("Certificate will be saved to %s", *output)

	if certInfo.Organization == "" {
		log.Fatal("You must provide an organization name")
	}

	if len(flag.Args()) == 0 {
		log.Fatal("You must provide DNS name list as argument")
	}
	certInfo.DNS = flag.Args()

	caCert, caPrivateKey, err := certificates.CreateSelfSignedCA(certInfo)

	if err != nil {
		log.Fatal(err.Error())
	}

	certificates.WriteCertificate(filepath.Join(*output, fmt.Sprintf("%s-ca-cert.pem", *filePrefix)), caCert)
	certificates.WriteKey(filepath.Join(*output, fmt.Sprintf("%s-ca-key.pem", *filePrefix)), caPrivateKey)

	serverCert, serverPrivateKey, err := certificates.CreateServerCertificate(caCert, caPrivateKey, certInfo)

	if err != nil {
		log.Fatal(err.Error())
	}

	certificates.WriteCertificate(filepath.Join(*output, fmt.Sprintf("%s-server-cert.pem", *filePrefix)), serverCert)
	certificates.WriteKey(filepath.Join(*output, fmt.Sprintf("%s-server-key.pem", *filePrefix)), serverPrivateKey)

}
