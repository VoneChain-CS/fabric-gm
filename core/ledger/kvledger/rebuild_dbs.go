/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"github.com/VoneChain-CS/fabric-gm/common/ledger/blkstorage"
	"github.com/VoneChain-CS/fabric-gm/common/ledger/util/leveldbhelper"
	"github.com/VoneChain-CS/fabric-gm/core/ledger"
	"github.com/VoneChain-CS/fabric-gm/core/ledger/kvledger/txmgmt/statedb/statecouchdb"
	"github.com/pkg/errors"
)

// RebuildDBs drops existing ledger databases.
// Dropped database will be rebuilt upon server restart
func RebuildDBs(config *ledger.Config) error {
	rootFSPath := config.RootFSPath
	fileLockPath := fileLockPath(rootFSPath)
	fileLock := leveldbhelper.NewFileLock(fileLockPath)
	if err := fileLock.Lock(); err != nil {
		return errors.Wrap(err, "as another peer node command is executing,"+
			" wait for that command to complete its execution or terminate it before retrying")
	}
	defer fileLock.Unlock()

	if config.StateDBConfig.StateDatabase == "CouchDB" {
		if err := statecouchdb.DropApplicationDBs(config.StateDBConfig.CouchDB); err != nil {
			return err
		}
	}
	if err := dropDBs(rootFSPath); err != nil {
		return err
	}

	blockstorePath := BlockStorePath(rootFSPath)
	return blkstorage.DeleteBlockStoreIndex(blockstorePath)
}
