package main

import (
	"fmt"
	"sync"
	"time"
)

// Ttype представляет таск
type Ttype struct {
	id         int
	cT         time.Time // время создания
	fT         time.Time // время выполнения
	taskRESULT string    // результат выполнения задания
}

func main() {
	taskCreturer := func(a chan<- Ttype) {
		for {
			ft := time.Now()
			if ft.Nanosecond()%2 > 0 { // условие для появления ошибочных тасков
				ft = time.Time{} // сбрасываем время в случае ошибки
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	taskWorker := func(a Ttype, wg *sync.WaitGroup) {
		defer wg.Done()

		if a.cT.IsZero() {
			a.taskRESULT = "something went wrong"
		} else {
			a.taskRESULT = "task has been succeeded"
		}

		a.fT = time.Now()

		time.Sleep(time.Millisecond * 150)
	}

	wg := sync.WaitGroup{}
	wg.Add(10) // количество рабочих горутин

	go func() {
		for t := range superChan {
			go taskWorker(t, &wg)
		}
	}()

	wg.Wait() // ожидание завершения всех рабочих горутин
	close(superChan)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan Ttype)

	go func() {
		for t := range superChan {
			if t.taskRESULT == "task has been succeeded" {
				doneTasks <- t
			} else {
				undoneTasks <- t
			}
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	fmt.Println("Errors:")
	for r := range undoneTasks {
		fmt.Println(r)
	}

	fmt.Println("Done tasks:")
	for r := range doneTasks {
		fmt.Println(r)
	}
}
