package factory

import (
	"errors"
	"fmt"

	"github.com/VoneChain-CS/fabric-gm/bccsp"
	"github.com/VoneChain-CS/fabric-gm/bccsp/gm"
)

const (
	// GMBasedFactoryName is the name of the factory of the software-based BCCSP implementation
	GMBasedFactoryName = "GM"
)

// GMFactory is the factory of the GMbased BCCSP.
type GMFactory struct{}

// Name returns the name of this factory
func (f *GMFactory) Name() string {
	return GMBasedFactoryName
}

// Get returns an instance of BCCSP using Opts.
func (f *GMFactory) Get(config *FactoryOpts) (bccsp.BCCSP, error) {
	// Validate arguments
	if config == nil || config.SwOpts == nil {
		return nil, errors.New("Invalid config. It must not be nil.")
	}

	gmOpts := config.SwOpts

	var ks bccsp.KeyStore
	if gmOpts.Ephemeral == true {
		ks = gm.NewDummyKeyStore()
	} else if gmOpts.FileKeystore != nil {
		fks, err := gm.NewFileBasedKeyStore(nil, gmOpts.FileKeystore.KeyStorePath, false)
		if err != nil {
			return nil, fmt.Errorf("Failed to initialize gm software key store: %s", err)
		}
		ks = fks
	} else {
		// Default to DummyKeystore
		ks = gm.NewDummyKeyStore()
	}

	return gm.New(gmOpts.SecLevel, "GMSM3", ks)
	//return gm.New(gmOpts.SecLevel, gmOpts.HashFamily, ks)
}
