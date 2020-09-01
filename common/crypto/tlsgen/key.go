/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tlsgen

import (
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
	"net"
	"time"

	"github.com/pkg/errors"
)

func newPrivKey() (*sm2.PrivateKey, []byte, error) {
	privateKey, err := sm2.GenerateKey()
	if err != nil {
		return nil, nil, err
	}
	privBytes, err := sm2.MarshalSm2UnecryptedPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, privBytes, nil
}

func newCertTemplate() (sm2.Certificate, error) {
	sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return sm2.Certificate{}, err
	}
	return sm2.Certificate{
		Subject:      pkix.Name{SerialNumber: sn.String()},
		NotBefore:    time.Now().Add(time.Hour * (-24)),
		NotAfter:     time.Now().Add(time.Hour * 24),
		KeyUsage:     sm2.KeyUsageKeyEncipherment | sm2.KeyUsageDigitalSignature,
		SerialNumber: sn,
	}, nil
}

func newCertKeyPair(isCA bool, isServer bool, host string, certSigner crypto.Signer, parent *sm2.Certificate) (*CertKeyPair, error) {
	privateKey, privBytes, err := newPrivKey()
	if err != nil {
		return nil, err
	}

	template, err := newCertTemplate()
	if err != nil {
		return nil, err
	}

	tenYearsFromNow := time.Now().Add(time.Hour * 24 * 365 * 10)
	if isCA {
		template.NotAfter = tenYearsFromNow
		template.IsCA = true
		template.KeyUsage |= sm2.KeyUsageCertSign | sm2.KeyUsageCRLSign
		template.ExtKeyUsage = []sm2.ExtKeyUsage{
			sm2.ExtKeyUsageClientAuth,
			sm2.ExtKeyUsageServerAuth,
		}
		template.BasicConstraintsValid = true
	} else {
		template.ExtKeyUsage = []sm2.ExtKeyUsage{sm2.ExtKeyUsageClientAuth}
	}
	template.SignatureAlgorithm = sm2.SM2WithSM3
	template.SubjectKeyId = computeSKI(privateKey)
	if isServer {
		template.NotAfter = tenYearsFromNow
		template.ExtKeyUsage = append(template.ExtKeyUsage, sm2.ExtKeyUsageServerAuth)
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}
	// If no parent cert, it's a self signed cert
	if parent == nil || certSigner == nil {
		parent = &template
		certSigner = privateKey
	}
	rawBytes, err := sm2.CreateCertificate(rand.Reader, &template, parent, &privateKey.PublicKey, certSigner)
	if err != nil {
		return nil, err
	}
	pubKey := encodePEM("CERTIFICATE", rawBytes)

	block, _ := pem.Decode(pubKey)
	if block == nil { // Never comes unless x509 or pem has bug
		return nil, errors.Errorf("%s: wrong PEM encoding", pubKey)
	}
	cert, err := sm2.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	privKey := encodePEM("EC PRIVATE KEY", privBytes)
	return &CertKeyPair{
		Key:     privKey,
		Cert:    pubKey,
		Signer:  privateKey,
		TLSCert: cert,
	}, nil
}


// compute Subject Key Identifier //TODO Important
func computeSKI(privKey *sm2.PrivateKey) []byte {
	// Marshall the public key
	raw := elliptic.Marshal(privKey.Curve, privKey.PublicKey.X, privKey.PublicKey.Y)
	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
	//return hash[:]
}

func encodePEM(keyType string, data []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: keyType, Bytes: data})
}
