package filter

import (
	"math"
	"sort"
	"time"

	types "github.com/gordonpn/hot-flag-deals/mailer/pkg/data"
	log "github.com/sirupsen/logrus"
)

var (
	standDevThreshold   = 0.9
	viewsMean           float64
	viewsMedian         int
	viewsSkewness       float64
	viewsSlice          []int
	viewsStandDev       float64
	viewsThresholdCoeff = 1.05
	votesMean           float64
	votesMedian         int
	votesSkewness       float64
	votesSlice          []int
	votesStandDev       float64
	votesThresholdCoeff = 3.0
)

func getThresholds(threads []types.Thread) (viewsThreshold, votesThreshold int) {
	for _, thread := range threads {
		viewsSlice = append(viewsSlice, thread.Views)
		votesSlice = append(votesSlice, thread.Votes)
	}
	viewsMean = getMean(viewsSlice)
	votesMean = getMean(votesSlice)

	viewsMedian = getMedian(viewsSlice)
	votesMedian = getMedian(votesSlice)

	viewsStandDev = getStandDev(viewsSlice, viewsMean)
	votesStandDev = getStandDev(votesSlice, votesMean)

	viewsSkewness = getSkewness(viewsMean, viewsMedian, viewsStandDev)
	votesSkewness = getSkewness(votesMean, votesMedian, votesStandDev)

	if math.Abs(viewsSkewness) >= standDevThreshold {
		viewsThreshold = viewsMedian
	} else {
		viewsThreshold = round(viewsMean)
	}
	if math.Abs(votesSkewness) >= standDevThreshold {
		votesThreshold = votesMedian
	} else {
		votesThreshold = round(votesMean)
	}
	viewsThreshold = round(float64(viewsThreshold) * viewsThresholdCoeff)
	votesThreshold = round(float64(votesThreshold) * votesThresholdCoeff)
	log.WithFields(log.Fields{
		"viewsMean":           viewsMean,
		"viewsMedian":         viewsMedian,
		"viewsSkewness":       viewsSkewness,
		"viewsThreshold":      viewsThreshold,
		"viewsThresholdCoeff": viewsThresholdCoeff,
	}).Debug()
	log.WithFields(log.Fields{
		"votesMean":           votesMean,
		"votesMedian":         votesMedian,
		"votesSkewness":       votesSkewness,
		"votesThreshold":      votesThreshold,
		"votesThresholdCoeff": votesThresholdCoeff,
	}).Debug()
	return
}

func getMean(intSlice []int) (mean float64) {
	sum := 0
	for _, num := range intSlice {
		sum += num
	}
	mean = float64(sum) / float64(len(intSlice))
	return
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

func getStandDev(intSlice []int, mean float64) (standDev float64) {
	for i := range intSlice {
		standDev += math.Pow(float64(intSlice[i])-mean, 2)
	}
	standDev = math.Sqrt(standDev / float64(len(intSlice)))
	return
}

func getSkewness(mean float64, median int, standDev float64) (skewness float64) {
	skewness = (mean - float64(median)) * 3 / standDev
	return
}

func Filter(threads []types.Thread) (filteredThreads []types.Thread) {
	viewsThreshold, votesThreshold := getThresholds(threads)
	const TimeThreshold = 72

	for _, thread := range threads {
		timeNow := time.Now()
		diffHours := timeNow.Sub(thread.DatePosted).Hours()
		if thread.Seen || diffHours > TimeThreshold {
			continue
		}
		if (thread.Views >= viewsThreshold && thread.Votes > int(votesMean)) || thread.Votes >= votesThreshold {
			filteredThreads = append(filteredThreads, thread)
		}
	}
	sort.SliceStable(filteredThreads, func(this, that int) bool {
		return filteredThreads[this].Votes > filteredThreads[that].Votes
	})
	log.WithFields(log.Fields{
		"len(filteredThreads)": len(filteredThreads),
		"cap(filteredThreads)": cap(filteredThreads)},
	).Debug("Length and capacity of filtered threads")
	return
}

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}
