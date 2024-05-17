package client

type Client struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ClientUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ClientPatchRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}
