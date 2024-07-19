package client

import (
	"fmt"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
)

type Client struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (client *Client) getHash() string {
	return common.GenerateETag([]byte(fmt.Sprintf("%v", client)))
}

type ClientDataRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ClientPatchRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}
