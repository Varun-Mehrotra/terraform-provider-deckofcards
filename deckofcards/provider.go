package deckofcards

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"deckofcards_deck": resourceDeck(),
			"deckofcards_pile": resourcePile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"deckofcards_deck": dataSourceDeck(),
		},
	}
}
