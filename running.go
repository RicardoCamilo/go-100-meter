package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type Track struct {
	runningTrack []Lane
}

func NewTrack() *Track {
	return &Track{runningTrack: make([]Lane, 10)}
}

type Runner struct {
	lane     Lane
	name     string
	avgSpeed int
}

func NewRunner(name string, vel int, track *Track) Runner {
	var lane Lane = append([]string{name}, strings.Split(strings.Repeat("-", 99), "")...)
	track.runningTrack = append(track.runningTrack, lane)
	return Runner{name: name, lane: lane, avgSpeed: vel}
}

func calculateVelocity(avgSpeed int) time.Duration {
	return time.Duration(avgSpeed + (rand.Intn(40) - 20))
}

func (r Runner) run(finished chan<- string) {
	for i, _ := range r.lane {
		if i > 0 {
			<-time.After(calculateVelocity(r.avgSpeed) * time.Millisecond)
			r.lane[i], r.lane[i-1] = r.lane[i-1], r.lane[i]
		}
	}
	finished <- r.name
}

type Lane []string

func (l Lane) printLane() string {
	return strings.Trim(fmt.Sprintf( "%v\n", l), "[]")
}

func (t Track) printRunningTrack() {
	var buffer bytes.Buffer
	for _, lane := range t.runningTrack {
		buffer.WriteString(lane.printLane())
	}
	fmt.Fprint(os.Stdout, buffer.String())
}

var runnersInfo = []Runner{
	{name: "Usain Bolt", avgSpeed: 100},
	{name: "Jesse Owens", avgSpeed: 108},
	{name: "Justin Gatlin", avgSpeed: 116},
	{name: "Tyson Gay", avgSpeed: 112},
	{name: "Tommie Smith", avgSpeed: 122},
	{name: "Carl Lewis", avgSpeed: 101},
	{name: "Florence", avgSpeed: 102},
	{name: "Maurice Greene", avgSpeed: 117},
	{name: "Donovan Bailey", avgSpeed: 125},
	{name: "Michael Johnson", avgSpeed: 119},
}

func main() {
	track := NewTrack()
	var runners []Runner

	for _, ri := range runnersInfo {
		name := fmt.Sprintf("%15v", ri.name)
		runners = append(runners, NewRunner(name, ri.avgSpeed, track))
	}

	raceIsOver := make(chan struct{})
	finished := make(chan string, len(runners))

	var wg sync.WaitGroup

	for _, r := range runners {
		go func(runner Runner) {
			wg.Add(1)
			runner.run(finished)
			wg.Done()
		}(r)
	}

	// Printer
	go func() {
	loop:
		for {
			select {
			case <-raceIsOver:
				break loop
			default:
				track.printRunningTrack()
			}
		}
	}()

	var winners []string
	go func() {
		for runner := range finished {
			if len(winners) < 3 {
				winners = append(winners, strings.Trim(runner, " "))
			}
		}
	}()

	wg.Wait()

	track.printRunningTrack()
	close(raceIsOver)
	// signals that all runners cross finish line
	close(finished)

	fmt.Printf("\n1: %s,\n2: %s,\n3: %s\n", winners[0], winners[1], winners[2])
}
