package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAviatrixAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixAccountRead,

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"aws_account_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//"aws_iam": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},		# REST API needs to support this
			"aws_role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_role_ec2": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//"aws_secret_key": {
			//	Type:      schema.TypeString,
			//	Computed:  true,
			//	Sensitive: true,
			//},
		},
	}
}

func dataSourceAviatrixAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	account := &goaviatrix.Account{
		AccountName: d.Get("account_name").(string),
	}
	log.Printf("[INFO] Looking for Aviatrix account: %#v", account)
	acc, err := client.GetAccount(account)
	if err != nil {
		return fmt.Errorf("aviatrix Account: %s", err)
	}

	if acc != nil {
		d.Set("account_name", acc.AccountName)
		d.Set("cloud_type", acc.CloudType)
		d.Set("aws_account_number", acc.AwsAccountNumber)
		d.Set("aws_access_key", acc.AwsAccessKey)
		d.Set("aws_secret_key", acc.AwsSecretKey)
		d.Set("aws_role_arn", acc.AwsRoleApp)
		d.Set("aws_role_ec2", acc.AwsRoleEc2)
		d.SetId(acc.AccountName)
	} else {
		d.SetId("")
	}
	return nil
}
