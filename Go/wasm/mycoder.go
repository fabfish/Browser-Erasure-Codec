// Copyright 2015, Klaus Post, see LICENSE for details.
// Copyright 2020 OSH-xdontpanic
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
	// use []interface{} directly. the arg (if array) for js.ValueOf() must be []interface{}
	jsContent := make([]interface{},len(content))
	for i:=0; i<len(content); i++{
		jsContent[i] = js.Global().Get("Uint8Array").New(len(content[0]))
		js.CopyBytesToJS(jsContent[i].(js.Value),content[i])
	}
	return js.ValueOf(jsContent)
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
			//os.Exit(1)
		}
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted.")
			//os.Exit(1)
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
		// if args[0][i]==null, set buffer[i] as nil.
		if !args[0].Index(i).Equal(js.Null()) {
			buffer[i] = make([]byte, args[0].Index(i).Length())
			js.CopyBytesToGo(buffer[i], args[0].Index(i))
		}else {
			buffer[i]=nil;
		}
	}
	content := goDecoder(buffer, args[1].Int(), args[2].Int())
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
	// 调用js端的方法，将字符串返回给js端
	return fmt.Sprintf("%x", res)
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("callMd5",js.FuncOf(callMd5))
	js.Global().Set("callEncoder",js.FuncOf(callEncoder))
	js.Global().Set("callDecoder",js.FuncOf(callDecoder))
	<-c
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		//os.Exit(2)
	}
}