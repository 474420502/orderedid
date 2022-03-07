# 简单的序列id只用nodeid区分的唯一id. 有序. 基于数据戳.
# 多线程也需要提供 nodeid < 1024的id. 特点就是序列和速度.

```golang
package main

import (
	"log"

	"github.com/474420502/orderedid"
)

func main() {
	var id orderedid.OrderedID
	creator := orderedid.New(1)
	id = creator.Create()
	log.Println(id.Uint64(), id.Base58(), id.Timestamp(), id.NodeID()) // 142125288653825 27noD5f5R 1646623082475 1
}

```

# 性能测试
```go
func BenchmarkCase(b *testing.B) {
	var id OrderedID
	creator := New(1)
	for i := 0; i < b.N; i++ {
		id = creator.Create()
	}

	b.Log(id)
}
```


```shell
goos: linux
goarch: amd64
pkg: github.com/474420502/orderedid
cpu: AMD Ryzen 7 5700G with Radeon Graphics         
BenchmarkCaseMID-16    	24551696	        42.77 ns/op	       0 B/op	       0 allocs/op
```