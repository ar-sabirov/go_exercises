package main

import (
	"fmt"
	"time"
)

type trafficLights struct {
	state bool
}

func (tl *trafficLights) operate() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		<-ticker.C
		if tl.state {
			fmt.Println("Open directions 2,3")
		} else {
			fmt.Println("Open directions 0,1")
		}
		tl.state = !tl.state
	}
}

type RoadIntersection struct {
	tl         *trafficLights
	directions [4]chan int
}

func New() RoadIntersection {
	tl := trafficLights{}
	go tl.operate()

	var chans [4]chan int
	for i := range chans {
		chans[i] = make(chan int)
	}
	return RoadIntersection{tl: &tl, directions: chans}
}

func (ri *RoadIntersection) Spawn(cars []int, directions []int, arrivalTimes []int) {
	for i := range arrivalTimes {
		go func(ix int) {
			car, dir := cars[ix], directions[ix]
			seconds := time.Duration(arrivalTimes[ix])
			<-time.After(seconds * time.Second)
			fmt.Printf("Car %v arrived at Direction %v\n", car, dir)
			ri.directions[dir] <- car
		}(i)
	}
}

func (ri *RoadIntersection) Operate() {
	for {
		if ri.tl.state {
			select {
			case x := <-ri.directions[0]:
				fmt.Printf("Car %v crossed Direction %v\n", x, 0)
			case x := <-ri.directions[1]:
				fmt.Printf("Car %v crossed Direction %v\n", x, 1)
			default:
				// do nothing
			}
		} else {
			select {
			case x := <-ri.directions[2]:
				fmt.Printf("Car %v crossed Direction %v\n", x, 2)
			case x := <-ri.directions[3]:
				fmt.Printf("Car %v crossed Direction %v\n", x, 3)
			default:
				//do nothing
			}
		}
	}
}

func main() {
	is := New()

	cars := []int{1, 3, 5, 2, 4}
	directions := []int{0, 2, 1, 3, 4}
	arrivalTimes := []int{0, 5, 8, 15, 20}

	go is.Spawn(cars, directions, arrivalTimes)
	go is.Operate()

	time.Sleep(60 * time.Second)
}
