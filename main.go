package main

import "github.com/jbdoumenjou/go-sandbox-concurency/confinement"

func main() {
	confinement.Run([]int{1, 2, 3, 4, 5})
}
