package utils

import (
	"fmt"
	"reflect"
)

func PrintType(v interface{}) {
	fmt.Printf("%v", reflect.TypeOf(v))
}

func Print(v interface{}) {
	fmt.Printf("%v", v)
}
