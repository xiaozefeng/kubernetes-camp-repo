## 思考：当我们写代码的时候，函数的参数传递应该用 值还是指针？
从功能实现上来讲：需要修改数据的情况用指针接受者，不需修改数据的时候用值。
但是像 slice 和 map 就不应该 用指针接受者，因为他们已经被设计成指针使用了

当由于是用值还是指针时，用指针就对了。


## 1.1 practice:  ["I", "am", "stupid", "and", "weak"]  使用for 改成 ["I", "am", "smart", "and", "strong"]
```go
	var s = []string{"I", "am", "stupid", "and", "weak"}
	var mapping = map[string]string{
		"stupid": "smart",
		"weak":   "strong",
	}
	for i := 0; i < len(s); i++ {
		if val, ok := mapping[s[i]]; ok {
			s[i] = val
		}
	}
	fmt.Println(s)
```

## 1.2 practice 生产者消费者
```go
package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	var c = make(chan int, 10)
	timeout:= time.After(time.Second * 5)
	for {
		select {
		case val := <-c:
			log.Printf("received val from channel : %v\n", val)
		case <-time.After(time.Second):
			var x = rand.Intn(100)
			log.Printf("send val %d into channel", x)
			c <- x
		case <- timeout :
			log.Println("timeout !!! ")
			return 
		}
	}

}
```