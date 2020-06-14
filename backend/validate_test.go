package main

import (
	"testing"
)

func TestSubscriberValidation(t *testing.T) {
	ross := subscriber{Name: "Ross", Email: "ross.geller@friends.com"}
	t.Run("Ross", testValidation(ross, true))

	ash := subscriber{Name: "X Æ A-12", Email: "ash.musk@spacex.com"}
	t.Run("Ash", testValidation(ash, false))

	quebecois := subscriber{Name: "Jean-Michel", Email: "jean.michel@gmail.com"}
	t.Run("Quebecois", testValidation(quebecois, true))

	space := subscriber{Name: "Jean Michel", Email: "jean.michel@gmail.com"}
	t.Run("space", testValidation(space, true))

	incomplete := subscriber{Name: "Jean Michel", Email: "jean.michel"}
	t.Run("incomplete", testValidation(incomplete, false))

	sql := subscriber{Name: "ayylmao", Email: "ayylmao@gmail.com OR 1=1"}
	t.Run("sql", testValidation(sql, false))

	sql2 := subscriber{Name: "ayylmao", Email: "ayylmao@gmail.com; DROP TABLE subscribers"}
	t.Run("sql", testValidation(sql2, false))

}

func testValidation(testStruct subscriber, pass bool) func(*testing.T) {
	return func(t *testing.T) {
		if err := testStruct.Validate(); err != nil && pass {
			t.Errorf("expected validation on %v, got %v", testStruct, err)
		} else if err == nil && !pass {
			t.Errorf("expected to fail on %v, got %v", testStruct, err)
		}
	}
}
