// +build !pkcs11

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package factory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitFactories(t *testing.T) {
	err := initFactories(&FactoryOpts{
		ProviderName: "GM",
		SwOpts:       &SwOpts{},
	})
	assert.EqualError(t, err, "Failed initializing BCCSP: Could not initialize BCCSP GM [Failed initializing configuration at [0,]: Hash Family not supported []]")

	err = initFactories(&FactoryOpts{
		ProviderName: "PKCS11",
	})
	assert.EqualError(t, err, "Could not find default `PKCS11` BCCSP")
}
