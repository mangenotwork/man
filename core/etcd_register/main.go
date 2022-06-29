package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main(){
	//Case1()

	Case2()
}

// 注册服务基础测试
func Case1(){


	// 连接 ETCD
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.0.192:2379"},
		DialTimeout: time.Second,
		//DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	Name := "helloworld"
	ID := "0"

	r := New(client)

	// ===========================================  监控
	w, err := r.Watch(ctx, Name)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()
	go func() {
		for {
			res, err1 := w.Next()
			if err1 != nil {
				return
			}
			log.Printf("watch: %v  %d", res, len(res))
			for _, r := range res {
				log.Printf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)


	//  ===========================================  服务注册
	if err1 := r.Register(ctx, Name, ID, "aaaa"); err1 != nil {
		log.Fatal(err1)
	}
	time.Sleep(time.Second)

	res, err := r.GetService(ctx, Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("res = ", res)

	//  ===========================================  删除注册
	if err1 := r.Deregister(ctx, Name, ID); err1 != nil {
		log.Fatal(err1)
	}
	time.Sleep(time.Second)

	res, err = r.GetService(ctx, Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("res = ", res)
	if len(res) != 0 {
		log.Println("not expected empty")
	}
}

// 注册服务租约测试
func Case2(){
	// 连接 ETCD
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.0.192:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	ID :=   "0"
	Name :=  "helloworld"


	// ===========================================  监控
	go func() {
		r := New(client)
		w, err1 := r.Watch(ctx, Name)
		if err1 != nil {
			return
		}
		defer func() {
			_ = w.Stop()
		}()
		for {
			res, err2 := w.Next()
			if err2 != nil {
				return
			}
			log.Printf("watch: %d", len(res))
			for _, r := range res {
				log.Printf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)




	// 新建一个etcd实例
	r := New(client,
		RegisterTTL(2*time.Second), // 租约时间2秒
		MaxRetry(0),
	)

	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, Name, ID)
	value := "aaaaa"
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		log.Fatal(err)
	}

	// 测试 ： 休眠3秒
	time.Sleep(3 * time.Second)

	res, err := r.GetService(ctx, Name)
	if err != nil {
		log.Fatal(err)
	}
	if len(res) != 0 {
		log.Println("not expected empty")
	}

	// 租约心跳
	go r.heartBeat(ctx, leaseID, key, value)

	time.Sleep(time.Second)
	res, err = r.GetService(ctx, Name)
	if err != nil {
		log.Fatal(err)
	}
	if len(res) == 0 {
		log.Println("key过期，没有key")
	}
}



// Option  扩展属性参数
type Option func(o *options)

type options struct {
	ctx       context.Context
	namespace string
	ttl       time.Duration
	maxRetry  int  // 最大重试次数
}

// Context with registry context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Namespace with registry namespace.
func Namespace(ns string) Option {
	return func(o *options) { o.namespace = ns }
}

// RegisterTTL with register ttl.
func RegisterTTL(ttl time.Duration) Option {
	return func(o *options) { o.ttl = ttl }
}

// MaxRetry 最大重试次数
func MaxRetry(num int) Option {
	return func(o *options) { o.maxRetry = num }
}

// Registry etcd注册服务
type Registry struct {
	opts   *options
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease  // 租约
}


// New 实例化注册服务
func New(client *clientv3.Client, opts ...Option) (r *Registry) {
	op := &options{
		ctx:       context.Background(),
		namespace: "/service",
		ttl:       time.Second * 15,
		maxRetry:  5,
	}
	for _, o := range opts {
		o(op)
	}
	return &Registry{
		opts:   op,
		client: client,
		kv:     clientv3.NewKV(client),
	}
}


// Register 服务注册
// 原理:
//		创建一个租约： clientv3.NewLease，
//		然后创建key写入值： r.registerWithKV(ctx, key, value)，
//		同时监听心跳： r.heartBeat(r.opts.ctx, leaseID, key, value)
func (r *Registry) Register(ctx context.Context, name, id, value string) error {
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace,name, id)
	if r.lease != nil {
		r.lease.Close()
	}
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		return err
	}

	go r.heartBeat(r.opts.ctx, leaseID, key, value)
	return nil
}


// Deregister 取消注册
// 原理: 使用 Delete 删除 key
func (r *Registry) Deregister(ctx context.Context, name, id string) error {
	defer func() {
		if r.lease != nil {
			r.lease.Close()
		}
	}()
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, name, id)
	_, err := r.client.Delete(ctx, key)
	return err
}


// GetService 通过key 获取值，
//	这个key 的结构是  namespace + name + id, 查询是  通过 namespace + name 获取所有的 namespace/name/id
func (r *Registry) GetService(ctx context.Context, name string) ([]string, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	resp, err := r.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		items = append(items, string(kv.Key)+":"+string(kv.Value))
	}
	return items, nil
}


