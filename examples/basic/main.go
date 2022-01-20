package main

import (
	"fmt"

	bonesay "github.com/anthonycuervo23/bonesay/v2"
)

func main() {
	if false {
		simple()
	} else {
		complex()
	}
}

func simple() {
	say, err := bonesay.Say(
		"Hello",
		bonesay.Type("default"),
		bonesay.BallonWidth(40),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}

func complex() {
	bone, err := bonesay.New(
		bonesay.BallonWidth(40),
		//bonesay.Thinking(),
		bonesay.Random(),
	)
	if err != nil {
		panic(err)
	}
	say, err := bone.Say("Hello")
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}
