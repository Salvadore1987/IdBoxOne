package lib

type PersonalInfo struct {
	Name string
	Surname string
	BirthDate string
	Country string
	Sex string
	ValidityDate string
	DocNumber string
	DocType string
	IssuingState string
	OptionalData string
}

func NewPersonalInfo() *PersonalInfo {
	return &PersonalInfo{}
}