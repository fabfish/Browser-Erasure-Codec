# 使用说明

## 内容和用法

最新内容在 Go 文件夹中

打开 Go 文件夹

scripts 文件夹里面是曾经用过的 JavaScript 函数，新的里面基本没用到。

wasm 文件夹里面是 Go-Webassembly 的相关文件，都要用到。

myerasure-go-wasm.html 是使用的网页。

需要用到一个 Web Server，比如说我是在 Go 文件夹下跑

```bash
goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'
```

然后在 :8080 打开 html 网页才能正常运作

## 关于代码

```javascript
    <script src="wasm/wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("wasm/mycoder.wasm"), go.importObject)
            .then((result) => go.run(result.instance));
    </script>
```

这一段会找到 wasm 文件，里面写好的函数在 js 里用。其他的基本没啥变化。

## 备用

GOOS=js GOARCH=wasm go build -o mycoder.wasm mycoder.go
goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'