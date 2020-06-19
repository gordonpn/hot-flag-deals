package filter

import (
	"math"
	"testing"
)

func TestRound(t *testing.T) {
	ceiling := 5.75
	expected := 6
	t.Run("ceiling", testRoundFunc(ceiling, expected))

	floor := 5.25
	expected = 5
	t.Run("floor", testRoundFunc(floor, expected))

	middle := 5.5
	expected = 6
	t.Run("middle", testRoundFunc(middle, expected))
}

func testRoundFunc(testFloat float64, expected int) func(*testing.T) {
	return func(t *testing.T) {
		result := round(testFloat)
		if result != expected {
			t.Errorf("Expected %d, but got %d", expected, result)
		}
	}
}

func TestGetMean(t *testing.T) {
	t.Run("test one", testGetMeanFunc([]int{5, 5, 5, 5, 5, 5, 5}, 5))
	t.Run("test two", testGetMeanFunc([]int{1, 2, 3, 4, 5, 6}, 3.5))
}

func testGetMeanFunc(testSlice []int, expected float64) func(*testing.T) {
	return func(t *testing.T) {
		result := getMean(testSlice)
		if result != expected {
			t.Errorf("Expected %f, but got %f", expected, result)
		}
	}
}

func TestGetMedian(t *testing.T) {
	t.Run("test one", testGetMedianFunc([]int{5, 3, 6, 3, 4}, 4))
	t.Run("test two", testGetMedianFunc([]int{10, 19, 28, 16, 17, 30}, 18))
}

func testGetMedianFunc(testSlice []int, expected int) func(*testing.T) {
	return func(t *testing.T) {
		result := getMedian(testSlice)
		if result != expected {
			t.Errorf("Expected %d, but got %d", expected, result)
		}
	}
}

func TestGetStandDev(t *testing.T) {
	t.Run("test one", testGetStandDevFunc([]int{10, 12, 23, 23, 16, 23, 21, 16}, 18, 4.8989794855664))
	t.Run("test two", testGetStandDevFunc([]int{10, 12, 23, 23, 16, 23, 21, 16, 89, 38}, 27.1, 21.920082116635))
}

func testGetStandDevFunc(testSlice []int, testMean float64, expected float64) func(*testing.T) {
	return func(t *testing.T) {
		result := getStandDev(testSlice, testMean)
		if !withTolerance(result, expected) {
			t.Errorf("Expected %f, but got %f", expected, result)
		}
	}
}

func withTolerance(a, b float64) bool {
	tolerance := 0.01
	if diff := math.Abs(a - b); diff < tolerance {
		return true
	}
	return false
}

func TestGetSkewness(t *testing.T) {
	t.Run("test one", testGetSkewnessFunc(70.5, 80, 19.33, -1.47))
}

func testGetSkewnessFunc(testMean float64, testMedian int, testStandDev float64, expected float64) func(*testing.T) {
	return func(t *testing.T) {
		result := getSkewness(testMean, testMedian, testStandDev)
		if !withTolerance(result, expected) {
			t.Errorf("Expected %f, but got %f", expected, result)
		}
	}
}
