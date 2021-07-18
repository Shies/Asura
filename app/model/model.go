package model

import _ "sort"

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

func (*Welcome) Less(i, j int) bool {
	return false
}

func (*Welcome) Swap(i, j int) {
	return
}
