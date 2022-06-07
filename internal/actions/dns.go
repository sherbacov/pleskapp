package actions

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/api/factory"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"os"
	"sort"
	"text/tabwriter"
)

func DnsList(host types.Server, domain string) error {
	api := factory.GetDnsRecordManagement(host.GetServerAuth())
	records, err := api.ListDnsRecords(domain)
	if err != nil {
		return err
	}

	sort.Slice(records, func(i, j int) bool { return records[i].Host < records[j].Host })

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "Host\tType\tValue")
	for _, domain := range records {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", domain.Host, domain.Type, domain.Value)
	}

	return w.Flush()
}

func DnsRecordCreate(host types.Server, domain string, sub string, value string) error {
	api := factory.GetDnsRecordManagement(host.GetServerAuth())
	err := api.CreateDnsRecord(domain, sub, value)
	if err != nil {
		return err
	}

	color.Cyan("record created.")

	return nil
}

func DnsRecordRemove(host types.Server, domain types.Domain, record api.DnsRecord) error {
	api := factory.GetDnsRecordManagement(host.GetServerAuth())

	if record.Type == "A" {
		err := api.RemoveDnsRecordA(domain, record.Host, record.Value)
		if err != nil {
			return err
		}

		color.Cyan("record deleted.")
	}

	return nil
}
