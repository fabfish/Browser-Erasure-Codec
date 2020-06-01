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

/*字符串转化成ArrayBuffer
 * maybe useful when providing API to server
 */
function str2ab(str) {
    var buf = new ArrayBuffer(str.length); // 每个字符占用1个字节
    var bufView = new Uint8Array(buf);
    for (var i = 0, strLen = str.length; i < strLen; i++) {
        bufView[i] = str.charCodeAt(i);
    }
    return buf;
}