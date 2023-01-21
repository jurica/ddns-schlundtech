package schlundtech

import (
	"bytes"
	"embed"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

var apiUrl = "https://gateway.schlundtech.de"

var (
	//go:embed templates
	res embed.FS
)

func genZoneInquireTasks0205(user string, password string, context string, token string, zone string) (bytes.Buffer, error) {
	var raw bytes.Buffer
	tpl, err := template.ParseFS(res, "templates/0205.xml")
	if err != nil {
		return raw, err
	}

	data := map[string]interface{}{
		"user":     user,
		"password": password,
		"context":  context,
		// "token":    token,
		"zone": zone,
	}

	err = tpl.Execute(&raw, data)
	if err != nil {
		return raw, err
	}

	return raw, nil
}

func sendZoneInquireTasks0205(user string, password string, context string, token string, zone string) (Zone, error) {
	raw, err := genZoneInquireTasks0205(user, password, context, token, zone)
	if err != nil {
		return Zone{}, err
	}

	resp, err := http.Post(apiUrl, "application/xml", bytes.NewReader(raw.Bytes()))
	if err != nil {
		return Zone{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var response Response
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return Zone{}, err
	}
	if response.Result.Data.Zone.Name == "" {
		return Zone{}, errors.New("response does not contain zone information")
	}

	return response.Result.Data.Zone, nil
}

func genZoneUpdateBulk0202001(user string, password string, context string, token string, zone *Zone, rrName string, newIp string) (bytes.Buffer, error) {
	var raw bytes.Buffer
	tpl, err := template.ParseFS(res, "templates/0202001.xml")
	if err != nil {
		return raw, err
	}

	rr, err := getCurrentRR(zone, rrName)
	if err != nil {
		return raw, err
	}

	data := map[string]interface{}{
		"user":     user,
		"password": password,
		"context":  context,
		// "token":    token,
		"zone":             zone.Name,
		"system_ns":        zone.SystemNS,
		"rr_name":          rr.Name,
		"rr_type":          rr.Type,
		"rr_ttl":           rr.Ttl,
		"rr_value_current": rr.Value,
		"rr_value_updated": newIp,
	}

	tpl.Execute(&raw, data)

	return raw, nil
}

func sendZoneUpdateBulk0202001(user string, password string, context string, token string, zone *Zone, rrName string, newIp string) error {
	rawReq, err := genZoneUpdateBulk0202001(user, password, context, token, zone, rrName, newIp)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiUrl, "application/xml", bytes.NewReader(rawReq.Bytes()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func getZone(domain string) (string, error) {
	idx := strings.Index(domain, ".")
	if idx < 0 {
		return "", errors.New("invalid domain given")
	}

	return domain[idx+1:], nil
}

func getRrName(domain string) (string, error) {
	idx := strings.Index(domain, ".")
	if idx < 0 {
		return "", errors.New("invalid domain given")
	}

	return domain[0:idx], nil
}

func getCurrentRR(zone *Zone, subdomain string) (*RR, error) {
	for _, s := range zone.Rrs {
		if s.Name == subdomain {
			return &s, nil
		}
	}

	msg := fmt.Sprintf("no rr found for '%s' in zone '%s'", subdomain, zone.Name)
	return nil, errors.New(msg)
}

func UpdateDdnsRecord(user string, password string, context string, token string, domain string, newIp string) error {
	rrName, err := getRrName(domain)
	if err != nil {
		return err
	}
	zoneName, err := getZone(domain)
	if err != nil {
		return err
	}

	zone, err := sendZoneInquireTasks0205(user, password, context, token, zoneName)
	if err != nil {
		return err
	}

	err = sendZoneUpdateBulk0202001(user, password, context, token, &zone, rrName, newIp)
	if err != nil {
		return err
	}

	return nil
}
