
# 简单的序列id只用nodeid区分的唯一id, 有序并基于数据戳
## 算法原理
该算法实现一个简单的有序数字ID生成器,主要依赖节点ID和时间戳。每个ID生成节点被分配一个节点ID(nodeid),节点ID范围在1-1023之间。每个ID在生成时会记录一个时间戳(timestamp)。算法通过节点ID和时间戳算出一个序号(sequence),并与节点ID拼接生成最终ID。
该算法可以保证:
1. ID全局唯一:依赖节点ID和时间戳唯一性
2. ID有序:后生成的ID序号必定大于先生成的ID
3. 高性能:没有依赖数据库等外部资源,纯内存计算
## 实现步骤
1. 设置节点ID(nodeid),范围1-1023
2. 获取当前时间戳(timestamp)
3. 计算序号(sequence):
   - sequence = (timestamp - 节点启动时间戳) / 时间戳步长
   - 时间戳步长建议设置为1ms
4. 生成ID:ID = 节点ID * 1024 + sequence

 
## 使用示例
```golang
package main

import (
	"log"

	"github.com/474420502/orderedid"
)

func main() {
	var id orderedid.OrderedID
	creator := orderedid.New() // orderedid.NewWith(1)
	id = creator.Create()
	log.Println(id.Uint64(), id.Base58(), id.Timestamp(), id.NodeID()) // 142125288653825 27noD5f5R 1646623082475 0
}
```

# 性能测试
```go
func BenchmarkCase(b *testing.B) {
	var id OrderedID
	creator := New()
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