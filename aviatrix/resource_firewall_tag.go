package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAviatrixFirewallTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallTagCreate,
		Read:   resourceAviatrixFirewallTagRead,
		Update: resourceAviatrixFirewallTagUpdate,
		Delete: resourceAviatrixFirewallTagDelete,

		Schema: map[string]*schema.Schema{
			"firewall_tag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_tag_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixFirewallTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewall_tag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	err := client.CreateFirewallTag(firewall_tag)
	if err != nil {
		return fmt.Errorf("failed to create firewall tag: %s", err)
	}
	//If cidr list is present, update cidr list
	if _, ok := d.GetOk("cidr_list"); ok {
		cidrList := d.Get("cidr_list").([]interface{})
		for _, currCIDR := range cidrList {
			cm := currCIDR.(map[string]interface{})
			cidrMember := goaviatrix.CIDRMember{
				CIDRTag: cm["cidr_tag_name"].(string),
				CIDR:    cm["cidr"].(string),
			}
			firewall_tag.CIDRList = append(firewall_tag.CIDRList, cidrMember)
		}
		err := client.UpdateFirewallTag(firewall_tag)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix FirewallTag: %s", err)
		}
	}
	d.SetId(firewall_tag.Name)
	return nil
}

func resourceAviatrixFirewallTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewallTag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	fwt, err := client.GetFirewallTag(firewallTag)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error fetching firewall tag %s: %s", firewallTag.Name, err)
	}
	log.Printf("[TRACE] Reading cidr list for tag %s: %#v", firewallTag.Name, fwt)
	if fwt != nil {

		var cidrList []map[string]interface{}
		for _, cidrMember := range fwt.CIDRList {
			cm := make(map[string]interface{})
			cm["cidr_tag_name"] = cidrMember.CIDRTag
			cm["cidr"] = cidrMember.CIDR

			cidrList = append(cidrList, cm)
		}
		d.Set("cidr_list", cidrList)
	}

	return nil
}

func resourceAviatrixFirewallTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewallTag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	d.Partial(true)
	log.Printf("[INFO] Creating Aviatrix firewall: %#v", firewallTag)
	//Update cidr list
	cidrList := d.Get("cidr_list").([]interface{})
	for _, currCIDR := range cidrList {
		cm := currCIDR.(map[string]interface{})
		cidrMember := goaviatrix.CIDRMember{
			CIDRTag: cm["cidr_tag_name"].(string),
			CIDR:    cm["cidr"].(string),
		}
		firewallTag.CIDRList = append(firewallTag.CIDRList, cidrMember)
	}
	err := client.UpdateFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix FirewallTag: %s", err)
	}
	if _, ok := d.GetOk("cidr_list"); ok {
		d.SetPartial("cidr_list")
	}
	d.Partial(false)
	return nil
}

func resourceAviatrixFirewallTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewallTag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	//firewall_tag.CIDRList = make([]*goaviatrix.CIDRMember, 0)
	err := client.UpdateFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix FirewallTag policy list: %s", err)
	}
	err = client.DeleteFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix FirewallTag policy list: %s", err)
	}
	return nil
}
