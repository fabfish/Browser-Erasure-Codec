//+build ignore

// Copyright 2015, Klaus Post, see LICENSE for details.
//
// Simple encoder example
//
// The encoder encodes a simgle file into a number of shards
// To reverse the process see "simpledecoder.go"
//
// To build an executable use:
//
// go build simple-decoder.go
//
// Simple Encoder/Decoder Shortcomings:
// * If the file size of the input isn't divisible by the number of data shards
//   the output will contain extra zeroes
//
// * If the shard numbers isn't the same for the decoder as in the
//   encoder, invalid output will be generated.
//
// * If values have changed in a shard, it cannot be reconstructed.
//
// * If two shards have been swapped, reconstruction will always fail.
//   You need to supply the shards in the same order as they were given to you.
//
// The solution for this is to save a metadata file containing:
//
// * File size.
// * The number of data/parity shards.
// * HASH of each shard.
// * Order of the shards.
//
// If you save these properties, you should abe able to detect file corruption
// in a shard and be able to reconstruct your data if you have the needed number of shards left.

package main

import (
	"flag"
	"fmt"
	//"io/ioutil"
	"os"
	//"path/filepath"

	"syscall/js"
	"github.com/klauspost/reedsolomon"
)

var dataShards = flag.Int("data", 4, "Number of shards to split the data into, must be below 257.")
var parShards = flag.Int("par", 2, "Number of parity shards")
var outDir = flag.String("out", "", "Alternative output directory")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  simple-encoder [-flags] filename.ext\n\n")
		fmt.Fprintf(os.Stderr, "Valid flags:\n")
		flag.PrintDefaults()
	}
}

//sendFragments((str)fileName,(str)fileType,(int)numOfDivision,(int)numOfAppend,(byte[][])content(content),(string[])digest,(int)fileSize);
func goEncoder(raw []byte, numOfDivision int, numOfAppend int)(content [][]byte){
//func goEncoder(rawJS js.Value, numOfDivisionJS js.Value, numOfAppendJS js.Value)(content [][]byte){
	/*
	var raw = new([]byte)
	var numOfDivision = new(int)
	var numOfAppend = new(int)
	raw = rawJS.Int()
	numOfDivision = numOfDivisionJS.Int()
	numOfAppend = numOfAppendJS.Int()
	*/
	//fmt.Println(raw)
	enc, err := reedsolomon.New(numOfDivision, numOfAppend)
	checkErr(err)
	content, err = enc.Split(raw)
	checkErr(err)
	err = enc.Encode(content)
	checkErr(err)
	//fmt.Println(content)
	//fmt.Println("hi")
	return content
}

func callEncoder(this js.Value, args []js.Value) interface{}{
	buffer := make([]byte, args[0].Length())
	//fmt.Println(args[0].Length)
	//fmt.Println(args[0])
	js.CopyBytesToGo(buffer, args[0])
	content := goEncoder(buffer, args[1].Int(), args[2].Int())
	fmt.Println("content len = ",len(content)," content[0] len =", len(content[0]))
	//fmt.Println("content = ",content)

	jsBuffer := make([]js.Value, len(content))
	jsInterface := make([]interface{},len(content))
	for  i:=0; i<len(content); i++{
		jsBuffer[i] =  js.Global().Get("Uint8Array").New(len(content[0]))
		js.CopyBytesToJS(jsBuffer[i], content[i])
		jsInterface[i] = js.ValueOf(jsBuffer[i])
	}
	//fmt.Println(jsBuffer)
	return js.ValueOf(jsInterface)

} 

//func mydecoder(content [][]byte)()

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("callEncoder",js.FuncOf(callEncoder))
	<-c
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}

