package main

import (
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

func (a *App) getThresholds() {
	var (
		postsSlice []int
		viewsSlice []int
		votesSlice []int
	)

	for _, thread := range a.threads {
		viewsSlice = append(viewsSlice, thread.Views)
		votesSlice = append(votesSlice, thread.Votes)
		postsSlice = append(postsSlice, thread.Posts)
	}

	a.viewsMedian = getMedian(viewsSlice)
	a.votesMedian = getMedian(votesSlice)
	a.postsMedian = getMedian(postsSlice)
}

func (a *App) filter() []thread {
	log.Info("Filtering deals")
	const TimeThreshold = 120
	var threads []thread
	for _, thread := range a.threads {
		timeNow := time.Now()
		diffHours := timeNow.Sub(thread.DatePosted).Hours()
		if thread.Notified || diffHours >= TimeThreshold {
			continue
		}
		if thread.Posts >= a.postsMedian && thread.Views >= a.viewsMedian && thread.Votes >= a.votesMedian {
			threads = append(threads, thread)
		}
	}
	return threads
}

func getMedian(intSlice []int) (median int) {
	sort.Ints(intSlice)
	middle := len(intSlice) / 2
	if len(intSlice)%2 == 0 {
		median = (intSlice[middle-1] + intSlice[middle]) / 2
	} else {
		median = intSlice[middle]
	}
	return
}
