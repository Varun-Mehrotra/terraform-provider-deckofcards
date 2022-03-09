package deckofcards

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDeck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeckCreate,
		ReadContext:   resourceDeckRead,
		UpdateContext: resourceDeckUpdate,
		DeleteContext: resourceDeckDelete,
		Schema: map[string]*schema.Schema{
			// TODO: Refactor this to just ID
			"deck_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			// TODO: Implement multiple decks
			"shuffled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true, // TODO: Figure out if this is the right option. Not sure the right functionality for updates
				Default:  true,
			},
			"remaining": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"jokers_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			// TODO: Add explicit Cards as Optional
		},
	}
}

func resourceDeckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	shuffleStr := ""

	if isShuffled := d.Get("shuffled").(bool); isShuffled {
		shuffleStr = "/shuffle"
	}

	tflog.Debug(ctx, "Create Deck via /deck/new/")

	req, err := http.NewRequest("GET", fmt.Sprintf(
		"%s/deck/new%s/?jokers_enabled=%t",
		"https://deckofcardsapi.com/api",
		shuffleStr,
		d.Get("jokers_enabled").(bool),
	), nil)

	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	deck := make(map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&deck)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("remaining", int(deck["remaining"].(float64))); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("deck_id", deck["deck_id"].(string)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(deck["deck_id"].(string))

	return diags
}

func resourceDeckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	deckID := d.Get("deck_id").(string)

	tflog.Debug(ctx, "GET Deck")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/deck/%s", "https://deckofcardsapi.com/api", deckID), nil)

	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	deck := make(map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&deck)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("shuffled", deck["shuffled"].(bool)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("remaining", int(deck["remaining"].(float64))); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(deckID)

	return diags
}

func resourceDeckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceDeckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
