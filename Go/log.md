现在完成的内容：编写成了 wasm
主要的坑：go get reedsolomon
现在要做：为了避免什么CORS，需要自己建服务器
npm install -g simplehttpserver
https://medium.com/starbugs/run-golang-on-browser-using-wasm-c0db53d89775

go get -u github.com/shurcooL/goexec

<script src="static/wasm_exec.js"></script>
<script>
	const go = new Go();
	WebAssembly.instantiateStreaming(fetch("static/main.wasm"), go.importObject)
		.then((result) => go.run(result.instance));
</script>

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./examples
GOOS=js GOARCH=wasm go build -o myencoder.wasm myencoder.go
goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'

目前进度：小文件有输出，大文件直接崩溃
输出是什么？我也想知道 总之编码应该成功了
