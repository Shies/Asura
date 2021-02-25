package model

type Test struct {
	ID string	`json:"id"`
	SN string	`json:"sn"`
}

func (*Test) TableName() string {
	return "test"
}