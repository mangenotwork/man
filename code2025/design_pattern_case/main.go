package main

import (
	"fmt"
	"sync"
)

// 设计模式

func main() {
	//case1()
	//case2()
	//case3()
	//case4()
	//case5()
	//case6()
	//case7()
	//case8()
	//case9()
	case10()
}

// 单例模式

// 定义一个结构体Singleton，用于存储单例的实例数据

type singleton struct {
	value string // 这里可以存储单例对象的任何数据
}

// 定义一个Once对象，用于确保初始化操作只执行一次
var once sync.Once

// 定义一个全局变量instance，用于存储单例的实例
var instance *singleton

// 初始化函数，由Once.Do调用
func initSingleton() {
	instance = &singleton{value: "unique instance"} // 这里初始化singleton实例
}

// getInstance函数用于获取单例的实例
func getInstance() *singleton {
	// 执行initSingleton，确保instance只被初始化一次
	once.Do(initSingleton)
	return instance // 返回单例的实例
}

func case1() {
	// 获取单例的实例
	singletonInstance := getInstance()
	fmt.Println(singletonInstance.value) // 输出: unique instance

	// 再次获取单例的实例，将返回相同的实例
	anotherInstance := getInstance()
	if singletonInstance == anotherInstance {
		fmt.Println("Both instances are the same") // 输出: Both instances are the same
	}

	// 测试并发环境下的单例模式
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			singletonInstance := getInstance()
			fmt.Println(singletonInstance.value)
		}()
	}
	wg.Wait()
}

// 工厂模式

// 定义一个接口Product，它声明了所有具体产品对象必须实现的操作
type Product interface {
	operation() // 产品对象的操作
}

// 定义具体产品ConcreteProductA，实现了Product接口
type ConcreteProductA struct{}

func (p *ConcreteProductA) operation() {
	fmt.Println("Operation of ConcreteProductA")
}

// 定义另一个具体产品ConcreteProductB，也实现了Product接口
type ConcreteProductB struct{}

func (p *ConcreteProductB) operation() {
	fmt.Println("Operation of ConcreteProductB")
}

// 定义一个抽象工厂Creator，它声明了工厂方法factoryMethod，用于创建产品对象
type Creator interface {
	factoryMethod() Product // 工厂方法，用于创建产品对象
}

// 定义具体工厂CreatorA，实现了Creator接口
type CreatorA struct{}

func (c *CreatorA) factoryMethod() Product {
	return &ConcreteProductA{} // 具体工厂CreatorA返回ConcreteProductA的实例
}

// 定义另一个具体工厂CreatorB，也实现了Creator接口
type CreatorB struct{}

func (c *CreatorB) factoryMethod() Product {
	return &ConcreteProductB{} // 具体工厂CreatorB返回ConcreteProductB的实例
}

func case2() {
	// 创建具体工厂CreatorA的实例
	creatorA := &CreatorA{}
	productA := creatorA.factoryMethod()
	productA.operation() // 调用产品A的操作

	// 创建具体工厂CreatorB的实例
	creatorB := &CreatorB{}
	productB := creatorB.factoryMethod()
	productB.operation() // 调用产品B的操作
}

// 观察者模式

// 定义Observer接口，它声明了观察者需要实现的Update方法
type Observer interface {
	Update(string) // 当主题状态改变时，此方法会被调用
}

// 定义Subject结构体，它包含一个观察者列表和方法来添加或通知观察者
type Subject struct {
	observers []Observer // 存储观察者的列表
}

// Attach方法用于将一个观察者添加到观察者列表中
func (s *Subject) Attach(observer Observer) {
	s.observers = append(s.observers, observer)
}

// Notify方法用于通知所有观察者主题状态的改变
func (s *Subject) Notify(message string) {
	for _, observer := range s.observers {
		observer.Update(message) // 调用每个观察者的Update方法
	}
}

// 定义一个具体观察者ConcreteObserver，它实现了Observer接口
type ConcreteObserverA struct {
	name string
}

