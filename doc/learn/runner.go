package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"time"
)

type Runner struct {
	interrupt chan os.Signal   //接收信号的通道
	complete  chan error       //接收完成的通道
	timeout   <-chan time.Time //接收超时的通道
	tasks     []func(int)      //任务数组，保存的是函数，函数作为一个元素使用。
}

//统一错误
var ErrTimeout = errors.New("received timeout")
var ErrInterrupt = errors.New("received interrupt")
//
func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1), //创建信号通道
		complete:  make(chan error),        //任务完成通道
		timeout:   time.After(d),           //返回一个时间通道，该通道在无操作时，d时间后激发超时
	}
}
func (r *Runner) Add(tasks ...func(int)) { //添加任务，可以多参数方式，可以添加多个。
	r.tasks = append(r.tasks, tasks...) //切片可以这样使用。数组不可以使用
}
func (r *Runner) Start() error {
	// Notify函数让signal包将输入信号转发到r.interrupt。
	// 如果没有列出要传递的信号，会将所有输入信号传递到r.interrupt；
	// 否则只传递列出的输入信号。
	// signal包不会为了向c发送信息而阻塞（就是说如果发送时c阻塞了，signal包会直接放弃）：
	// 调用者应该保证c有足够的缓存空间可以跟上期望的信号频率。
	// 对使用单一信号用于通知的通道，缓存为1就足够了。
	signal.Notify(r.interrupt, os.Interrupt) //信号接收通知

	// 用不同的goroutine执行不同的任务, 并将结果放到 complete 通道.
	go func() {
		r.complete <- r.run()
	}()

	// 如果complete通道接收到消息，就继续运行。
	// 如果complete通道没有接收到消息就阻塞。
	// 此时主线程就阻塞，等待信号到来。
	select {
	// 当任务处理完成时发出的信号
	case err := <-r.complete:
		return err
	case <-r.timeout: // 当任务处理程序运行超时时发出的信号
		return ErrTimeout
	}
}
func (r *Runner) run() error {
	for id, task := range r.tasks { //遍历数组，运行各个函数
		// 检测操作系统的中断信号
		if r.gotInterrupt() {
			return ErrInterrupt //接收到CTRL-C信号
		}
		task(id)
	}
	return nil //如果无误，则返回nil
}
func (r *Runner) gotInterrupt() bool {
	select {
	case <-r.interrupt:          //接收到os.Interrupt
		signal.Stop(r.interrupt) //停止接收以后的信号
		return true

	default:// 继续正常运行
		return false
	}
}

//================================================================================

const timeout = 5 * time.Second

func main() {

	log.Println("Starting work.")
	r := New(timeout)
	r.Add(
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask())
	if err := r.Start(); err != nil {
		switch err {
		case ErrTimeout:
			log.Println("Terminating due to timeout.")
			os.Exit(1)
		case ErrInterrupt:
			log.Println("Terminating due to interrupt.")
			os.Exit(2)
		}
	}
	log.Println("Process ended.")
}

func createTask() func(int) {
	return func(id int) {
		log.Printf("Processor - Task #%d.", id)
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}
