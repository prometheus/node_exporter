package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	s, _ := ioutil.ReadFile("./collector/fixtures/megacli_disks.txt")
	fmt.Println(string(s))
}
