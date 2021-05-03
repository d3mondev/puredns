package massdns

// JSONRecord contains a record from a massdns output file.
type JSONRecord struct {
	TTL   int    `json:"ttl"`
	Type  string `json:"type"`
	Class string `json:"class"`
	Name  string `json:"name"`
	Data  string `json:"data"`
}

// JSONResponseData contains the response data from a massdns output file.
type JSONResponseData struct {
	Answers     []JSONRecord `json:"answers"`
	Authorities []JSONRecord `json:"authorities"`
	Additionals []JSONRecord `json:"additionals"`
}

// JSONResponse contains the response from a massdns output file.
type JSONResponse struct {
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Class    string           `json:"class"`
	Status   string           `json:"status"`
	Data     JSONResponseData `json:"data"`
	Resolver string           `json:"resolver"`
}
