/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric/integration/lifecycle/chaincode/caller"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func main() {
	err := shim.Start(&caller.CC{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exiting caller chaincode: %s", err)
		os.Exit(2)
	}
}
