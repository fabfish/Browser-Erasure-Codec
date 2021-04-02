0402
又回来工作力
现在的目标是对文件进行分片，按照指令的512KB左右分成许多小块
现状是文件会被放进一个worker里，直接传给wasm
今天考虑改变worker实现上述目标，有两种办法，第一种是进来的raw先分成片，然后for循环放进worker
第二种是在worker里面创建subworker
我认为创建subworker虽然难搞一点，但是有利于未来改进。比如说你有多个文件输入，理想的情况是用两个worker
添加了用户参数512（每个文件512）

还有一个问题，js内容在哪里分块？
go的库虽然有优化，但是编译成wasm之后就没用了
所以就在js里面做

把文件分片：
https://stackoverflow.com/questions/26224597/how-to-parse-uint8array-into-object
现在学习worker怎么用
