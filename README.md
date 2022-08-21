## 搭配Gin使用示例

```
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/with-zeal/restart"
)

func main() {
    router := gin.Default()
    router.GET("/hello", func(ctx *gin.Context) {
        ctx.JSON(http.StatusOK, gin.H{
            "msg": "Hello",
        })
    })
    // 旧服务处理请求最大时间，默认为1分钟
    restart.TimeOut = 2 * time.Minute
    // 延时任务，所以调用位置决定了后续定义的任务会被推迟
    // restart.StartAt(time.Now().Add(5 * time.Second))
    restart.WaitFor(5 * time.Second)
    // 第二个参数表示是否异步执行该函数
    restart.AddBeforeTask(func() {
        fmt.Println("服务启动前")
    }, false)
    restart.AddAfterTask(func() {
        fmt.Println("服务重启前")
    })
    server := restart.NewServer(router, ":8080")
    // err := server.RunTLS("./cert.pem","./key.pem")
    err := server.Run()
    if err != nil {
        log.Fatal(err)
    }
}

```
