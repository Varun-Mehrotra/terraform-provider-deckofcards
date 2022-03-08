package deckofcards

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"time"
)

func dataSourceDeck() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeckRead,
		Schema: map[string]*schema.Schema{
			"deck_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"shuffled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"remaining": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDeckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
