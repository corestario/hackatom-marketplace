package hh

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BaseNFT non fungible token definition
type BaseNFT struct {
	ID          string         `json:"id,omitempty"`       // id of the token; not exported to clients
	Owner       sdk.AccAddress `json:"owner,string"`       // account address that owns the NFT
	Name        string         `json:"name,string"`        // name of the token
	Description string         `json:"description,string"` // unique description of the NFT
	Image       string         `json:"image,string"`       // image path
	TokenURI    string         `json:"token_uri,string"`   // optional extra properties available fo querying
}

// NewBaseNFT creates a new NFT instance
func NewBaseNFT(ID string, owner sdk.AccAddress, tokenURI, description, image, name string) BaseNFT {
	return BaseNFT{
		ID:          ID,
		Owner:       owner,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Image:       strings.TrimSpace(image),
		TokenURI:    strings.TrimSpace(tokenURI),
	}
}

func (m BaseNFT) String() string {
	return fmt.Sprintf(`ID:				%s
Owner:			%s
Name:			%s
Description: 	%s
Image:			%s
TokenURI:		%s`,
		m.ID,
		m.Owner,
		m.Name,
		m.Description,
		m.Image,
		m.TokenURI,
	)
}

type NFT struct {
	BaseNFT
	OnSale bool      `json:"on_sale"`
	Price  sdk.Coins `json:"price"`
}

func (m NFT) String() string {
	return fmt.Sprintf(`
%s
OnSale:			%v`,
		m.BaseNFT.String(),
		m.OnSale,
	)
}

type Transfer struct {
	// TODO: obviously there should be at least "complete", "in progress" and "failed"
	// TODO: states, but I'm not sure we should bother about that.
	Complete bool `json:"complete"`
}

func NewTransfer() Transfer {
	return Transfer{}
}

func (m Transfer) String() string {
	if m.Complete {
		return "Transfer is complete"
	}

	return "Transfer is in progress"
}
