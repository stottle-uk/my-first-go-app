package scannertasksdomains

type ScannerTask struct {
	ID   int      `json:"id"`
	List []string `json:"list"`
}
