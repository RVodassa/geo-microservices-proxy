package entity

type ProfileRequest struct {
	ID uint64 `json:"id"`
}

type ProfileResponse struct {
	ID    uint64 `json:"id,omitempty"`
	Login string `json:"login"`
}

type ListRequest struct {
	Offset uint64 `json:"offset,omitempty"`
	Limit  uint64 `json:"limit,omitempty"`
}

type ListResponse struct {
	Profiles []*ProfileResponse `json:"profiles"`
	Count    uint64             `json:"count"`
}
