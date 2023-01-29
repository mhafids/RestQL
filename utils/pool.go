package utils

import (
	"bytes"
	"sync"
)

var Bufpool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}
