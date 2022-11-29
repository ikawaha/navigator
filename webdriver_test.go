package navigator

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestChromeDriver(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, err := os.ReadFile("testdata/hello.html")
		if err != nil {
			t.Fatalf("os.ReadFile() failed: unexpected error %v", err)
		}
		w.Write(b)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	d := ChromeDriver(ChromeOptions("args", []string{"--headless"}))
	if err := d.Start(); err != nil {
		t.Errorf("d.Start() failed: unexpected error %v", err)
		return
	}
	page, err := d.NewPage()
	if err != nil {
		t.Errorf("d.NewPage() failed: unexpected error %v", err)
		return
	}
	if err := page.Navigate(ts.URL + "/hello"); err != nil {
		t.Errorf("page.Navigate() failed: unexpected error %v", err)
		return
	}
	if got, err := page.Title(); err != nil {
		t.Errorf("page.Title() failed: unexpected error %v", err)
	} else if want := "Hello"; got != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	if err := d.Stop(); err != nil {
		t.Errorf("d.Stop() failed: unexpected error %v", err)
	}
}