// 实现Observer接口的Update方法
func (c *ConcreteObserverA) Update(message string) {
	fmt.Printf("%s received message: %s\n", c.name, message)
}

func case3() {
	// 创建主题对象
	subject := &Subject{}

	// 创建具体观察者对象
	observerA := &ConcreteObserverA{name: "Observer A"}
	observerB := &ConcreteObserverA{name: "Observer B"}

	// 将观察者添加到主题的观察者列表中
	subject.Attach(observerA)
	subject.Attach(observerB)

	// 当主题状态改变时，通知所有观察者
	subject.Notify("State changed to State 1")
}

// 装饰者模式

// 定义Component接口，它是所有组件和装饰者的基类
type Component interface {
	operation() // 组件执行的操作
}

// 定义具体组件ConcreteComponent，实现了Component接口
type ConcreteComponent struct{}

func (c *ConcreteComponent) operation() {
	fmt.Println("ConcreteComponent: performing basic operation")
}

type ConcreteComponent2 struct{}

func (c *ConcreteComponent2) operation() {
	fmt.Println("ConcreteComponent2: performing basic operation2")
}

// 定义Decorator抽象结构体，它包含一个Component接口类型的字段
type Decorator struct {
	component Component // 用于组合Component接口
}

// 实现Decorator的operation方法，调用其Component的operation方法
func (d *Decorator) operation() {
	if d.component != nil {
		d.component.operation() // 调用被装饰者的operation
	}
}

// 定义具体装饰者ConcreteDecoratorA，它嵌入了Decorator结构体
type ConcreteDecoratorA struct {
	Decorator // 继承Decorator，实现装饰功能
}

// 为ConcreteDecoratorA实现operation方法，添加额外的职责
func (cda *ConcreteDecoratorA) operation() {
	cda.Decorator.operation() // 首先调用被装饰者的operation
	fmt.Println("ConcreteDecoratorA: added additional responsibilities")
}

func case4() {
	// 创建具体组件
	component := &ConcreteComponent{}

	// 创建装饰者并关联具体组件
	decoratorA := &ConcreteDecoratorA{Decorator{component}}

	// 执行装饰后的组件操作
	decoratorA.operation()

	component2 := &ConcreteComponent2{}
	decoratorA.Decorator.component = component2
	decoratorA.operation()
}

// 策略模式

// 定义Strategy接口，它声明了所有具体策略必须实现的algorithm方法
type Strategy interface {
	algorithm() // 策略的算法方法
}

// 定义具体策略ConcreteStrategyA，实现了Strategy接口
type ConcreteStrategyA struct{}

func (c *ConcreteStrategyA) algorithm() {
	fmt.Println("Executing Algorithm A")
}

// 定义另一个具体策略ConcreteStrategyB，也实现了Strategy接口
type ConcreteStrategyB struct{}

func (c *ConcreteStrategyB) algorithm() {
	fmt.Println("Executing Algorithm B")
}

// 定义Context结构体，它包含一个Strategy接口类型的字段
type Context struct {
	strategy Strategy // 用于存储当前使用的策略
}

// 执行策略的方法，通过Context中的strategy字段调用algorithm方法
func (c *Context) executeStrategy() {
	c.strategy.algorithm() // 执行当前策略的算法
}

func case5() {
	// 创建Context对象
	context := &Context{}

	// 创建具体策略对象
	strategyA := &ConcreteStrategyA{}
	strategyB := &ConcreteStrategyB{}

	// 将Context的策略设置为策略A
	context.strategy = strategyA
	context.executeStrategy() // 输出: Executing Algorithm A

	// 更换策略为策略B
	context.strategy = strategyB
	context.executeStrategy() // 输出: Executing Algorithm B
}

// 适配器模式

// 定义Target接口，表示客户端使用的特定领域相关的接口
type Target interface {
	request() // 客户端期望调用的方法
}

// 定义一个已经存在的类Adaptee，它有自己的接口
type Adaptee struct{}

func (a *Adaptee) specificRequest() {
	fmt.Println("Adaptee performs a specific request")
}

