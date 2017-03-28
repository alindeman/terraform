package heroku

import (
	"context"

	heroku "github.com/cyberdelia/heroku-go/v3"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHerokuSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceHerokuSpaceCreate,
		Read:   resourceHerokuSpaceRead,
		Update: resourceHerokuSpaceUpdate,
		Delete: resourceHerokuSpaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"shield": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"nat_sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceHerokuSpaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	opts := heroku.SpaceCreateOpts{
		Name:         d.Get("name").(string),
		Organization: d.Get("organization").(string),
		Region:       heroku.String(d.Get("region").(string)),
		Shield:       heroku.Bool(d.Get("shield").(bool)),
	}

	space, err := client.SpaceCreate(context.TODO(), opts)
	if err != nil {
		return err
	}

	d.SetId(space.ID)
	return resourceHerokuSpaceRead(d, meta)
}

func resourceHerokuSpaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	space, err := client.SpaceInfo(context.TODO(), d.Id())
	if err != nil {
		return err
	}

	if space.State == "allocated" {
		nat, err := client.SpaceNatInfo(context.TODO(), d.Id())
		if err != nil {
			return err
		}

		d.Set("nat_sources", nat.Sources)
	}

	d.Set("name", space.Name)
	d.Set("organization", space.Organization.Name)
	d.Set("region", space.Region.Name)
	d.Set("shield", space.Shield)
	return nil
}

func resourceHerokuSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	opts := heroku.SpaceUpdateOpts{
		Name: heroku.String(d.Get("name").(string)),
	}

	space, err := client.SpaceUpdate(context.TODO(), d.Id(), opts)
	if err != nil {
		return err
	}

	d.SetId(space.ID)
	return resourceHerokuSpaceRead(d, meta)
}

func resourceHerokuSpaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	_, err := client.SpaceDelete(context.TODO(), d.Id())
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
