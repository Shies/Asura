package model

type Welcome struct {
	Hello string
	World string
}

func (*Welcome) TableName() string {
	return "welcome"
}

func (*Welcome) Len() int {
	return 0
}

func (*Welcome) Swap() bool {
	return false
}
