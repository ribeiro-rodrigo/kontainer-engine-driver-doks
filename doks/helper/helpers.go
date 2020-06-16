package helper

import "time"

type Sleeper interface {
	Sleep(duration time.Duration)
}

type TimerSleeper struct {}

func NewTimerSleeper()TimerSleeper{
	return TimerSleeper{}
}

func (TimerSleeper) Sleep(duration time.Duration){
	time.Sleep(duration)
}
