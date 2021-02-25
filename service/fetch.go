package service

type Fetcher struct {
	cnt uint64
}

func (f *Fetcher) Fetch(m Method, url string) string {
	res, err := m.Request(url)
	if err != nil {
		return "An error occured: " + err.Error()
	}
	return res
}
