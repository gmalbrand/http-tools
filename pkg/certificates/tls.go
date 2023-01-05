package certificates

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewCertificateTemplate(caCert bool, timeToLive int) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(int64(time.Now().Year())),
		Subject: pkix.Name{
			Organization:  []string{""},
			Country:       []string{""},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, timeToLive),
		IsCA:                  caCert,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

func ReadCertificate(filepath string) (*x509.Certificate, error) {
	f, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Errorf(err.Error())
		return nil, fmt.Errorf("Unable to read file %s", filepath)
	}

	block, _ := pem.Decode(f)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Errorf(err.Error())
		return nil, fmt.Errorf("Unable to parse certificate from file : %s", filepath)
	}

	return cert, nil
}

func ReadPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	f, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Errorf(err.Error())
		return nil, fmt.Errorf("Unable to read file %s", filepath)
	}

	block, _ := pem.Decode(f)

	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Errorf(err.Error())
		return nil, fmt.Errorf("Unable to parse certificate from file : %s", filepath)
	}

	key, _ := k.(*rsa.PrivateKey)
	return key, nil
}

func CreateSelfSignedCA(timeToLive int) (*x509.Certificate, *rsa.PrivateKey, error) {
	caTemplate := NewCertificateTemplate(true, timeToLive)
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		log.Errorf(err.Error())
		return nil, nil, errors.New("Fails at generating CA Private key")
	}

	c, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPrivateKey.PublicKey, caPrivateKey)

	if err != nil {
		log.Error(err.Error())
		return nil, nil, errors.New("Fails at generating CA Certificate")
	}

	caCert, err := x509.ParseCertificate(c)

	if err != nil {
		log.Error(err.Error())
		return nil, nil, errors.New("Fails at parsing generated certificate")
	}

	return caCert, caPrivateKey, nil
}

func CreateServerCertificate(caCert *x509.Certificate, caPrivateKey *rsa.PrivateKey, timeToLive int, dnsNames []string) (*x509.Certificate, *rsa.PrivateKey, error) {
	certTemplate := NewCertificateTemplate(false, timeToLive)
	certTemplate.DNSNames = dnsNames
	certTemplate.Subject.CommonName = dnsNames[0]

	serverPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		log.Errorf(err.Error())
		return nil, nil, errors.New("Fails at generating server private key")
	}

	c, err := x509.CreateCertificate(rand.Reader, certTemplate, caCert, &serverPrivateKey.PublicKey, caPrivateKey)

	if err != nil {
		log.Error(err.Error())
		return nil, nil, errors.New("Fails at generating CA Certificate")
	}

	serverCertificate, err := x509.ParseCertificate(c)

	if err != nil {
		log.Error(err.Error())
		return nil, nil, errors.New("Fails at parsing generated certificate")
	}

	return serverCertificate, serverPrivateKey, nil
}

func WriteCertificate(filepath string, cert *x509.Certificate) error {
	f, err := os.Create(filepath)

	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("Fails at creating file %s", filepath)
	}
	defer f.Close()

	privateKeyBuf := new(bytes.Buffer)
	_ = pem.Encode(privateKeyBuf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	f.Write(privateKeyBuf.Bytes())

	return nil
}

func WriteKey(filepath string, key *rsa.PrivateKey) error {
	f, err := os.Create(filepath)

	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("Fails at creating file %s", filepath)
	}
	defer f.Close()

	privateKeyBuf := new(bytes.Buffer)
	_ = pem.Encode(privateKeyBuf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	f.Write(privateKeyBuf.Bytes())

	return nil
}
