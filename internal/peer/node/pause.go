/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"github.com/VoneChain-CS/fabric-gm/core/ledger/kvledger"
	"github.com/VoneChain-CS/fabric-gm/internal/peer/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path/filepath"
)

func pauseCmd() *cobra.Command {
	pauseChannelCmd.ResetFlags()
	flags := pauseChannelCmd.Flags()
	flags.StringVarP(&channelID, "channelID", "c", common.UndefinedParamValue, "Channel ID")
	flags.StringVarP(&rootFSPath, "rootFSPath", "p", common.UndefinedParamValue, "File system path")

	return pauseChannelCmd
}

var pauseChannelCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pauses a channel on the peer.",
	Long:  `Pauses a channel on the peer. When the command is executed, the peer must be offline. When the peer starts after pause, it will not receive blocks for the paused channel.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if channelID == common.UndefinedParamValue {
			return errors.New("Must supply channel ID")
		}
		if rootFSPath == common.UndefinedParamValue {
			return errors.New("Must supply file system path")
		}
		ledgersPath := filepath.Join(rootFSPath, "ledgersData")
		return kvledger.PauseChannel(ledgersPath, channelID)
	},
}
