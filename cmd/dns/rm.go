package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/plesk/pleskapp/plesk/internal/actions"
	api2 "github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/api/factory"
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/locales"
	"github.com/spf13/cobra"
	"strings"
)

var rmCmd = &cobra.Command{
	Use:   "rm DOMAIN",
	Short: locales.L.Get("domain.delete.description"),
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

		//color.Cyan("Domain: ", domain)

		// removing last dot.
		domain = strings.TrimRight(domain, ".")

		lookupDomain, _, errLookup := findDomain(domain, server)
		if errLookup != nil {
			return errLookup
		}

		//searching records
		api := factory.GetDnsRecordManagement(server.GetServerAuth())
		records, err := api.ListDnsRecords(lookupDomain.Name)
		if err != nil {
			return err
		}

		var matchRecords []api2.DnsRecord

		n := 0
		for _, val := range records {
			if strings.TrimRight(val.Host, ".") == domain {
				matchRecords = append(matchRecords, val)
				n++
			}
		}

		if n == 0 {
			color.Red("Record not found.")
			return nil
		}

		// only one match
		if n == 1 {
			record := matchRecords[0]
			//fmt.Println(matchRecords)
			return actions.DnsRecordRemove(*server, *lookupDomain, record)

		} else {
			ip, _ := cmd.Flags().GetString("ip")
			if ip != "" {
				for _, val := range records {
					if val.Value == ip && strings.TrimRight(val.Host, ".") == domain {
						//We got record
						return actions.DnsRecordRemove(*server, *lookupDomain, val)
					}
				}
			}

			// multiple match
			color.Red("We found multiple records, please provide --ip or --all for remove all records")
			fmt.Println(matchRecords)
		}

		return nil
	},

	Args: cobra.MaximumNArgs(4),
	//Args: cobra.ExactArgs(3),
}

func init() {
	rmCmd.Flags().String("ip", "", "IP address for filter record")
	rmCmd.Flags().Bool("all", false, "Remove all matched record")
	DnsCmd.AddCommand(rmCmd)
}
