# ErasureCodeforBrowserSide

余致远

## 介绍

懒得写

## 使用

懒得写

## 现有问题

懒得写

## 日记

##### 6.1

儿童节快乐

今天把仓库搬上 GitHub 不过忘了设置什么来防止大家发现（被发现了就要被督促进度了！！

然后建仓库之后方便在虚拟机上用 快乐！

##### 6.3

写点啥

本次的工作是写好了收文件、计算 hash 然后传给后端这一部分。还没有做从后端接收。

loli 给的接口：
object encodeFile();
void decodeFile((str)fileName,(str)fileType,(int)numOfDivision,(int)numOfAppend,(byte[][])content,(string[])digest);

目前做的接口则是

##### 6.6

做完了 lab4，熬个夜

目前写完了 loli 要求的接口，但还有很多可以优化：

1. 输入数字的方式有点傻，需要按按钮
2. 把东西放进 *.js 里面应该会好看一点，等会儿试试
3. 还是 xk 老师说的，要注意分块大小。
4. Hash 的方式可以再探索一下。

然后应该就是开始做 wasm 了，明天可以试试。

##### 6.20

这段时间的日志断了，但是现在已经实现了在浏览器用 go wasm。感觉 100M 的文件差不多是性能上限。

##### 6.21 

考虑了一下，用 go 读文件没问题，而且说不定能省略一点步骤

https://juejin.im/post/5eb2191df265da7bbf21a0f4

## TODO 日志

TODO: 05 / 27
a. 和 loli 对接，大概问题有
格式：编解码都要 uint8 ，这个比较简单。
附加的信息：
    错误校验：需要哈希，可以放后面
    文件名/大小/块长 等等

b. 本周工作
优先级
! 寻找外包
1 写个效果测评
1 擦除数据（模拟传输时的丢失/出错）改成随机的
2 把当前html里的函数放进scripts文件夹中
2 分块大小控制：xk老师建议控制在512KB几百块左右（其实忘了），回头看看会议记录他说了啥，然后试试
3 试用wasm，毕竟目前这个太慢。（稍微搞懂了一点js，可能会方便一些。

界面和组织很丑陋

现在最新的是 myera4.html

erasure code 主要脚本在 scripts 文件夹中，测试主要写在 html 里还没移走

0530
生成不重复随机数
https://www.cnblogs.com/fishyao/p/RandomNumber.html

0601
TODO
1 按 loli 的接口重写
应当明确 js 性能不好，不能满足我们的期望，所以使用 wasm 势在必行。
问题是现有的 go 实现也有一定问题，即文本最后会有多余的 0 字符，应该想办法清除这些字符。
2 给新虚拟机配置好 emscripten 和 golang 等工具，写一下配置过程。
2 按接口改写 go wasm，查阅 stream 和 file 什么的关系。
2 查看 flutter 相关，准备修改 UI。

sendFragments((str)fileName,(str)fileType,(int)numOfDivision,(int)numOfAppend,(byte[][])content,(string[])digest);
recvFragments((str)fileName,(str)fileType,(int)numOfDivision,(int)numOfAppend,(byte[][])content,(string[])digest);

0603
TODO
想到一个问题：js 脚本里是如何自动忽略结尾的 0 的？
（补充）观察过了，是用原文件的 fileSize 和 decode 生成的数组元素个数比较，用 .subarray 方法消去末尾的 0.

0606
TODO
见日记。

0620
TODO
增加 md5 功能，考虑在 go 里面增加按钮上传文件的事件。需要看一下那个 md5 的示范。