package service

import "sync"

type Dispatcher struct {
	lock             sync.RWMutex // 加锁
	EventObserverMap map[string][]chan string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{EventObserverMap: make(map[string][]chan string)}
}

func (dis *Dispatcher) Subscribe(e string, ch chan string) {
	dis.lock.Lock()
	defer dis.lock.Unlock()
	if _, ok := dis.EventObserverMap[e]; ok {
		dis.EventObserverMap[e] = append(dis.EventObserverMap[e], ch)
	} else {
		dis.EventObserverMap[e] = []chan string{ch}
	}
}
func (dis *Dispatcher) ClearSubscribe(url string) {
	if data, exist := dis.EventObserverMap[url]; exist && len(data) > 0 {
		for _, val := range dis.EventObserverMap[url] {
			close(val)
		}
		delete(dis.EventObserverMap, url)
	}
}
func (dis *Dispatcher) Unsubscribe(e string, ch chan string) {
	dis.lock.Lock()
	defer dis.lock.Unlock()
	if _, ok := dis.EventObserverMap[e]; ok {
		for i := range dis.EventObserverMap[e] {
			if dis.EventObserverMap[e][i] == ch {
				if i == 0 {
					if len(dis.EventObserverMap[e]) == 1 {
						delete(dis.EventObserverMap, e)
					} else {
						dis.EventObserverMap[e] = dis.EventObserverMap[e][i+1:]
					}
				} else if i == len(dis.EventObserverMap[e]) {
					dis.EventObserverMap[e] = dis.EventObserverMap[e][:i]
				} else {
					dis.EventObserverMap[e] = append(dis.EventObserverMap[e][:i], dis.EventObserverMap[e][i+1:]...)
				}
			}
		}
	}
}
func (dis *Dispatcher) Post(e string, response string) bool {
	eventMap := dis.EventObserverMap
	flag := false
	if _, ok := eventMap[e]; ok {
		for _, handler := range eventMap[e] {
			if handler != nil {
				flag = true
				go func(handler chan string) {
					handler <- response
				}(handler)
			}
		}
	}
	return flag
}
