package goaviatrix

import (
	"context"
	"encoding/json"
)

type OptionSet struct {
	CID           string `json:"CID"`
	Action        string `json:"action"`
	Name          string `json:"option_set_name"`
	SearchDomains []SearchDomains
	NameServers   []NameServers
}

type Result struct {
	Return    bool          `json:"return"`
	Results   OptionSetRead `json:"results"`
	Reason    string        `json:"reason"`
	ErrorType string        `json:"errortype"`
}

type OptionSetRead struct {
	Name          string          `json:"_id"`
	NameServers   []NameServers   `json:"name_servers"`
	SearchDomains []SearchDomains `json:"search_domains"`
}

type SearchDomains struct {
	SearchDomain string `json:"search_domain"`
	Server       string `json:"server"`
}

type NameServers struct {
	Server    string `json:"server"`
	Dot       bool   `json:"dot"`
	Transport string `json:"transport"`
}

func (c *Client) CreateOptionSet(ctx context.Context, optionSet *OptionSet) error {
	optionSet.CID = c.CID
	optionSet.Action = "add_option_set"
	searchDomainsJson, err := json.Marshal(optionSet.SearchDomains)
	if err != nil {
		return err
	}
	nameServersJson, err := json.Marshal(optionSet.NameServers)
	if err != nil {
		return err
	}
	form := map[string]interface{}{
		"CID":             c.CID,
		"action":          optionSet.Action,
		"option_set_name": optionSet.Name,
		"search_domains":  string(searchDomainsJson),
		"name_servers":    string(nameServersJson),
	}
	return c.PostAPIContext2(ctx, nil, "add_option_set", form, BasicCheck)
}

func (c *Client) GetOptionSet(ctx context.Context, name string) (*OptionSetRead, error) {
	action := "get_option_set"
	form := map[string]string{
		"CID":             c.CID,
		"action":          action,
		"option_set_name": name,
	}

	var resp Result
	err := c.PostAPIContext2(ctx, &resp, "get_option_set", form, BasicCheck)
	if err != nil {
		return nil, err
	}
	return &resp.Results, nil
}

func (c *Client) UpdateOptionSet(ctx context.Context, optionSet *OptionSet) error {
	//optionSet.CID = c.CID
	//action := "update_option_set"

	optionSet.CID = c.CID
	optionSet.Action = "update_option_set"
	searchDomainsJson, err := json.Marshal(optionSet.SearchDomains)
	if err != nil {
		return err
	}
	nameServersJson, err := json.Marshal(optionSet.NameServers)
	if err != nil {
		return err
	}
	form := map[string]interface{}{
		"CID":             c.CID,
		"action":          optionSet.Action,
		"option_set_name": optionSet.Name,
		"search_domains":  string(searchDomainsJson),
		"name_servers":    string(nameServersJson),
	}

	return c.PostAPIContext2(ctx, nil, "update_option_set", form, BasicCheck)
}

func (c *Client) DeleteOptionSet(ctx context.Context, name string) error {
	action := "delete_option_set"

	form := map[string]string{
		"CID":             c.CID,
		"action":          action,
		"option_set_name": name,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}
