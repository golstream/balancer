package httputils

type (
	statusCode int
)

type (
	Response struct {
		code statusCode
		body []byte
	}
)

func (r *Response) GetStatusCode() statusCode {
	return r.code
}

func (sc statusCode) Is2xxStatusCode() bool {
	return sc/100 == 2
}

func (sc statusCode) Is4xxStatusCode() bool {
	return sc/100 == 4
}

func (sc statusCode) Is5xxStatusCode() bool {
	return sc/100 == 5
}

func (r *Response) GetIntStatusCode() int {
	return int(r.code)
}

func (r *Response) GetBody() []byte {
	return r.body
}
