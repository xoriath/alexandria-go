package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/golden"
)

type TestProductsListener struct {
	data []string
}

func (t *TestProductsListener) Products() []string {
	return t.data
}

func TestCatalogHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/catalogs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewCatalogHandler(&TestProductsListener{
		data: []string{"Product1", "Product2"},
	}, "../templates")

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	assert.Assert(t, golden.String(rr.Body.String(), "catalogs.golden"))
}
