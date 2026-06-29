// internal/infrastructure/crawler/wb_test.go
package crawler

import "testing"

func TestParseWbNm(t *testing.T) {
	ok := map[string]int64{
		"https://www.wildberries.ru/catalog/420939815/detail.aspx":         420939815,
		"https://wildberries.ru/catalog/420939815/detail.aspx?targetUrl=X": 420939815,
	}
	for u, want := range ok {
		got, err := parseWbNm(u)
		if err != nil || got != want {
			t.Fatalf("parseWbNm(%q) = %d, %v; ожидали %d", u, got, err, want)
		}
	}
	if _, err := parseWbNm("https://www.wildberries.ru/"); err == nil {
		t.Fatal("ожидалась ошибка для ссылки без артикула")
	}
}

func TestVolToBasket(t *testing.T) {
	cases := map[int64]int{0: 1, 143: 1, 144: 2, 4209: 24, 100000: 31}
	for vol, want := range cases {
		if got := volToBasket(vol); got != want {
			t.Fatalf("volToBasket(%d) = %d; ожидали %d", vol, got, want)
		}
	}
}

func TestParseWbCardName(t *testing.T) {
	body := `{"nm_id":420939815,"imt_name":"MIYOO Mini Plus 64GB","photo_count":3}`
	if got := parseWbCardName(body); got != "MIYOO Mini Plus 64GB" {
		t.Fatalf("imt_name: %q", got)
	}
	if parseWbCardName(`{"nm_id":1}`) != "" {
		t.Fatal("без imt_name ожидали пустую строку")
	}
}

func TestWbImageURL(t *testing.T) {
	got := wbImageURL(24, 4209, 420939, 420939815)
	want := "https://basket-24.wbbasket.ru/vol4209/part420939/420939815/images/big/1.webp"
	if got != want {
		t.Fatalf("wbImageURL = %q; ожидали %q", got, want)
	}
}

func TestWbClient_Supports(t *testing.T) {
	var _ Adapter = NewWbClient()
	if !NewWbClient().Supports("wildberries.ru") || NewWbClient().Supports("ozon.ru") {
		t.Fatal("Supports wildberries.ru неверен")
	}
}
