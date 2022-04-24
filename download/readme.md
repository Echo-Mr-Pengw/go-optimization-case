# 背景
前端上传Excel至oss，返回给予后端Excel的下载地址；后端通过Excel的地址下载Excel文件至本地，导出数据进行相关操作。
# 技术方案选型
1.  一次性把数据读入memory
2. 按照一定大小把数据读到memory

### 一次性把数据读入memory
#### 代码示例
```
func download() {

	file := "./download.txt"
	f, err  := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("openfile err:", err)
		return
	}

	defer func() {
		_ = f.Close()
	}()
	
	resp, err := http.Get("https://xxxx/31650620292790.xlsx")
	if err != nil {
		fmt.Println("get err=", err)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body) // 一次性把数据读记载到内存
	_, err = f.Write(body)
	fmt.Println("write ", err)

}
```
#### pprof分析
对内存进行分析
```
➜  ~ go tool pprof http://127.0.0.1:6161/debug/pprof/heap
Fetching profile over HTTP from http://127.0.0.1:6161/debug/pprof/heap
Saved profile in /Users/pengwei/pprof/pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.053.pb.gz
Type: inuse_space
Time: Apr 24, 2022 at 10:10am (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 3119.25kB, 100% of 3119.25kB total
Showing top 10 nodes out of 30
      flat  flat%   sum%        cum   cum%
 2086.21kB 66.88% 66.88%  2086.21kB 66.88%  bytes.makeSlice
  520.04kB 16.67% 83.55%   520.04kB 16.67%  net/http.glob..func5
     513kB 16.45%   100%      513kB 16.45%  vendor/golang.org/x/net/http2/hpack.newInternalNode (inline)
         0     0%   100%  2086.21kB 66.88%  bytes.(*Buffer).ReadFrom
         0     0%   100%  2086.21kB 66.88%  bytes.(*Buffer).grow
         0     0%   100%  2086.21kB 66.88%  io/ioutil.ReadAll (inline)
         0     0%   100%  2086.21kB 66.88%  io/ioutil.readAll
         0     0%   100%  2086.21kB 66.88%  main.download
         0     0%   100%  2086.21kB 66.88%  main.main
         0     0%   100%  1033.04kB 33.12%  net/http.(*http2ClientConn).readLoop
(pprof)
```
通过上面对内存的分析结构可以看出，占用memory较多的函数main.download、realAll等。通过list追踪上面函数，最终定位到最终耗时的函数为：bytes.makeSlice
```
(pprof) list bytes.makeSlice
Total: 3.05MB
ROUTINE ======================== bytes.makeSlice in /Users/pengwei/.g/go/src/bytes/buffer.go
    2.04MB     2.04MB (flat, cum) 66.88% of Total
         .          .    224:	defer func() {
         .          .    225:		if recover() != nil {
         .          .    226:			panic(ErrTooLarge)
         .          .    227:		}
         .          .    228:	}()
    2.04MB     2.04MB    229:	return make([]byte, n)
         .          .    230:}
         .          .    231:
         .          .    232:// WriteTo writes data to w until the buffer is drained or an error occurs.
         .          .    233:// The return value n is the number of bytes written; it always fits into an
         .          .    234:// int, but it is int64 to match the io.WriterTo interface. Any error
(pprof)
```
最终定位到最终函数的操作是`make([]byte, n)`

### 分配读取
#### 代码示例
```
func download1() {
	file := "./download.txt"
	f, err  := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("openfile err:", err)
		return
	}

	defer func() {
		_ = f.Close()
	}()

	resp, err := http.Get("https://xxxx/31650620292790.xlsx")
	if err != nil {
		fmt.Println("get err=", err)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()


	_, err = io.Copy(f, resp.Body)  // 分批
	fmt.Println(err)
}
```
#### pprof分析
```
➜  ~ go tool pprof http://127.0.0.1:6161/debug/pprof/heap
Fetching profile over HTTP from http://127.0.0.1:6161/debug/pprof/heap
Saved profile in /Users/pengwei/pprof/pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.055.pb.gz
Type: inuse_space
Time: Apr 24, 2022 at 10:50am (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1042.17kB, 100% of 1042.17kB total
      flat  flat%   sum%        cum   cum%
  528.17kB 50.68% 50.68%   528.17kB 50.68%  io.copyBuffer
     514kB 49.32%   100%      514kB 49.32%  bufio.NewWriterSize (inline)
         0     0%   100%   528.17kB 50.68%  io.Copy (inline)
         0     0%   100%   528.17kB 50.68%  main.download1
         0     0%   100%   528.17kB 50.68%  main.main
         0     0%   100%      514kB 49.32%  net/http.(*conn).serve
         0     0%   100%      514kB 49.32%  net/http.newBufioWriterSize
         0     0%   100%   528.17kB 50.68%  runtime.main
(pprof)
```
占用内存比例明显下降，下载6.4M的excel，从占用内存`2086.21kB`降低到` 528.17kB`，差不多降低`4倍`

# 源码分析
todo



