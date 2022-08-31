terraform {
  required_providers {
    deckofcards = {
      source = "hashicorp.com/Varun-Mehrotra/deckofcards"
    }
  }
}

provider "deckofcards" {}

data "deckofcards_deck" "example" {
  deck_id = "70ybwopagdqi"
}

resource "deckofcards_deck" "example" {}

resource "deckofcards_pile" "example" {
  name    = "example"
  deck_id = deckofcards_deck.example.id
  size    = 4
}

output "resource_example_deck_id" {
  value = deckofcards_deck.example.id
}
