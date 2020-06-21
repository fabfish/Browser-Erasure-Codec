// Copyright 2015, Klaus Post, see LICENSE for details.
// Copyright 2020 XD
package main

import (
	"fmt"
	"os"
	"bytes"
	//"reflect"
	"crypto/md5"

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
	// use buffer to receive the args[0] (js content)
	buffer := make([]byte, args[0].Length())
	js.CopyBytesToGo(buffer, args[0])
	content := goEncoder(buffer, args[1].Int(), args[2].Int())
	// trans content into js value
	jsBuffer := make([]js.Value, len(content))
	jsInterface := make([]interface{},len(content))
	for  i:=0; i<len(content); i++{
		jsBuffer[i] =  js.Global().Get("Uint8Array").New(len(content[0]))
		js.CopyBytesToJS(jsBuffer[i], content[i])
		jsInterface[i] = js.ValueOf(jsBuffer[i])
	}
	return js.ValueOf(jsInterface)
} 


//decodeFile(fileName, fileType, numOfDivision, numOfAppend, content, digest, fileSize);
func goDecoder(shards [][]byte, numOfDivision int, numOfAppend int)(content []byte){
	enc, err := reedsolomon.New(numOfDivision, numOfAppend)
	checkErr(err)
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
	content = bytes.Join(shards,[]byte(""))
	return content
}

func callDecoder(this js.Value, args []js.Value) interface{}{
	//var decoded = erasure.recombine(content, fileSize, numOfDivision, numOfAppend);
	// use buffer and the .Index(i int) func to index args[0]
	buffer := make([][]byte, args[0].Length())
	for i:=0; i<len(buffer); i++ {
		buffer[i] = make([]byte, args[0].Index(0).Length())
		js.CopyBytesToGo(buffer[i], args[0].Index(i))
	}
	content := goDecoder(buffer, args[1].Int(), args[2].Int())
	//fmt.Println(content)
	jsBuffer :=  js.Global().Get("Uint8Array").New(len(content))
	js.CopyBytesToJS(jsBuffer, content)
	return js.ValueOf(jsBuffer)
} 

func callMd5(this js.Value, args []js.Value) interface{} {
	// 声明一个和文件大小一样的切片
	buffer := make([]byte, args[0].Length())
	// 将文件的bytes数据复制到切片中，这里传进来的是一个Uint8Array类型
	js.CopyBytesToGo(buffer, args[0])
	// 计算md5的值
	res := md5.Sum(buffer)
	// 调用js端的方法，将结构返回给js端
	// js.Global().Get("target").Get("callback").Invoke(fmt.Sprintf("%d", res))
	// return nil
	return fmt.Sprintf("%x", res)
}

func main() {
	// 声明一个函数，用来导出到js端，供js端调用
	calcMd5 := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// 声明一个和文件大小一样的切片
		buffer := make([]byte, args[0].Length())
		// 将文件的bytes数据复制到切片中，这里传进来的是一个Uint8Array类型
		js.CopyBytesToGo(buffer, args[0])
		// 计算md5的值
		res := md5.Sum(buffer)
		// 调用js端的方法，将结构返回给js端
		// js.Global().Get("target").Get("callback").Invoke(fmt.Sprintf("%d", res))
		// return nil
		return fmt.Sprintf("%x", res)
	})
	c := make(chan struct{}, 0)
	js.Global().Get("target").Set("calcMd5", calcMd5)
	js.Global().Set("callMd5",js.FuncOf(callMd5))
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