/* Simplified way to detect file change,
 * a better way is to use Hash key
 */
function arraysEqual(a, b) {
    if (a.length != b.length)
        return false;

    for (var i = 0; i < a.length; i++)
        if (a[i] != b[i])
            return false;
    return true;
}

/*ArrayBuffer
 * maybe useful when providing API to server
 */
function str2ab(str) {
    var buf = new ArrayBuffer(str.length); // 
    var bufView = new Uint8Array(buf);
    for (var i = 0, strLen = str.length; i < strLen; i++) {
        bufView[i] = str.charCodeAt(i);
    }
    return buf;
}

function c() {
    var t = document.getElementById("txt");
    t.value = "ÎÒºÜºÃ£¡";
}

//0530 add random
function randomize(n) {
    for (var j = 0; j < n; j++) {
        randomArray[j] = j;
    }
    randomArray.sort(function () {
        return 0.5 - Math.random();
    });
    //0530show random
    if (document.getElementById("random") != null) {
        for (var j = 0; j < n; j++) {
            document.getElementById("random").innerHTML += (randomArray[j] + " , ");
        }
        document.getElementById("random").innerHTML += "</br>";
    }
}