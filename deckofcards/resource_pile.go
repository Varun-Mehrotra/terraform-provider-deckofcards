package deckofcards

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePileCreate,
		ReadContext:   resourcePileRead,
		UpdateContext: resourcePileUpdate,
		DeleteContext: resourcePileDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// This should be optional, you can init (i.e. draw) from deck or pile
			"deck_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"shuffled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// I think I don't need this? Just let Size exist and be partially read?
			// "remaining": &schema.Schema{
			// 	Type:     schema.TypeInt,
			// 	Computed: true,
			// },
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			// TODO: Implement Card list
		},
	}
}

func resourcePileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	baseURL := "https://deckofcardsapi.com/api"

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	deckID := d.Get("deck_id").(string)
	pileName := d.Get("name").(string)
	pileSize := d.Get("size").(int)

	// Draw cards from deck

	// TODO: Add some spicy concurrency for the reqs (maybe)
	tflog.Debug(ctx, "Draw Cards from Deck")
	req, err := http.NewRequest("GET", fmt.Sprintf(
		"%s/deck/%s/draw/?count=%d",
		baseURL,
		deckID,
		pileSize,
	), nil)

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	defer r.Body.Close()

	cards := make(map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&cards)
	if err != nil {
		return diag.FromErr(err)
	}

	cardsList := []string{}

	// Parse list of cards into string
	for _, i := range cards["cards"].([]interface{}) {
		cardsList = append(cardsList, i.(map[string]interface{})["code"].(string))
	}

	// Make Pile using new cards
	req2, err := http.NewRequest("GET", fmt.Sprintf(
		"%s/deck/%s/pile/%s/add/?cards=%s",
		baseURL,
		deckID,
		pileName,
		strings.Join(cardsList, ","),
	), nil)

	r2, err := client.Do(req2)
	if err != nil {
		return diag.FromErr(err)
	}

	defer r.Body.Close()

	pile := make(map[string]interface{}, 0)
	err = json.NewDecoder(r2.Body).Decode(&pile)
	if err != nil {
		return diag.FromErr(err)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pileName)

	return diags
}

func resourcePileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	baseURL := "https://deckofcardsapi.com/api"

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	deckID := d.Get("deck_id").(string)
	pileName := d.Get("name").(string)

	tflog.Debug(ctx, "GET pile")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/deck/%s/pile/%s/list", baseURL, deckID, pileName), nil)

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

	// pile := make(map[string]interface{}, 0)
	// json.NewDecoder(deck[])

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

func resourcePileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// client := &http.Client{Timeout: 10 * time.Second}

	// ? no idea what to do with this one.

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourcePileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Put all cards back in deck

	return diags
}
