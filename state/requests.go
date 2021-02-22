package state

type Method int

const (
	GET Method = iota
	POST
	DELETE
	PUT
)

func (m Method) String() string {
	return map[Method]string{
		GET:    "GET",
		POST:   "POST",
		DELETE: "DELETE",
		PUT:    "PUT",
	}[m]
}

type Request struct {
	Method Method
	URL    string
	Name   string
}

type Requests []Request
