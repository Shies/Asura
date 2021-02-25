package dao

import (
	"encoding/gob"
	"bytes"
	"fmt"
)

type Dummy interface {
	Foobar()
}

type Shies struct{
	// to do sth.
}

type Danko struct{
	// to do sth.
}

func (s *Shies) Foobar() {
	fmt.Println("Shies")
}

func (d *Danko) Foobar() {
	fmt.Println("Danko")
}

func init() {
	// This type must match exactly what youre going to be using,
	// down to whether or not its a pointer
	gob.Register(&Shies{})
	gob.Register(&Danko{})
}

func handle() {
	network := new(bytes.Buffer)
	enc := gob.NewEncoder(network)

	var inter Dummy
	inter = new(Shies)

	// Note: pointer to the interface
	err := enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	inter = new(Danko)
	err = enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	// Now lets get them back out
	dec := gob.NewDecoder(network)

	var get Dummy
	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}

	// Should meow
	get.Foobar()

	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}

	// Should woof
	get.Foobar()

}