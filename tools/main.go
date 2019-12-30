package main

import (
	"fmt"
	"reflect"
)

type Em string

const (
	a1 Em = "1"
	a2 Em = "2"
)

// func (e Em) Valid(s string) {
// 	t := Em(s)

// }

func appedMap(a map[string]string, c string) {
	a[c] = c
}

func main() {
	ss := Em("1ss")
	vs, ok := reflect.ValueOf(ss).Interface().(Em)
	if ok {
		fmt.Println("ok", vs)
	} else {
		fmt.Println("not")
	}
	reflect.TypeOf(ss)
	fmt.Println(Em(ss), reflect.TypeOf(Em(ss)))
	var errs []error
	fmt.Println(len(errs))

	cc := make(map[string]string)
	appedMap(cc, "sss")
	fmt.Println(cc)
}
