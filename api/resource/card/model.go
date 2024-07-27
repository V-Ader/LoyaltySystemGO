package card

import (
	"fmt"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
)

type Card struct {
	Id        int  `json:"id"`
	Issuer_id int  `json:"issuer_id"`
	Owner_id  int  `json:"owner_id"`
	Active    bool `json:"active"`
	Usages    int  `json:"usages"`
	Capacity  int  `json:"capacity"`
}

func (card *Card) GetHash() string {
	return common.GenerateETag([]byte(fmt.Sprintf("%v", card)))
}

type CardDataRequest struct {
	Issuer_id int  `json:"issuer_id"`
	Owner_id  int  `json:"owner_id"`
	Active    bool `json:"active"`
	Usages    int  `json:"usages"`
	Capacity  int  `json:"capacity"`
}

type CardPatchRequest struct {
	Issuer_id *int  `json:"issuer_id,omitempty"`
	Owner_id  *int  `json:"owner_id,omitempty"`
	Active    *bool `json:"active,omitempty"`
	Usages    *int  `json:"usages,omitempty"`
	Capacity  *int  `json:"capacity,omitempty"`
}
