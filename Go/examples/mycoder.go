// Copyright 2015, Klaus Post, see LICENSE for details.
// Copyright 2020 XD
////+build ignore0689bgh897j6uy
package main

import (
	//"flag"
	"fmt"
	//"io/ioutil"
	"os"
	//"path/filepath"
	"bytes"
	//"reflect"

	"syscall/js"
	"github.com/klauspost/reedsolomon"
)

//sendFragments((str)fileName,(str)fileType,(int)numOfDivision,(int)numOfAppend,(byte[][])content(content),(string[])digest,(int)fileSize);
func goEncoder(raw []byte, numOfDivision int, numOfAppend int)(content [][]byte){
	enc, err := reedsolomon.New(numOfDivision, numOfAppend)
	checkErr(err)
	content, err = enc.Split(raw)
	checkErr(err)
	err = enc.Encode(content)
	checkErr(err)
	return content
}

func callEncoder(this js.Value, args []js.Value) interface{}{
	buffer := make([]byte, args[0].Length())
	js.CopyBytesToGo(buffer, args[0])
	content := goEncoder(buffer, args[1].Int(), args[2].Int())
	//fmt.Println("content len = ",len(content)," content[0] len =", len(content[0]))
	jsBuffer := make([]js.Value, len(content))
	jsInterface := make([]interface{},len(content))
	for  i:=0; i<len(content); i++{
		jsBuffer[i] =  js.Global().Get("Uint8Array").New(len(content[0]))
		js.CopyBytesToJS(jsBuffer[i], content[i])
		jsInterface[i] = js.ValueOf(jsBuffer[i])
	}
	return js.ValueOf(jsInterface)
	//return js.ValueOf([]interface{}(jsBuffer))
} 


//decodeFile(fileName, fileType, numOfDivision, numOfAppend, content, digest, fileSize);
func goDecoder(shards [][]byte, numOfDivision int, numOfAppend int)(content []byte){
	enc, err := reedsolomon.New(numOfDivision, numOfAppend)
	checkErr(err)
	//fmt.Println("shards = ",shards)
	//fmt.Println(numOfDivision, numOfAppend)
	// Verify the shards
	ok, err := enc.Verify(shards)
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data")
		err = enc.Reconstruct(shards)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			os.Exit(1)
		}
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted.")
			os.Exit(1)
		}
		checkErr(err)
	}
	//fmt.Println(shards);
	content = bytes.Join(shards,[]byte(""))
	return content
}

func callDecoder(this js.Value, args []js.Value) interface{}{
	//var decoded = erasure.recombine(content, fileSize, numOfDivision, numOfAppend);
	//fmt.Println(reflect.TypeOf(args[0].Index(0).Index(0)))
	//fmt.Println(args[0].Index(0).Length(), args[0].Index(0).Index(0))
	buffer := make([][]byte, args[0].Length())
	for i:=0; i<len(buffer); i++ {
		buffer[i] = make([]byte, args[0].Index(0).Length())
		js.CopyBytesToGo(buffer[i], args[0].Index(i))
	}
	//fmt.Println("buffer = ",buffer)
	content := goDecoder(buffer, args[1].Int(), args[2].Int())
	//fmt.Println(content)
	jsBuffer :=  js.Global().Get("Uint8Array").New(len(content))
	js.CopyBytesToJS(jsBuffer, content)
	return js.ValueOf(jsBuffer)
	//js.CopyBytesToGo(buffer, args[0])
	/*
	bfrags := make([][]js.Value,args[0].Length())
	bfrags = args[0].JSValue
	buffer := make([][]byte, args[0].Length())
	for i:=0; i<len(buffer); i++{
		js.CopyBytesToGo(buffer[i], args[0][i])
	}
	content := goDecoder(buffer, args[2].Int(), args[3].Int())
	fmt.Println("content len = ",len(content)," content[0] len =", len(content[0]))
	jsBuffer :=  js.Global().Get("Uint8Array").New(len(content))
	js.CopyBytesToJS(jsBuffer, content)
	return js.ValueOf(jsBuffer)
	*/
	/*
	var x js.Wrapper
	x = args[0]
	fmt.Println("x = ", args[0])
	content ,err := interface{}(x.JSValue()).(js.Value)
	fmt.Println(content,err, reflect.TypeOf(content), args[0].Length())
	fmt.Println(content.Index(0))
	buffer :=make([]byte, args[0].Length())
	js.CopyBytesToGo(buffer, content)
	fmt.Println(buffer)
	return js.ValueOf(args[0])
	*/
} 

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("callEncoder",js.FuncOf(callEncoder))
	js.Global().Set("callDecoder",js.FuncOf(callDecoder))
	<-c
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}

