package filter

import (
	types "github.com/gordonpn/hot-flag-deals/internal/data"
	log "github.com/sirupsen/logrus"
	"math"
	"sort"
	"time"
)

func getThresholds(threads []types.Thread) (viewsThreshold, votesThreshold int) {
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
	for _, thread := range threads {
		viewsSlice = append(viewsSlice, thread.Views)
		votesSlice = append(votesSlice, thread.Votes)
	}
	viewsMean = GetMean(viewsSlice)
	votesMean = GetMean(votesSlice)

	viewsMedian = GetMedian(viewsSlice)
	votesMedian = GetMedian(votesSlice)

	viewsStandDev = GetStandDev(viewsSlice, viewsMean)
	votesStandDev = GetStandDev(votesSlice, votesMean)

	viewsSkewness = GetSkewness(viewsMean, viewsMedian, viewsStandDev)
	votesSkewness = GetSkewness(votesMean, votesMedian, votesStandDev)

	if math.Abs(viewsSkewness) >= standDevThreshold {
		viewsThreshold = viewsMedian
	} else {
		viewsThreshold = Round(viewsMean)
	}
	if math.Abs(votesSkewness) >= standDevThreshold {
		votesThreshold = votesMedian
	} else {
		votesThreshold = Round(votesMean)
	}
	viewsThreshold = Round(float64(viewsThreshold) * viewsThresholdCoeff)
	votesThreshold = Round(float64(votesThreshold) * votesThresholdCoeff)
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

func GetMean(intSlice []int) (mean float64) {
	sum := 0
	for _, num := range intSlice {
		sum += num
	}
	mean = float64(sum) / float64(len(intSlice))
	return
}

func GetMedian(intSlice []int) (median int) {
	sort.Ints(intSlice)
	middle := len(intSlice) / 2
	if len(intSlice)%2 == 0 {
		median = (intSlice[middle-1] + intSlice[middle]) / 2
	} else {
		median = intSlice[middle]
	}
	return
}

func GetStandDev(intSlice []int, mean float64) (standDev float64) {
	for i := range intSlice {
		standDev += math.Pow(float64(intSlice[i])-mean, 2)
	}
	standDev = math.Sqrt(standDev / float64(len(intSlice)))
	return
}

func GetSkewness(mean float64, median int, standDev float64) (skewness float64) {
	skewness = (mean - float64(median)) * 3 / standDev
	return
}

func Filter(threads []types.Thread) (filteredThreads []types.Thread) {
	viewsThreshold, votesThreshold := getThresholds(threads)
	const TimeThreshold = 72

	for _, thread := range threads {
		if (thread.Views >= viewsThreshold && thread.Votes >= votesThreshold) && !thread.Seen {
			timeNow := time.Now()
			diffHours := timeNow.Sub(thread.DatePosted).Hours()
			if diffHours <= TimeThreshold {
				filteredThreads = append(filteredThreads, thread)
			}
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

func Round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}
