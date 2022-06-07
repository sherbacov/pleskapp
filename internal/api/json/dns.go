package json

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/plesk/pleskapp/plesk/internal/api"
	"github.com/plesk/pleskapp/plesk/internal/types"
)

type jsonDns struct {
	client *resty.Client
}

func NewDns(c *resty.Client) jsonDns {
	return jsonDns{
		client: c,
	}
}

type createDnsRecordRequest struct {
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	ParentDomain domainReference `json:"parent_domain"`
	ServerID     int             `json:"server_id"`
}

func (j jsonDns) ListDnsRecords(domain string) ([]api.DnsRecord, error) {
	var res, err = j.client.R().
		SetResult([]dnsRecord{}).
		SetError(&jsonError{}).
		SetQueryParam("domain", domain).
		Get("/api/v2/dns/records")

	if err != nil {
		return []api.DnsRecord{}, err
	}

	if res.IsSuccess() {
		var r = res.Result().(*[]dnsRecord)
		return jsonDnsRecordToDnsRecords(*r), nil
	}

	if res.StatusCode() == 403 {
		return []api.DnsRecord{}, authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return []api.DnsRecord{}, jsonErrorToError(*r)
}

func (j jsonDns) CreateDnsRecord(domain string, host string, value string) error {
	var p cliGateRequest

	p = cliGateRequest{
		Params: []string{
			"--add",
			domain,
			"-a",
			host,
			"-ip",
			value,
		},
		Env: map[string]string{},
	}

	req, err := json.Marshal(p)
	if err != nil {
		return err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&cliGateResponce{}).
		SetError(&jsonError{}).
		Post("/api/v2/cli/dns/call")

	if err != nil {
		return err
	}

	if res.IsSuccess() {
		var r *cliGateResponce = res.Result().(*cliGateResponce)
		if r.Code != 0 || len(r.Stderr) != 0 {
			return jsonCliGateResponceToError(*r)
		}

		return nil
	}

	if res.StatusCode() == 403 {
		return authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return jsonErrorToError(*r)
}

func (j jsonDns) RemoveDnsRecordA(domain types.Domain, record string, ip string) error {
	var p cliGateRequest

	p = cliGateRequest{
		Params: []string{
			"--del",
			domain.Name,
			"-a",
			record,
			"-ip",
			ip,
		},
		Env: map[string]string{},
	}

	req, err := json.Marshal(p)
	if err != nil {
		return err
	}

	res, err := j.client.R().
		SetBody(req).
		SetResult(&cliGateResponce{}).
		SetError(&jsonError{}).
		Post("/api/v2/cli/dns/call")

	if err != nil {
		return err
	}

	if res.IsSuccess() {
		var r *cliGateResponce = res.Result().(*cliGateResponce)
		if r.Code != 0 || len(r.Stderr) != 0 {
			return jsonCliGateResponceToError(*r)
		}

		return nil
	}

	if res.StatusCode() == 403 {
		return authError{server: j.client.HostURL, needReauth: true}
	}

	var r *jsonError = res.Error().(*jsonError)
	return jsonErrorToError(*r)
}