// Watch 根据服务名称创建观察者。
// 主要功能是监控 key结构namespace + name + id 下 id 的个数，也就是说 /namespace/name/* key的个数
// 场景: 监控与上报， 服务熔断， 配置中心监听配置， 等等
// 原理: newWatcher(ctx, key, name, r.client)
func (r *Registry) Watch(ctx context.Context, name string) (*watcher, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	return newWatcher(ctx, key, name, r.client)
}


// registerWithKV 创建新租约，返回租约ID
// 原理:
//		lease.Grant(ctx, int64(r.opts.ttl.Seconds()))  设置租约过期时间为ttl
// 	注意: 续租时间大概自动为租约的三分之一时间(官方)
//		client.Put(ctx, key, value, clientv3.WithLease(grant.ID))   创建key写入值
func (r *Registry) registerWithKV(ctx context.Context, key string, value string) (clientv3.LeaseID, error) {
	grant, err := r.lease.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return 0, err
	}
	_, err = r.client.Put(ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return 0, err
	}
	return grant.ID, nil
}

// heartBeat 监听心跳  也可以理解为租约的业务实现
func (r *Registry) heartBeat(ctx context.Context, leaseID clientv3.LeaseID, key string, value string) {
	curLeaseID := leaseID
	kac, err := r.client.KeepAlive(ctx, leaseID)
	if err != nil {
		curLeaseID = 0
	}
	rand.Seed(time.Now().Unix())

	for {
		if curLeaseID == 0 {
			// 没有租约了，尝试注册
			retreat := []int{}
			for retryCnt := 0; retryCnt < r.opts.maxRetry; retryCnt++ {
				if ctx.Err() != nil {
					return
				}
				// 防止无限阻塞
				idChan := make(chan clientv3.LeaseID, 1)
				errChan := make(chan error, 1)
				cancelCtx, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()
					id, registerErr := r.registerWithKV(cancelCtx, key, value)
					if registerErr != nil {
						errChan <- registerErr
					} else {
						idChan <- id
					}
				}()

				select {
				case <-time.After(3 * time.Second):
					cancel()
					continue
				case <-errChan:
					continue
				case curLeaseID = <-idChan:
				}

				kac, err = r.client.KeepAlive(ctx, curLeaseID)
				if err == nil {
					break
				}
				retreat = append(retreat, 1<<retryCnt)
				time.Sleep(time.Duration(retreat[rand.Intn(len(retreat))]) * time.Second)
			}
			if _, ok := <-kac; !ok {
				// retry failed
				return
			}
		}

		select {
		case _, ok := <-kac:
			if !ok {
				if ctx.Err() != nil {
					// channel closed due to context cancel
					return
				}
				// 需要重试注册, 强行让 租约id=0 执行上面的逻辑
				curLeaseID = 0
				continue
			}
		case <-r.opts.ctx.Done():
			return
		}
	}
}



// =========================================================  watcher

// 监控
type watcher struct {
	key         string
	ctx         context.Context
	cancel      context.CancelFunc
	watchChan   clientv3.WatchChan
	watcher     clientv3.Watcher
	kv          clientv3.KV
	first       bool
	serviceName string
}

// newWatcher 新建监控
// 监控 /namespace/name
func newWatcher(ctx context.Context, key, name string, client *clientv3.Client) (*watcher, error) {
	w := &watcher{
		key:         key,
		watcher:     clientv3.NewWatcher(client),
		kv:          clientv3.NewKV(client),
		first:       true,
		serviceName: name,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.watchChan = w.watcher.Watch(w.ctx, key, clientv3.WithPrefix(), clientv3.WithRev(0))
	err := w.watcher.RequestProgress(context.Background())
	if err != nil {
		return nil, err
	}
	return w, nil
}

// 迭代
func (w *watcher) Next() ([]string, error) {
	if w.first {
		item, err := w.getInstance()
		w.first = false
		return item, err
	}

	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
		return w.getInstance()
	}
}

// 关闭监听的 chan
func (w *watcher) Stop() error {
	w.cancel()
	return w.watcher.Close()
}

// 获取值
func (w *watcher) getInstance() ([]string, error) {
	resp, err := w.kv.Get(w.ctx, w.key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		items = append(items, string(kv.Key)+":"+string(kv.Value))
	}
	return items, nil
}
