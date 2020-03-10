package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func wait() {
	limiter := rate.NewLimiter(3, 5)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for i := 0; ; i++ {
		fmt.Printf("%03d %s\n", i, time.Now().Format("2006-01-02 15:04:05.000"))
		err := limiter.Wait(ctx)
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}
	}
}

func reserve() {
	limiter := rate.NewLimiter(3, 5)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for i := 0; ; i++ {
		fmt.Printf("%03d %s\n", i, time.Now().Format("2006-01-02 15:04:05.000"))
		reserve := limiter.Reserve()
		if !reserve.OK() {
			//返回是异常的，不能正常使用
			fmt.Println("Not allowed to act! Did you remember to set lim.burst to be > 0 ?")
			return
		}
		delayD := reserve.Delay()
		fmt.Println("sleep delay ", delayD)
		time.Sleep(delayD)
		select {
		case <-ctx.Done():
			fmt.Println("timeout, quit")
			return
		default:
		}
	}
}

func allow() {
	limiter := rate.NewLimiter(3, 5)
	for i := 0; i < 50; i++ {
		if limiter.Allow() {
			fmt.Printf("%03d Ok  %s\n", i, time.Now().Format("2006-01-02 15:04:05.000"))
		} else {
			fmt.Printf("%03d Err %s\n", i, time.Now().Format("2006-01-02 15:04:05.000"))
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	allow()
}
