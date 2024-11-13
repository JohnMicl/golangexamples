## 构建
```
go build -gcflags "-N -l" ./helloworld.go 
```

具体编译器标志
```
-N
含义：禁用优化。
作用：告诉编译器不要对代码进行优化。这对于调试非常有用，因为未优化的代码更容易与源代码对应，便于调试工具（如 GDB）正确地显示变量值和执行流程。
-l
含义：禁用内联。
作用：告诉编译器不要对函数进行内联展开。内联是一种优化技术，可以减少函数调用的开销，但会使生成的二进制文件更大，且调试时更难以跟踪函数调用。
```

## 创建一个空目录
 mkdir  /mnt/myfs

## 挂载运行
```
./helloworld --mountpoint=/mnt/myfs --fuse.debug=true
```

## 验证
```
df -aTh | grep myfs
helloworld     fuse.hellofs     0.0K  0.0K  0.0K     - /mnt/myfs
```
