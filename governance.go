package dcrlibwallet

import (
	"time"

	pwww "github.com/planetdecred/dcrlibwallet/politeiawww"
)

type Politeia struct {
	inventory *pwww.TokenInventory
	client    pwww.Client
}

func NewPoliteia(timeoutSeconds int64) Politeia {
	return Politeia{
		client: pwww.NewClient(time.Duration(timeoutSeconds * int64(time.Second))),
	}
}
