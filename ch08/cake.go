// 流水线做面包
// author: baoqiang
// time: 2019/3/20 下午8:21
package ch08

import (
	"time"
	"fmt"
	"math/rand"
)

type Shop struct {
	Verbose        bool
	Cakes          int
	BakeTime       time.Duration
	BakeStdDev     time.Duration
	BakeBuf        int
	NumIcers       int
	IceTime        time.Duration
	IceStdDev      time.Duration
	IceBuf         int
	InscribeTime   time.Duration
	InscribeStdDev time.Duration
}

type cake int

func (s *Shop) baker(backed chan<- cake) {
	for i := 0; i < s.Cakes; i++ {
		c := cake(i)
		if s.Verbose {
			fmt.Println("baking", i)
		}
		work(s.BakeTime, s.BakeStdDev)
		backed <- c
	}
	close(backed)
}

func (s *Shop) icer(iced chan<- cake, baked <-chan cake) {
	for c := range baked {
		if s.Verbose {
			fmt.Println("icing", c)
		}
		work(s.IceTime, s.IceStdDev)
		iced <- c
	}
}

func (s *Shop) inscriber(iced <-chan cake) {
	for i := 0; i < s.Cakes; i++ {
		c := <-iced
		if s.Verbose {
			fmt.Println("inscribing", c)
		}
		work(s.InscribeTime, s.InscribeStdDev)
		if s.Verbose {
			fmt.Println("finished", c)
		}
	}
}

func (s *Shop) Work(runs int) {
	for run := 0; run < runs; run++ {
		baked := make(chan cake, s.BakeBuf)
		iced := make(chan cake, s.IceBuf)

		go s.baker(baked)

		for i := 0; i < s.NumIcers; i++ {
			go s.icer(iced, baked)
		}

		s.inscriber(iced)
	}
}

func work(d, stddev time.Duration) {
	delay := d + time.Duration(rand.NormFloat64()*float64(stddev))
	time.Sleep(delay)
}
