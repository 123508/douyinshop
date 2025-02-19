package els

import "testing"

func TestSearchProduct(t *testing.T) {
	name := "手机"
	got, err := SearchProduct(name, 1, 2)
	if err != nil {
		t.Errorf("SearchProduct() error = %v", err)
		return
	}
	t.Logf("SearchProduct() got = %v", got)
}
