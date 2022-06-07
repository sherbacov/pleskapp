// Copyright 1999-2021. Plesk International GmbH.

package cmd

import (
	"errors"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [DOMAIN]",
	Short: locales.L.Get("domain.list.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		if len(args) == 0 {
			return errors.New("domain is not provided")
		}

		var serverName, err = config.DefaultServer()
		if err != nil {
			return err
		}

		var server, errServer = config.GetServer(serverName)
		if errServer != nil {
			return errServer
		}

		return actions.DnsList(*server, args[0])
	},
	Args: cobra.MaximumNArgs(2),
}

func init() {
	DnsCmd.AddCommand(listCmd)
}
