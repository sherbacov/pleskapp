package cmd

import (
	"errors"
	"fmt"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [DOMAIN] [IPv4] || create [SERVER] [DOMAIN] [IPv4]",
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

		domain := args[0]
		ipv4 := args[1]

		lookupDomain, sub, errLookup := findDomain(domain, server)
		if errLookup != nil {
			return errLookup
		}

		fmt.Println("domain: ", lookupDomain)

		return actions.DnsRecordCreate(*server, lookupDomain.Name, sub, ipv4)

		return nil
	},
	Args: cobra.MaximumNArgs(2),
}

var createCmdOld = &cobra.Command{
	Use:   "create [DOMAIN] [IPv4] || create [SERVER] [DOMAIN] [IPv4]",
	Short: locales.L.Get("database.create.description"),
	RunE: func(cmd *cobra.Command, args []string) error {
		var realType string = ""
		dbt, _ := cmd.Flags().GetString("type")
		for _, i := range []string{"A", "NS", "MX"} {
			if dbt == i {
				realType = dbt
			}
		}

		if realType == "" {
			return errors.New(locales.L.Get("errors.unknown.database.type", dbt))
		}

		server, err := config.GetServer(args[0])
		if err != nil {
			return err
		}

		dbs := server.GetDatabaseServerByType(realType)
		if dbs == nil {
			return types.DatabaseServerNotFound{
				DbType: realType,
				Server: server.Host,
			}
		}

		domain, err := config.GetDomain(*server, args[1])
		if err != nil {
			return err
		}

		db := types.NewDatabase{
			Name:             args[2],
			Type:             realType,
			ParentDomain:     domain.Name,
			DatabaseServerID: dbs.ID,
		}

		cmd.SilenceUsage = true
		err = actions.DatabaseAdd(*server, *domain, *dbs, db)

		if err == nil {
			fmt.Println(locales.L.Get("database.create.success", db.Name))
		}

		return err
	},
	Args: cobra.ExactArgs(3),
}

func init() {
	//createCmd.Flags().String("type", "A", locales.L.Get("database.create.flag.type"))
	//DnsCmd.AddCommand(createCmd)
	DnsCmd.AddCommand(createCmd)
}
