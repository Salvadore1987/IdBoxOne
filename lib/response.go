package lib

type Response struct {
	ICCID string				`json:"iccid"`
	PassportData *PersonalInfo	`json:"passportData"`
	Photo string				`json:"photo"`
}

func NewResponse() *Response {
	return &Response{}
}