// 定义Adapter结构体，它作为Target接口和Adaptee类之间的桥梁
type Adapter struct {
	adaptee *Adaptee // 引用Adaptee对象
}

// Adapter实现了Target接口的request方法
// 该方法内部委托给Adaptee的specificRequest方法
func (a *Adapter) request() {
	if a.adaptee != nil {
		a.adaptee.specificRequest() // 委托调用Adaptee的方法
	}
}

func case6() {
	// 创建Adaptee对象
	adaptee := &Adaptee{}

	// 创建Adapter对象，并注入Adaptee对象
	adapter := &Adapter{adaptee: adaptee}

	// 客户端使用Target接口，这里通过Adapter实现
	var target Target = adapter
	target.request() // 通过Adapter调用Adaptee的方法
}

// 代理模式

// SubjectProxy 定义主题接口
type SubjectProxy interface {
	Request() string
}

// RealSubject 实现具体主题
type RealSubject struct{}

// Request 实现请求方法
func (r *RealSubject) Request() string {
	return "RealSubject: Handling request."
}

// Proxy 实现代理主题
type Proxy struct {
	realSubject *RealSubject
}

// NewProxy 创建代理实例
func NewProxy() *Proxy {
	return &Proxy{
		realSubject: &RealSubject{},
	}
}

// Request 代理请求方法
func (p *Proxy) Request() string {
	if p.checkAccess() {
		result := p.realSubject.Request()
		p.logAccess()
		return result
	}
	return "Proxy: Access denied."
}

// checkAccess 检查访问权限
func (p *Proxy) checkAccess() bool {
	fmt.Println("Proxy: Checking access prior to firing a real request.")
	return true // 简化示例，实际中可能包含复杂的权限检查
}

// logAccess 记录访问日志
func (p *Proxy) logAccess() {
	fmt.Println("Proxy: Logging the time of request.")
}

func case7() {
	proxy := NewProxy()
	result := proxy.Request()
	fmt.Println(result)
}

// 命令模式

// 定义Command接口，它声明了所有具体命令必须实现的Execute方法
type Command interface {
	Execute() // 执行命令的方法
}

// 定义Receiver结构体，它将执行命令的实际请求
type Receiver struct{}

func (r *Receiver) Action() {
	fmt.Println("Receiver: Action")
}

// 定义ConcreteCommand结构体，它实现了Command接口
// 每个具体命令都包含一个Receiver的引用，表示请求的接收者
type ConcreteCommand struct {
	receiver *Receiver // 命令执行的接收者
}

// ConcreteCommand实现Command接口的Execute方法
// 该方法调用Receiver的Action方法来执行请求
func (c *ConcreteCommand) Execute() {
	c.receiver.Action() // 执行请求
}

// 定义Invoker结构体，它负责调用命令对象的Execute方法
type Invoker struct {
	command Command // 存储命令对象
}

// 调用命令对象的Execute方法
func (i *Invoker) Invoke() {
	i.command.Execute() // 执行命令
}

func case8() {
	// 创建接收者对象
	receiver := &Receiver{}

	// 创建具体命令对象，并注入接收者
	command := &ConcreteCommand{receiver: receiver}

	// 创建调用者对象，并注入具体命令对象
	invoker := &Invoker{command: command}

	// 调用者执行命令
	invoker.Invoke() // 输出: Receiver: Action
}

// 组合模式

// 定义Component接口，作为组合中对象的一致性协议
type ComponentZH interface {
	Operation()               // 执行操作的方法
	Add(ComponentZH)          // 向组合中添加子节点的方法
	Remove(ComponentZH)       // 从组合中移除子节点的方法
	GetChild(int) ComponentZH // 根据索引获取子节点的方法
}

// 定义Leaf结构体，表示组合中的叶节点
type Leaf struct {
	name string
}

// Leaf实现Component接口的Operation方法
func (l *Leaf) Operation() {
	fmt.Println("Leaf:", l.name)
}

