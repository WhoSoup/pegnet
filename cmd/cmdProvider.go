package cmd

import (
	"fmt"

	"github.com/pegnet/pegnet/common"

	"github.com/spf13/cobra"
)

// CmdFlagProvider is able to pull values from the command line, and replace default config
// values.
type CmdFlagProvider struct {
	cmd *cobra.Command
}

func NewCmdFlagProvider(cmd *cobra.Command) *CmdFlagProvider {
	d := new(CmdFlagProvider)
	d.cmd = cmd

	return d
}

func (c *CmdFlagProvider) Load() (map[string]string, error) {
	var err error
	settings := map[string]string{}

	miners, _ := c.cmd.Flags().GetInt("miners")
	if miners != -1 {
		settings["Miner.NumberOfMiners"] = fmt.Sprintf("%d", miners)
	}

	top, _ := c.cmd.Flags().GetInt("top")
	if top != -1 {
		settings["Miner.RecordsPerBlock"] = fmt.Sprintf("%d", top)
	}

	id, _ := c.cmd.Flags().GetString("identity")
	if id != "" {
		settings["Miner.IdentityChain"] = id
	}

	// Use the same flag for the client and server.
	addr, _ := c.cmd.Flags().GetString("caddr")
	if addr != "" {
		settings[common.ConfigCoordinatorListen] = addr
		settings[common.ConfigCoordinatorListen] = addr
	}

	network, _ := c.cmd.Flags().GetString("network")
	if network != "" {
		settings[common.ConfigPegnetNetwork], err = common.GetNetwork(network)
		if err != nil {
			return settings, err
		}
	}

	return settings, nil
}
