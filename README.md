# goLimit
漏桶算法，令牌桶算法实现

### 下载安装
```
go get -u github.com/longmon/goLimit

//或使用glide
glide get github.com/longmon/goLimit
```

### 漏桶算法应用示例
```go
import(
    "github.com/longmon/goLimit"
)

// 定义一个每秒钟漏出100个水滴的桶 
leaky := goLimit.NewLeaking(100) 

//获取授权，马上返回不会阻塞！获取失败统一返回`rate limited`错误
tm, err := leaky.Take() 

//阻塞等待，直到获得一个授权
tm := leaky.Wait() 
```

### 令牌桶算法应用示例
```go
import(
    "github.com/longmon/goLimit"
)
//初始化令牌桶，容量1000， 每秒往桶投放100个令牌, 取令牌时如果桶中无令牌超时10ms
bucket := goLimit.NewBucket(1000, 100, time.Millisecond * 10) 

//从桶里取出一个令牌，如果桶为空，不阻塞等待，返回nil,并检查错误
t, err := bucket.Take() 

// 从桶里取出一个令牌，如果桶为空，阻塞等待直到桶里面有令牌可用或超时
t := bucket.Wait()
```

### 蓄水池算法
```go
import(
    "github.com/longmon/goLimit"
)
//初始化容量1000的桶，设置超时10ms, 如果想一直阻塞，第二个参数设为0，有阻塞的池应该有放回的操作
p := goLimit.NewPool(1000, time.Millisecond * 10 ) 

// 初始化池中数据，也可以调用p.Put()一次放一个数据进去。
p.Init(interface{}{}) 

// 从池中取出数据，如果确定不会再把元素放入池中，就使用这个方法取池中数据，直到返回empty pool错误，池中再无数据可用
p.Take()   

//阻塞等待池中数据，如果池中没有数据就会阻塞，直到超时。这个方法应用于客户端处理完事务后再把数据放入池中的场景
p.Wait()  

//把数据放回池中。也可以再初始把的时候把不同的数据放入池中
p.Put(interface{}{}) 
```
