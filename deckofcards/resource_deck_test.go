package deckofcards

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDeckofcardsDeckBasic(t *testing.T) {
	deckName := "foobar"

	resource.Test(t, resource.TestCase{
		// PreCheck: func() { testAccPreCheck(t)},
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckDeckofcardsDeckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDeckofcardsDeckConfigBasic(deckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeckofcardsDeckExists(deckName),
				),
			},
		},
	})
}

// func testAccCheckDeckofcardsDeckDestroy(s *terraform.State) error {
// 	c := testAccProvider.Meta().(*hc.Client)
//
// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "deckofcards_deck" {
// 			continue
// 		}
//
// 		deckID := rs.Primary.ID
//
// 		err := c.DeleteDeck(deckID)
// 		if err != nil {
// 			return err
// 		}
// 	}
// }

func testAccCheckDeckofcardsDeckConfigBasic(deckName string) string {
	return fmt.Sprintf(`resource "deckofcards_deck" "%s" {

    }`, deckName)
}

func testAccCheckDeckofcardsDeckExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DeckID set")
		}

		return nil
	}
}
