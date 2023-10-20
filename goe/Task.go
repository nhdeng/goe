package goe

import (
	"github.com/robfig/cron/v3"
	"sync"
)

// TaskFunc 任务函数类型
type TaskFunc func(params ...interface{})

// 携程任务
var taskList chan *TaskExecutor
var once sync.Once

// 定时任务
var onceCron sync.Once
var taskCron *cron.Cron

func init() {
	// 得到任务列表
	chList := getTaskList()
	go func() {
		for t := range chList {
			doTask(t)
		}
	}()
}

func doTask(t *TaskExecutor) {
	//TODO: 是否需要携程池/MQ？
	go func() {
		defer func() {
			if t.callback != nil {
				t.callback()
			}
		}()
		t.Exec()
	}()
}

func getCronTask() *cron.Cron {
	onceCron.Do(func() {
		taskCron = cron.New(cron.WithSeconds())
	})
	return taskCron
}

func getTaskList() chan *TaskExecutor {
	once.Do(func() {
		taskList = make(chan *TaskExecutor)
	})
	return taskList
}

// TaskExecutor 任务执行者
type TaskExecutor struct {
	f        TaskFunc
	p        []interface{}
	callback func()
}

func NewTaskExecutor(f TaskFunc, p []interface{}, callback func()) *TaskExecutor {
	return &TaskExecutor{f: f, p: p, callback: callback}
}

// Exec 执行任务
func (this *TaskExecutor) Exec() {
	this.f(this.p...)
}

func Task(f TaskFunc, cb func(), params ...interface{}) {
	if f == nil {
		return
	}
	// 增加任务队列
	go func() {
		getTaskList() <- NewTaskExecutor(f, params, cb)
	}()
}
