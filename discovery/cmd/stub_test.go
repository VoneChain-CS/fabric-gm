/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"fmt"
	"net"
	"path/filepath"
	"testing"

	"github.com/VoneChain-CS/fabric-gm/cmd/common"
	"github.com/VoneChain-CS/fabric-gm/cmd/common/comm"
	"github.com/VoneChain-CS/fabric-gm/cmd/common/signer"
	discovery "github.com/VoneChain-CS/fabric-gm/discovery/client"
	corecomm "github.com/VoneChain-CS/fabric-gm/internal/pkg/comm"
	"github.com/stretchr/testify/assert"
)

func TestClientStub(t *testing.T) {
	srv, err := corecomm.NewGRPCServer("127.0.0.1:", corecomm.ServerConfig{
		SecOpts: corecomm.SecureOptions{},
	})
	assert.NoError(t, err)
	go srv.Start()
	defer srv.Stop()

	_, portStr, _ := net.SplitHostPort(srv.Address())
	endpoint := fmt.Sprintf("localhost:%s", portStr)
	stub := &ClientStub{}

	req := discovery.NewRequest()

	_, err = stub.Send(endpoint, common.Config{
		SignerConfig: signer.Config{
			MSPID:        "Org1MSP",
			KeyPath:      filepath.Join("testdata", "8150cb2d09628ccc89727611ebb736189f6482747eff9b8aaaa27e9a382d2e93_sk"),
			IdentityPath: filepath.Join("testdata", "cert.pem"),
		},
		TLSConfig: comm.Config{},
	}, req)
	assert.Contains(t, err.Error(), "Unimplemented desc = unknown service discovery.Discovery")
}

func TestRawStub(t *testing.T) {
	srv, err := corecomm.NewGRPCServer("127.0.0.1:", corecomm.ServerConfig{
		SecOpts: corecomm.SecureOptions{},
	})
	assert.NoError(t, err)
	go srv.Start()
	defer srv.Stop()

	_, portStr, _ := net.SplitHostPort(srv.Address())
	endpoint := fmt.Sprintf("localhost:%s", portStr)
	stub := &RawStub{}

	req := discovery.NewRequest()

	_, err = stub.Send(endpoint, common.Config{
		SignerConfig: signer.Config{
			MSPID:        "Org1MSP",
			KeyPath:      filepath.Join("testdata", "8150cb2d09628ccc89727611ebb736189f6482747eff9b8aaaa27e9a382d2e93_sk"),
			IdentityPath: filepath.Join("testdata", "cert.pem"),
		},
		TLSConfig: comm.Config{},
	}, req)
	assert.Contains(t, err.Error(), "Unimplemented desc = unknown service discovery.Discovery")
}
