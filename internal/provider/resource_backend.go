package provider

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/matthisholleville/terraform-provider-haproxy/internal/haproxy"
	"github.com/matthisholleville/terraform-provider-haproxy/internal/haproxy/models"
)

func resourceBackend() *schema.Resource {
	return &schema.Resource{
		Description:   "`haproxy_backend` manage backend.",
		CreateContext: resourceBackendCreate,
		ReadContext:   resourceBackendRead,
		UpdateContext: resourceBackendUpdate,
		DeleteContext: resourceBackendDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the octal mode used to define access permissions on the UNIX socket. Possible value 'http' or 'tcp'. https://cbonte.github.io/haproxy-dconv/1.8/configuration.html#5.1-mode",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Backend name",
			},
			"balance_algorithm": {
				Type:  schema.TypeString,
				Required: true,
				Description: "Backend balance algorithm - Possible values: roundrobin, leastconn",
			},
		},
	}
}

func resourceBackendRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*haproxy.Client)

	backend := models.Backend{
		Name: d.Id(),
	}

	result, err := client.GetBackend(backend)
	if err != nil {
		return diag.FromErr(err)

	}

	d.Set("name", result.Name)
	d.Set("mode", result.Mode)
	d.Set("balance_algorithm", result.Balance.Algorithm)

	return nil
}

func resourceBackendCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*haproxy.Client)
	backend := *buildBackendFromResourceParameters(d)

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

			_, err = client.CreateBackend(transaction.Id, backend)
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
	d.SetId(backend.Name)
	return nil
}

func resourceBackendUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*haproxy.Client)

	backend := *buildBackendFromResourceParameters(d)
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
			_, err = client.UpdateBackend(transaction.Id, backend)
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

func resourceBackendDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*haproxy.Client)

	backend := *buildBackendFromResourceParameters(d)

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

			err = client.DeleteBackend(transaction.Id, backend)
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

func buildBackendFromResourceParameters(d *schema.ResourceData) *models.Backend {
	backend := &models.Backend{}

	if v, ok := d.GetOk("mode"); ok {
		backend.Mode = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		backend.Name = v.(string)
	}

	if v, ok := d.GetOk("balance_algorithm"); ok {
		backend.Balance.Algorithm = v.(string)
	}

	return backend
}
