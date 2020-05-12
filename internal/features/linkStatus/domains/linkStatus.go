package lnkstatusdomains

// LinkLookup : LinkLookup
type LinkLookup struct {
	Link string
	Key  string
}

// InfringingLinkStatus : InfringingLinkStatus
type InfringingLinkStatus struct {
	Link      string `json:"url"`
	Duplicate bool   `json:"duplicate"`
	Status    string `json:"status,omitempty"`
}
