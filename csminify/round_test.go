package csminify

import "testing"

func TestRoundTo(t *testing.T) {
	testRoundTo(0.12, 1, 0, t)
	testRoundTo(1.23, 1, 1, t)
	testRoundTo(2.34, 1, 2, t)

	testRoundTo(-0.12, 1, 0, t)
	testRoundTo(-1.23, 1, -1, t)
	testRoundTo(-2.34, 1, -2, t)

	testRoundTo(0.12, 0.5, 0, t)
	testRoundTo(1.23, 0.5, 1, t)
	testRoundTo(2.34, 0.5, 2.5, t)

	testRoundTo(-0.12, 0.5, 0, t)
	testRoundTo(-1.23, 0.5, -1, t)
	testRoundTo(-2.34, 0.5, -2.5, t)

	testRoundTo(0.123, 0.02, 0.12, t)
	testRoundTo(2.345, 0.02, 2.34, t)
	testRoundTo(3.456, 0.02, 3.46, t)

	testRoundTo(-0.123, 0.02, -0.12, t)
	testRoundTo(-2.345, 0.02, -2.34, t)
	testRoundTo(-3.456, 0.02, -3.46, t)
}

func testRoundTo(in, precision, exp float64, t *testing.T) {
	res := roundTo(in, precision)
	if res != exp {
		t.Errorf("Expected result %f for input %f but got %f", exp, in, res)
	}
}
