package main
import (
 "fmt"
 "time"
 "strings"
 "os"
 "io"
 "net/http"
)
type Pool struct {
    Queue  chan func() error
    Number int
    Total  int
    result chan error
    finishCallback func()
}
func (p *Pool) Init(number int, total int) {// 初始化
    p.Queue  = make(chan func() error, total)
    p.Number = number
    p.Total  = total
    p.result = make(chan error, total)
}
func (p *Pool) Start() {// 开启Number个goroutine
    for i := 0; i < p.Number; i++ {
        go func() {
            for {
                task, ok := <-p.Queue
                if !ok {
                        break
                }
                err := task()
                p.result <- err
            }
        }()
    }
    for j := 0; j < p.Total; j++ {// 获得每个work的执行结果
        res, ok := <-p.result
        if !ok {
                break
        }
        if res != nil {
                fmt.Println(res)
        }
    }    
    if p.finishCallback != nil {// 所有任务都执行完成，回调函数
            p.finishCallback()
    }
}
func (p *Pool) Stop() {// 关
        close(p.Queue)
        close(p.result)
}
func (p *Pool) AddTask(task func() error) { // 添加任务
        p.Queue <- task
}
func (p *Pool) SetFinishCallback(callback func()) {// 设置结束回调
        p.finishCallback = callback
}
func Download_test() {
    urls := []string{
        "http://it.facesoho.com/demos/index.html",
        "http://it.facesoho.com/demos/index_2.html",
        "http://it.facesoho.com/demos/index_3.html",
        "http://it.facesoho.com/demos/index_4.html",
        "http://it.facesoho.com/demos/index_5.html",
    }
    pool := new(Pool)
    pool.Init(3, len(urls))
    for i := range urls {
        url := urls[i]
        pool.AddTask(func() error {
                return download(url)
        })
    }  
    isFinish := false
    pool.SetFinishCallback(func() {
        func(isFinish *bool) {
                *isFinish = true
        }(&isFinish)
    })
    pool.Start()  
    for !isFinish {
            time.Sleep(time.Millisecond * 100)
    }  
    pool.Stop()
    fmt.Println("所有操作已完成！")
}
func download(url string) error {
        fmt.Println("开始下载... ", url)  
        sp := strings.Split(url, "=")
        filename := sp[len(sp)-1]  
        file, err := os.Create("./download/" + filename)
        if err != nil {
                return err
        }  
        res, err := http.Get(url)
        if err != nil {
                return err
        }  
        length, err := io.Copy(file, res.Body)
        if err != nil {
                return err
        }  
        fmt.Println("## 下载完成！ ", url, " 文件长度：", length)
        return nil
}
func main() { //主函数
  Download_test()
}