/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"io/ioutil"

	"github.com/VoneChain-CS/fabric-gm/common/ledger/blockledger"
	"github.com/VoneChain-CS/fabric-gm/common/ledger/blockledger/fileledger"
	"github.com/VoneChain-CS/fabric-gm/common/metrics"
	config "github.com/VoneChain-CS/fabric-gm/orderer/common/localconfig"
	"github.com/pkg/errors"
)

func createLedgerFactory(conf *config.TopLevel, metricsProvider metrics.Provider) (blockledger.Factory, string, error) {
	ld := conf.FileLedger.Location
	var err error
	if ld == "" {
		if ld, err = ioutil.TempDir("", conf.FileLedger.Prefix); err != nil {
			logger.Panic("Error creating temp dir:", err)
		}
	}

	logger.Debug("Ledger dir:", ld)
	lf, err := fileledger.New(ld, metricsProvider)
	if err != nil {
		return nil, "", errors.WithMessage(err, "Error in opening ledger factory")
	}
	return lf, ld, nil
}
