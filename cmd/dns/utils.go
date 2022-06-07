package cmd

import (
	"github.com/plesk/pleskapp/plesk/internal/config"
	"github.com/plesk/pleskapp/plesk/internal/types"
	"strings"
)

func findDomain(domain string, server *types.Server) (*types.Domain, string, error) {

	//split domain
	domainParts := strings.Split(domain, ".")

	// Generate domain for lookup
	for i := 0; i < len(domainParts)-1; i++ {
		lookupDomain := strings.Join(domainParts[i:], ".")

		domain, err := config.GetDomain(*server, lookupDomain)
		//Domain not found
		if err != nil {
			continue
		}

		sub := strings.Join(domainParts[0:i], ".")

		//fmt.Println("Sub for record: ", sub)
		//fmt.Println("Domain for record: ", domain.Name)

		return domain, sub, nil
	}

	return nil, "", types.DomainNotFound{Domain: domain, Server: server.Host}
}
