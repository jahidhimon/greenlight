package data

import (
	"fmt"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("\"%d mins\"", r)
	return []byte(jsonValue), nil
}
