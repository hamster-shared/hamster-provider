package events

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEventType(t *testing.T) {

	var e1 EventInterface
	e1 = &StartVm{}

	fmt.Println(reflect.TypeOf(e1).String())
}
