//invalid
$("#export").click(function () {
    var content = "这是直接使用HTML5进行导出的";
var blob = new Blob([content], { type: "text/plain;charset=utf-8" });
saveAs(blob, "file.txt");//saveAs(blob,filename)
});