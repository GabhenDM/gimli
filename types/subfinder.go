package types

//SubfinderHost contains each element found via subfinder
type SubfinderHost struct {
	Host   string `json:"host"`
	Source string `json:"source"`
}

// SubfinderOutput contain the output from the Subfinder tool
type SubfinderOutput struct {
	Subdomains []SubfinderHost
}
