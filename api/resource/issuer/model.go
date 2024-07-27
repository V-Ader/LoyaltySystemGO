package issuer

import (
	"fmt"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
)

type Issuer struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (client *Issuer) GetHash() string {
	return common.GenerateETag([]byte(fmt.Sprintf("%v", client)))
}

type IssuerDataRequest struct {
	Name string `json:"name"`
}

type IssuerPatchRequest struct {
	Name *string `json:"name,omitempty"`
}
