terraform {
    required_providers {
        deckofcards = {
            source = "hashicorp.com/Varun-Mehrotra/deckofcards"
        }
    }
}

data "deckofcards_deck" "foobar" {
    deck_id = "zs2mh26mfpm3"
}

resource "deckofcards_deck" "foobar" { }
