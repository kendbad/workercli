package model

type Proxy struct {
	Protocol string
	IP       string
	Port     string
}

type ProxyResult struct {
	Proxy  Proxy
	IP     string
	Status string
	Error  string
}
