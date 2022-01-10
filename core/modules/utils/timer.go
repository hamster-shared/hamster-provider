package utils

import "time"

type TimerService struct {
	timerMap  map[uint64]*time.Timer
	tickerMap map[uint64]*time.Ticker
}

func NewTimerService() *TimerService {

	return &TimerService{
		timerMap:  make(map[uint64]*time.Timer),
		tickerMap: make(map[uint64]*time.Ticker),
	}
}

func (s *TimerService) SubTimer(id uint64, timer *time.Timer) {
	s.timerMap[id] = timer
}

func (s *TimerService) UnSubTimer(id uint64) {
	timer := s.timerMap[id]

	if timer != nil {
		timer.Stop()
	}

	delete(s.timerMap, id)
}

func (s *TimerService) GetTimer(id uint64) *time.Timer {
	return s.timerMap[id]
}

func (s *TimerService) SubTicker(id uint64, ticker *time.Ticker) {
	s.tickerMap[id] = ticker
}

func (s *TimerService) UnSubTicker(id uint64) {
	ticker := s.tickerMap[id]

	if ticker != nil {
		ticker.Stop()
	}

	delete(s.tickerMap, id)
}

func (s *TimerService) GetTicker(id uint64) *time.Ticker {
	return s.tickerMap[id]
}
