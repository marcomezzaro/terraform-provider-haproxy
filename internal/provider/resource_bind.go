package provider

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/matthisholleville/terraform-provider-haproxy/internal/haproxy"
	"github.com/matthisholleville/terraform-provider-haproxy/internal/haproxy/models"
)

func resourceBind() *schema.Resource {
	return &schema.Resource{
		Description:   "`haproxy_bind` manage bind.",
		CreateContext: resourceBindCreate,
		ReadContext:   resourceBindRead,
		UpdateContext: resourceBindUpdate,
		DeleteContext: resourceBindDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address to bind",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Frontend name",
			},
			"port": {
				Type:  schema.TypeInt,
				Required: true,
				Description: "Bind port",
			},
			"parent_name": {
				Type:  schema.TypeString,
				Required: true,
				Description: "Frontend name related to this bind",
			},
		},
	}
}

func resourceBindRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*haproxy.Client)
	pn:= fmt.Sprintf("%v", d.Get("parent_name"))
	bind := models.Bind{
		Name: d.Id(),
	}

	result, err := client.GetBind(bind, pn)
	if err != nil {
		return diag.FromErr(err)

	}

	d.Set("name", result.Name)
	d.Set("address", result.Address)
	d.Set("port", result.Port)
	d.Set("parent_name", pn)

	return nil
}

func resourceBindCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*haproxy.Client)
	bind := *buildBindFromResourceParameters(d)
	pn:= fmt.Sprintf("%v", d.Get("parent_name"))

	err := retry.Do(
		func() error {
			configuration, err := client.GetConfiguration()
			if err != nil {
				return err
			}
			transaction, err := client.CreateTransaction(configuration.Version)
			if err != nil {
				return err
			}

			_, err = client.CreateBind(transaction.Id, bind, pn)
			if err != nil {
				return err
			}
			_, err = client.CommitTransaction(transaction.Id)
			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(bind.Name)
	return nil
}

func resourceBindUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*haproxy.Client)
	bind := *buildBindFromResourceParameters(d)
	pn:= fmt.Sprintf("%v", d.Get("parent_name"))

	err := retry.Do(
		func() error {
			configuration, err := client.GetConfiguration()
			if err != nil {
				return err
			}
			transaction, err := client.CreateTransaction(configuration.Version)
			if err != nil {
				return err
			}
			_, err = client.UpdateBind(transaction.Id, bind, pn)
			if err != nil {
				return err
			}
			_, err = client.CommitTransaction(transaction.Id)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceBindDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*haproxy.Client)
	bind := *buildBindFromResourceParameters(d)
	pn:= fmt.Sprintf("%v", d.Get("parent_name"))


	err := retry.Do(
		func() error {
			configuration, err := client.GetConfiguration()
			if err != nil {
				return err
			}
			transaction, err := client.CreateTransaction(configuration.Version)
			if err != nil {
				return err
			}

			err = client.DeleteBind(transaction.Id, bind, pn)
			if err != nil {
				return err
			}
			_, err = client.CommitTransaction(transaction.Id)
			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func buildBindFromResourceParameters(d *schema.ResourceData) *models.Bind {
	bind := &models.Bind{}

	if v, ok := d.GetOk("port"); ok {
		bind.Port = v.(int)
	}

	if v, ok := d.GetOk("name"); ok {
		bind.Name = v.(string)
	}

	if v, ok := d.GetOk("address"); ok {
		bind.Address = v.(string)
	}

	return bind
}
