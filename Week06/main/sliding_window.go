package main

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	typeSuccess int = 1
	typeFail    int = 2
)

type metrics struct {
	success int64
	fail    int64
}

type SlidingWindow struct {
	bucket int
	curKey int64
	m      map[int64]*metrics
	data   *list.List
	sync.RWMutex
}

func (sw *SlidingWindow) incr(t int) {
	sw.Lock()
	defer sw.Unlock()
	nowTime := time.Now().Unix()
	if _, ok := sw.m[nowTime]; !ok {
		sw.m = make(map[int64]*metrics)
		sw.m[nowTime] = &metrics{}
	}
	if sw.curKey == 0 {
		sw.curKey = nowTime
	}
	if sw.curKey != nowTime {
		sw.data.PushBack(sw.m[nowTime])
		delete(sw.m, sw.curKey)
		sw.curKey = nowTime
		if sw.data.Len() > sw.bucket {
			for i := 0; i <= sw.data.Len()-sw.bucket; i++ {
				sw.data.Remove(sw.data.Front())
			}
		}
	}
	switch t {
	case typeSuccess:
		sw.m[nowTime].success++
	case typeFail:
		sw.m[nowTime].fail++
	default:
		log.Fatal("err type")
	}

}

func (sw *SlidingWindow) Len() int {
	return sw.data.Len()
}

//获取数据(space 如：5、10秒)
func (sw *SlidingWindow) Data(space int) []*metrics {
	sw.RLock()
	defer sw.RUnlock()
	var data []*metrics
	var num = 0
	var m = &metrics{}
	for i := sw.data.Front(); i != nil; i = i.Next() {
		one := i.Value.(*metrics)
		m.success += one.success
		m.fail += one.fail
		if num%space == 0 {
			data = append(data, m)
			m = &metrics{} //重置m
		}
		num++
	}
	return data
}

//创建滑动窗口
func NewSlidingWindow(bucket int) *SlidingWindow {
	sw := &SlidingWindow{}
	sw.bucket = bucket
	sw.data = list.New()
	return sw
}

//统计成功
func (sw *SlidingWindow) AddSuccess() {
	sw.incr(typeSuccess)
}

//统计失败
func (sw *SlidingWindow) AddFail() {
	sw.incr(typeFail)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	sw := NewSlidingWindow(100)
	var r int
	for i := 0; i < 1000; i++ {
		r = rand.Intn(3)
		if r == 1 {
			sw.AddSuccess()
		} else {
			sw.AddFail()
		}
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
	}
	fmt.Println("1秒的bucket长度", sw.Len())
	for _, item := range sw.Data(3) {
		fmt.Println(item.success, item.fail)
	}
	fmt.Println("==============")
	for _, item := range sw.Data(5) {
		fmt.Println(item.success, item.fail)
	}
}
