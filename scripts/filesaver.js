//invalid
$("#export").click(function () {
    var content = "����ֱ��ʹ��HTML5���е�����";
var blob = new Blob([content], { type: "text/plain;charset=utf-8" });
saveAs(blob, "file.txt");//saveAs(blob,filename)
});