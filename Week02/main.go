package main

import (
	"fmt"
	"GoProject/Week02/Service"
)

func main() {
	value, err := Service.WrapError()
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Printf("%s", value)
}
