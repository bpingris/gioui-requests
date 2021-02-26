package service

type Fetcher struct {
	cnt uint64
}

type Header struct {
	Key   string
	Value string
}

type Headers []Header

type FetchPayload struct {
	Method  Method
	URL     string
	Headers Headers
}

func (f *Fetcher) Fetch(p FetchPayload) string {
	res, err := p.Method.Request(p.URL, p.Headers)
	if err != nil {
		return "An error occured: " + err.Error()
	}
	return res
}
