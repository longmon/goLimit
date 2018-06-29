# goLimit
漏桶算法，令牌桶算法实现

## 漏桶算法应用示例
```go
leaky := goLimit.NewLeaking(100) // 100 QPS
tm, err := leaky.Take() //获取授权，马上返回不会阻塞！获取失败统一返回nil,errors.New("rate limited")
tm := leaky.Wait() //阻塞等待，直到获得一个授权
```

## 令牌桶算法应用示例
```go
bucket := goLimit.NewBucket(1000, 100) //桶容量1000， 每秒往桶投放100个令牌
t, err := bucket.Take() //从桶里取出一个令牌，如果桶为空，不阻塞等待，返回nil,并检查错误
t := bucket.Wait() // 从桶里取出一个令牌，如果桶为空，阻塞等待，直到桶里面有令牌可用为止
```