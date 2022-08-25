package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixOptionSet() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixOptionSetCreate,
		ReadWithoutTimeout:   resourceAviatrixOptionSetRead,
		UpdateWithoutTimeout: resourceAviatrixOptionSetUpdate,
		DeleteWithoutTimeout: resourceAviatrixOptionSetDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Option Set name.",
			},
			"search_domains": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
					},
				},
			},
			"name_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
						"dot": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "",
						},
						"transport": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixOptionSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)
	optionSet := &goaviatrix.OptionSet{
		Name:          d.Get("name").(string),
		SearchDomains: make([]goaviatrix.SearchDomains, 0),
		NameServers:   make([]goaviatrix.NameServers, 0),
	}

	if optionSet.Name == "" {
		return diag.Errorf("optionSet name can't be empty string")
	}

	searchDomains := d.Get("search_domains").([]interface{})

	for _, searchDomain := range searchDomains {
		if searchDomain != nil {
			domain := searchDomain.(map[string]interface{})
			searchDomainObject := &goaviatrix.SearchDomains{
				SearchDomain: domain["search_domain"].(string),
				Server:       domain["server"].(string),
			}
			optionSet.SearchDomains = append(optionSet.SearchDomains, *searchDomainObject)
		}
	}

	nameServers := d.Get("name_servers").([]interface{})

	for _, nameServer := range nameServers {
		if nameServer != nil {
			server := nameServer.(map[string]interface{})
			nameServerObject := &goaviatrix.NameServers{
				Server:    server["server"].(string),
				Dot:       server["dot"].(bool),
				Transport: server["transport"].(string),
			}
			optionSet.NameServers = append(optionSet.NameServers, *nameServerObject)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	flag := false
	defer resourceAviatrixOptionSetReadIfRequired(ctx, d, meta, &flag)
	err := client.CreateOptionSet(ctx, optionSet)
	if err != nil {
		return diag.Errorf("failed to create the option set: %s", err)
	}

	return resourceAviatrixOptionSetReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixOptionSetReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixOptionSetRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixOptionSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	optionSetName := d.Get("name").(string)
	optionSet, err := client.GetOptionSet(ctx, optionSetName)

	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read option set detail: %s", err)
	}

	d.Set("name", optionSet.Name)

	var searchDomains []map[string]interface{}
	if optionSet.SearchDomains != nil {
		for _, searchDomain := range optionSet.SearchDomains {
			searchDomainObject := make(map[string]interface{})
			searchDomainObject["search_domain"] = searchDomain.SearchDomain
			searchDomainObject["server"] = searchDomain.Server
			searchDomains = append(searchDomains, searchDomainObject)
		}
	}

	if err := d.Set("search_domains", searchDomains); err != nil {
		return diag.Errorf("fail to read search_domains detail: %s", err)
	}

	var nameServers []map[string]interface{}
	if optionSet.NameServers != nil {
		for _, nameServer := range optionSet.NameServers {
			nameServerObject := make(map[string]interface{})
			nameServerObject["server"] = nameServer.Server
			nameServerObject["dot"] = nameServer.Dot
			nameServerObject["transport"] = nameServer.Transport
			nameServers = append(nameServers, nameServerObject)
		}
	}
	if err := d.Set("name_servers", nameServers); err != nil {
		return diag.Errorf("fail to read name_servers detail: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixOptionSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	optionSet := &goaviatrix.OptionSet{
		Name: d.Get("name").(string),
	}
	searchDomainsUpdate := false
	nameServersUpdate := false

	if d.HasChanges("search_domains") {
		searchDomainsUpdate = true
	}
	if d.HasChange("name_servers") {
		nameServersUpdate = true
	}

	if searchDomainsUpdate || nameServersUpdate {
		searchDomains := d.Get("search_domains").([]interface{})
		for _, searchDomain := range searchDomains {
			if searchDomain != nil {
				domain := searchDomain.(map[string]interface{})
				searchDomainObject := &goaviatrix.SearchDomains{
					SearchDomain: domain["search_domain"].(string),
					Server:       domain["server"].(string),
				}
				optionSet.SearchDomains = append(optionSet.SearchDomains, *searchDomainObject)
			}
		}
		nameServers := d.Get("name_servers").([]interface{})
		for _, nameServer := range nameServers {
			if nameServer != nil {
				server := nameServer.(map[string]interface{})
				nameServerObject := &goaviatrix.NameServers{
					Server:    server["server"].(string),
					Dot:       server["dot"].(bool),
					Transport: server["transport"].(string),
				}
				optionSet.NameServers = append(optionSet.NameServers, *nameServerObject)
			}
		}
		err := client.UpdateOptionSet(ctx, optionSet)
		if err != nil {
			return diag.Errorf("failed to update option set detail: %s", err)
		}
	}

	return resourceAviatrixOptionSetRead(ctx, d, meta)
}

func resourceAviatrixOptionSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	name := d.Get("name").(string)

	err := client.DeleteOptionSet(ctx, name)
	if err != nil {
		return diag.Errorf("failed to delete option set detail: %s", err)
	}
	return nil
}
