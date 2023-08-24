package storage

import (
	"fmt"
)

var (
	ErrDateBusy = fmt.Errorf("this date was busy other event")
)
