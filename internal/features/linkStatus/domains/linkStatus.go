package lnkstatusdomains

// LinkLookup : LinkLookup
type LinkLookup struct {
	Link string
	Key  string
}

// LinkStatus : LinkStatus
type LinkStatus struct {
	URL         string `json:"url"`
	PageFoundOn string `json:"pageFoundOn"`
	TaskID      int    `json:"taskId"`
	ProductID   int    `json:"productId"`
	BookmarkID  string `json:"bookmarkId"`
}

// LinkStatusRequestAdmin : LinkStatusRequestAdmin
type LinkStatusRequestAdmin struct {
	Data LinkStatusAdmin `json:"data"`
}

// LinkStatusAdmin : LinkStatusAdmin
type LinkStatusAdmin struct {
	URL         string `json:"url"`
	PageFoundOn string `json:"page_found_on"`
	TaskID      int    `json:"task_id"`
	ProductID   int    `json:"product_id"`
}
