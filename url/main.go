package main

import (
	"fmt"
	"net/url"
)

func main() {
	c := url.Values{"method": {"get"}, "id": {"1"}}
	c.Add("phone", "v1qUhIZ/SDpA/4uLa7ReKW+w==")
	fmt.Println(c.Encode())
	fmt.Println(c.Get("method"))

	c.Set("method", "post")
	fmt.Println(c.Encode())
	fmt.Println(c.Get("method"))

	c.Del("method")
	fmt.Println(c.Encode())
	fmt.Println(c.Get("method"))

	c.Add("new", "hi")
	fmt.Println(c.Encode())
	fmt.Println(c.Encode())
}
