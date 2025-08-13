package main

import (
	"errors"
	"io"
	"math/rand"
	"sync"
	"time"
)

// 池化技术（Pooling）是一种优化资源管理的重要手段，核心思想是提前创建并复用资源，避免频繁创建和
// 销毁资源带来的性能开销（如内存分配、GC 压力、系统调用等）。常见的池化场景包括 goroutine 池、连接池、对象池等

func main() {
	// case1()

	case2()
}

// Goroutine 池的作用：限制并发 goroutine 数量，复用 goroutine 处理多个任务。

// Worker 工作协程，从任务通道接收任务并执行
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// 模拟任务处理
		time.Sleep(time.Millisecond * 100)
		println("worker", id, "processed job", job)
	}
}

func case1() {
	const (
		numJobs    = 100 // 总任务数
		numWorkers = 10  // 工作协程数量（池大小）
	)

	jobs := make(chan int, numJobs)
	var wg sync.WaitGroup

	// 启动工作协程（初始化池）
	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, &wg)
	}

	// 发送任务到池
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // 关闭任务通道，通知worker任务结束

	wg.Wait() // 等待所有任务完成
}

// 对象池（复用重量级对象）
// 对于创建成本高的对象（如大型结构体、解析器、加密器），频繁创建会增加内存分配和 GC 负担。对象池通过缓存这些对象实现复用。

// 定义一个重量级对象（示例）
type HeavyObject struct {
	Data []byte
}

var pn = 0

var objectPool = sync.Pool{
	// 当池为空时，创建新对象的函数
	New: func() interface{} {
		pn++
		println("创建新对象", pn)
		return &HeavyObject{
			Data: make([]byte, 1024*1024), // 1MB的对象，模拟创建成本高
		}
	},
}

func case2() {

	rand.NewSource(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(100)

	// 并发获取和释放对象
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			// 从池获取对象
			obj := objectPool.Get().(*HeavyObject)

			// 使用对象（模拟操作）
			obj.Data[0] = 1

			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

			// 用完放回池（注意：放回前需重置对象状态，避免污染）
			obj.Data[0] = 0
			objectPool.Put(obj)
		}()
	}

	wg.Wait()
}

// 基于channel的通用连接池
// 实现原理 :将连接句柄存入channel中，由于缓存channel的特性，获取连接时如果池中有连接，将直接返回，
// 如果池中没有连接，将阻塞或者新建连接（没超过最大限制的情况下）。

var (
	ErrInvalidConfig = errors.New("invalid pool config")
	ErrPoolClosed    = errors.New("pool closed")
)

type factory func() (io.Closer, error)

type Pool interface {
	Acquire() (io.Closer, error) // 获取资源
	Release(io.Closer) error     // 释放资源
	Close(io.Closer) error       // 关闭资源
	Shutdown() error             // 关闭池
}

type GenericPool struct {
	sync.Mutex
	pool        chan io.Closer
	maxOpen     int  // 池中最大资源数
	numOpen     int  // 当前池中资源数
	minOpen     int  // 池中最少资源数
	closed      bool // 池是否已关闭
	maxLifetime time.Duration
	factory     factory // 创建连接的方法
}

func NewGenericPool(minOpen, maxOpen int, maxLifetime time.Duration, factory factory) (*GenericPool, error) {
	if maxOpen <= 0 || minOpen > maxOpen {
		return nil, ErrInvalidConfig
	}
	p := &GenericPool{
		maxOpen:     maxOpen,
		minOpen:     minOpen,
		maxLifetime: maxLifetime,
		factory:     factory,
		pool:        make(chan io.Closer, maxOpen),
	}

	for i := 0; i < minOpen; i++ {
		closer, err := factory()
		if err != nil {
			continue
		}
		p.numOpen++
		p.pool <- closer
	}
	return p, nil
}

func (p *GenericPool) Acquire() (io.Closer, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}
	for {
		closer, err := p.getOrCreate()
		if err != nil {
			return nil, err
		}
		// todo maxLifttime处理
		return closer, nil
	}
}

func (p *GenericPool) getOrCreate() (io.Closer, error) {
	select {
	case closer := <-p.pool:
		return closer, nil
	default:
	}
	p.Lock()
	if p.numOpen >= p.maxOpen {
		closer := <-p.pool
		p.Unlock()
		return closer, nil
	}
	// 新建连接
	closer, err := p.factory()
	if err != nil {
		p.Unlock()
		return nil, err
	}
	p.numOpen++
	p.Unlock()
	return closer, nil
}

// 释放单个资源到连接池
func (p *GenericPool) Release(closer io.Closer) error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	p.pool <- closer
	p.Unlock()
	return nil
}

// 关闭单个资源
func (p *GenericPool) Close(closer io.Closer) error {
	p.Lock()
	closer.Close()
	p.numOpen--
	p.Unlock()
	return nil
}

// 关闭连接池，释放所有资源
func (p *GenericPool) Shutdown() error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	close(p.pool)
	for closer := range p.pool {
		closer.Close()
		p.numOpen--
	}
	p.closed = true
	p.Unlock()
	return nil
}
