package goe

import (
	"fmt"
	"github.com/gin-gonic/gin"
	Injector "github.com/shenyisyn/goft-ioc"
	"log"
	"reflect"
	"strings"
	"sync"
)

type Bean interface {
	Name() string
}

var innerRouter *GoeTree
var innerRouterOnce sync.Once

func getInnerRouter() *GoeTree {
	innerRouterOnce.Do(func() {
		innerRouter = NewGoeTree()
	})
	return innerRouter
}

type Goe struct {
	*gin.Engine
	g        *gin.RouterGroup
	exprData map[string]interface{}
	// 当前路由组
	currentGroup string
}

// Ignite Goe构造函数
func Ignite(ginMiddlewares ...gin.HandlerFunc) *Goe {
	g := &Goe{Engine: gin.New(), exprData: map[string]interface{}{}}
	// 强制加载异常处理中间件
	g.Use(ErrorHandle())
	for _, handler := range ginMiddlewares {
		g.Use(handler)
	}
	Injector.BeanFactory.Set(g)
	Injector.BeanFactory.Set(NewGPAUtil())
	return g
}

func (this *Goe) applyAll() {
	for t, v := range Injector.BeanFactory.GetBeanMapper() {
		if t.Elem().Kind() == reflect.Struct {
			Injector.BeanFactory.Apply(v.Interface())
		}
	}
}

// Launch 启动
func (this *Goe) Launch(port int32, funcs ...interface{}) {
	this.AfterLaunch(funcs...)
	this.applyAll()
	getCronTask().Start()
	err := this.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
	}
}

// AfterLaunch 启动之后执行的方法
func (this *Goe) AfterLaunch(exprs ...interface{}) {
	for _, expr := range exprs {
		if exp, ok := expr.(func()); ok {
			exp()
		} else {
			fmt.Println("AfterLaunch只能接受func类型的参数")
		}
	}
}

func (this *Goe) getPath(relativePath string) string {
	g := "/" + this.currentGroup
	if g == "/" {
		g = ""
	}
	g = g + relativePath
	g = strings.Replace(g, "//", "/", -1)
	return g
}

// Handle 重载gin.Handle方法
func (this *Goe) Handle(httpMethod, relativePath string, handler interface{}) *Goe {
	if h := Covert(handler); h != nil {
		methods := strings.Split(httpMethod, ",")
		for _, method := range methods {
			getInnerRouter().addRoute(method, this.getPath(relativePath), h)
			this.g.Handle(httpMethod, relativePath, h)
		}
	}
	return this
}

// Mount 挂载
func (this *Goe) Mount(group string, classes ...IClass) *Goe {
	this.g = this.Group(group)
	// 利用接口进行控制器挂载
	for _, class := range classes {
		this.currentGroup = group
		class.Build(this)
		// 将控制器也加入到bean容器中
		this.Beans(class)
	}
	// 返回自己方便链式调用
	return this
}

// Attach 添加中间件
func (this *Goe) Attach(f ...Fairing) *Goe {
	for _, f1 := range f {
		Injector.BeanFactory.Set(f1)
	}
	getFairingHandler().AddFairing(f...)
	return this
}

// Beans 设定数据库连接对象
func (this *Goe) Beans(beans ...Bean) *Goe {
	// 取出bean名称，加入到exprData里面
	for _, bean := range beans {
		this.exprData[bean.Name()] = bean
		Injector.BeanFactory.Set(bean)
	}
	return this
}

func (this *Goe) Config(cfgs ...interface{}) *Goe {
	Injector.BeanFactory.Config(cfgs...)
	return this
}

// Task 定时任务
func (this *Goe) Task(cron string, expr interface{}) *Goe {
	var err error
	if f, ok := expr.(func()); ok {
		_, err = getCronTask().AddFunc(cron, f)
	} else if exp, ok := expr.(Expr); ok {
		_, err = getCronTask().AddFunc(cron, func() {
			_, expErr := ExecExpr(exp, this.exprData)
			if expErr != nil {
				log.Println(expErr)
			}
		})
	}
	if err != nil {
		log.Println(err)
	}
	return this
}
