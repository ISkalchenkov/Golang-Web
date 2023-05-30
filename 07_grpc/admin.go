package main

import (
	"fmt"
	"time"
)

func NewAdminModule(subs *EventSubs, stats *StatTracker) *AdminModule {
	return &AdminModule{
		LogSubs: subs,
		Stats:   stats,
	}
}

type AdminModule struct {
	UnimplementedAdminServer

	LogSubs *EventSubs
	Stats   *StatTracker
}

func (adm *AdminModule) Logging(_ *Nothing, outStream Admin_LoggingServer) error {
	eventCh := adm.LogSubs.Subscribe(outStream)

Loop:
	for {
		select {
		case event := <-eventCh:
			if err := outStream.Send(event); err != nil {
				return fmt.Errorf("sending event to stream failed: %w", err)
			}
		case <-outStream.Context().Done():
			adm.LogSubs.Unsubscribe(outStream)
			break Loop
		}
	}

	return nil
}

func (adm *AdminModule) Statistics(interval *StatInterval, outStream Admin_StatisticsServer) error {
	adm.Stats.Subscribe(outStream)
	intervalDuration := time.Duration(interval.IntervalSeconds) * time.Second
	timer := time.NewTimer(intervalDuration)

Loop:
	for {
		select {
		case <-timer.C:
			stat, err := adm.Stats.Pull(outStream)
			if err != nil {
				return fmt.Errorf("pulling stat failed: %w", err)
			}

			if err := outStream.Send(stat); err != nil {
				return fmt.Errorf("sending stat to stream failed: %w", err)
			}
			timer.Reset(intervalDuration)
		case <-outStream.Context().Done():
			if !timer.Stop() {
				<-timer.C
			}
			adm.Stats.Unsubscribe(outStream)
			break Loop
		}
	}

	return nil
}