// Leaf实现Component接口的Add方法，叶节点不能有子节点，因此这里可以不实现或抛出错误
func (l *Leaf) Add(c ComponentZH) {
	fmt.Println("Cannot add to a leaf")
}

// Leaf实现Component接口的Remove方法，叶节点不能有子节点，因此这里可以不实现或抛出错误
func (l *Leaf) Remove(c ComponentZH) {
	fmt.Println("Cannot remove from a leaf")
}

// Leaf实现Component接口的GetChild方法，叶节点没有子节点，因此这里返回nil
func (l *Leaf) GetChild(i int) ComponentZH {
	return nil
}

// 定义Composite结构体，表示组合中的容器节点
type Composite struct {
	name     string
	Children []ComponentZH // 存储子节点的列表
}

// Composite实现Component接口的Operation方法
func (c *Composite) Operation() {
	fmt.Println("Composite:", c.name)
	for _, child := range c.Children {
		child.Operation() // 递归调用子节点的Operation方法
	}
}

// Composite实现Component接口的Add方法，向Children列表中添加子节点
func (c *Composite) Add(component ComponentZH) {
	c.Children = append(c.Children, component)
}

// Composite实现Component接口的Remove方法，从Children列表中移除子节点
func (c *Composite) Remove(component ComponentZH) {
	for i, child := range c.Children {
		if child == component {
			c.Children = append(c.Children[:i], c.Children[i+1:]...)
			break
		}
	}
}

// Composite实现Component接口的GetChild方法，根据索引获取子节点
func (c *Composite) GetChild(i int) ComponentZH {
	if i < 0 || i >= len(c.Children) {
		return nil // 索引超出范围，返回nil
	}
	return c.Children[i]
}

func case9() {
	// 创建叶节点
	leafA := &Leaf{name: "Leaf A"}
	leafB := &Leaf{name: "Leaf B"}

	// 创建组合节点
	composite := &Composite{name: "Composite Root"}
	composite.Add(leafA) // 向组合中添加叶节点A
	composite.Add(leafB) // 向组合中添加叶节点B

	// 执行组合节点的操作
	composite.Operation()
}

// 迭代器模式

// 定义Iterator接口，它声明了迭代器必须实现的Next和Current方法
type Iterator interface {
	Next() bool           // 移动到下一个元素，并返回是否成功移动
	Current() interface{} // 返回当前元素
}

// 定义ConcreteIterator结构体，它实现了Iterator接口
type ConcreteIterator struct {
	items []string // 存储聚合对象的元素列表
	index int      // 当前迭代到的元素索引
}

// Next方法实现，用于移动到下一个元素
func (c *ConcreteIterator) Next() bool {
	if c.index < len(c.items) {
		c.index++ // 索引递增
		return true
	}
	return false // 如果索引超出范围，返回false
}

// Current方法实现，用于返回当前元素
func (c *ConcreteIterator) Current() interface{} {
	if c.index > 0 && c.index <= len(c.items) {
		return c.items[c.index-1] // 返回当前索引的元素
	}
	return nil // 如果索引不在范围内，返回nil
}

// 定义Aggregate接口，表示聚合对象，它将负责创建迭代器
type Aggregate interface {
	CreateIterator() Iterator // 创建并返回迭代器
}

// 定义ConcreteAggregate结构体，它实现了Aggregate接口
type ConcreteAggregate struct {
	items []string // 聚合对象存储的元素列表
}

// CreateIterator方法实现，用于创建并返回一个迭代器
func (a *ConcreteAggregate) CreateIterator() Iterator {
	return &ConcreteIterator{items: a.items, index: 0} // 返回一个新的迭代器实例
}

func case10() {
	// 创建聚合对象并添加元素
	aggregate := &ConcreteAggregate{items: []string{"Item1", "Item2", "Item3"}}

	// 使用聚合对象创建迭代器
	iterator := aggregate.CreateIterator()

	// 使用迭代器遍历聚合对象中的所有元素
	for iterator.Next() {
		fmt.Println(iterator.Current())
	}
}
