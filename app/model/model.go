package model

type Welcome struct {
	Hello string
	World string
}

func (*Welcome) TableName() string {
	return ""
}

func (*Welcome) Len() int {
	return 0
}
