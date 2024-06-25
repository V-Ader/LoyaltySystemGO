package issuer

type Issuer struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type IssuerDataRequest struct {
	Name string `json:"name"`
}

type IssuerPatchRequest struct {
	Name *string `json:"name,omitempty"`
}
