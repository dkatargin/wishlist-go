package worker

import "testing"

func TestParsePrice(t *testing.T) {
	cases := map[string]float64{
		"12345 ₽":    12345,
		"12345.67 ₽": 12345.67,
		"999":        999,
		"7 ₽":        7,
	}
	for in, want := range cases {
		got, err := parsePrice(in)
		if err != nil || got != want {
			t.Fatalf("parsePrice(%q) = %v, %v; ожидали %v", in, got, err, want)
		}
	}

	if _, err := parsePrice("₽"); err == nil {
		t.Fatal("ожидалась ошибка для строки без цифр")
	}
}
