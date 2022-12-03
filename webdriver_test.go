package navigator

import (
	"context"
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
	mux.HandleFunc("/reserve_app", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, err := os.ReadFile("testdata/reserve_app.html")
		if err != nil {
			t.Fatalf("os.ReadFile() failed: unexpected error %v", err)
		}
		w.Write(b)
	})
	mux.HandleFunc("/check_info.html", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("r.ParseForm() failed: unexpected error, %v", err)
		}
		want := map[string]string{
			"reserve_y": "2011",
			"reserve_m": "12",
			"reserve_d": "13",
			"reserve_t": "4",
			"hc":        "5",
			"bf":        "off",
			"gname":     "名前",
		}
		for k, v := range want {
			if got, want := r.Form.Get(k), v; got != want {
				t.Errorf("want %q=%q, got %q", k, v, got)
			}
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, err := os.ReadFile("testdata/table.html")
		if err != nil {
			t.Fatalf("os.ReadFile() failed: unexpected error %v", err)
		}
		w.Write(b)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	d := ChromeDriver(ChromeOptions("args", []string{"--headless"}))
	if err := d.Start(context.Background()); err != nil {
		t.Errorf("d.Start() failed: unexpected error %v", err)
		return
	}
	page, err := d.NewPage()
	if err != nil {
		t.Errorf("d.NewPage() failed: unexpected error %v", err)
		return
	}
	t.Run("testdata/hello.html", func(t *testing.T) {
		if err := page.Navigate(ts.URL + "/hello"); err != nil {
			t.Errorf("page.Navigate() failed: unexpected error %v", err)
			return
		}
		if got, err := page.Title(); err != nil {
			t.Errorf("page.Title() failed: unexpected error %v", err)
		} else if want := "Hello"; got != want {
			t.Errorf("want %q, but got %q", want, got)
		}
		if got, err := page.FindByID("aloha").Text(); err != nil {
			t.Errorf("page.FindByID(\"aloha\").Text() failed: unexpected error %v", err)
		} else if want := "aloha!👋"; got != want {
			t.Errorf("want %q, but got %q", want, got)
		}
	})
	t.Run("testdata/reserve_app.html", func(t *testing.T) {
		if err := page.Navigate(ts.URL + "/reserve_app"); err != nil {
			t.Errorf("page.Navigate() failed: unexpected error %v", err)
			return
		}
		if got, err := page.Title(); err != nil {
			t.Errorf("page.Title() failed: unexpected error %v", err)
		} else if want := "予約情報入力"; got != want {
			t.Errorf("want %q, but got %q", want, got)
		}
		if err := page.FindByName("reserve_y").Fill("2011"); err != nil {
			t.Errorf("page.FindByName(\"reserve_y\").Fill() failed: unexpected error %v", err)
		}
		if err := page.FindByID("reserve_month").Fill("12"); err != nil {
			t.Errorf("page.FindByID(\"reserve_moth\").Fill() failed: unexpected error %v", err)
		}
		if err := page.FindByID("reserve_day").Fill("13"); err != nil {
			t.Errorf("page.FindByID(\"reserve_day\").Fill() failed: unexpected error %v", err)
		}
		if err := page.FindByID("reserve_term").Fill("4"); err != nil {
			t.Errorf("page.FindByID(\"reserve_t\").Fill() failed: unexpected error %v", err)
		}
		if err := page.FindByID("headcount").Fill("5"); err != nil {
			t.Errorf("page.FindByID(\"reserve_t\").Fill() failed: unexpected error %v", err)
		}
		if err := page.FindByID("breakfast_off").Click(); err != nil {
			t.Errorf("page.FindByID(\"breakfast_off\").Click() failed: unexpected error %v", err)
		}
		if err := page.FindByID("plan_b").Click(); err != nil {
			t.Errorf("page.FindByID(\"plan_b\").Click() failed: unexpected error %v", err)
		}
		if err := page.FindByID("guestname").Fill("名前"); err != nil {
			t.Errorf("page.FindByID(\"guestname\").Fill() failed: unexpected error %v", err)
		}
		if err := page.FindByID("goto_next").Submit(); err != nil {
			t.Errorf("page.FindByID(\"goto_next\").Submit() failed: unexpected error %v", err)
		}
	})

	t.Run("multi selection", func(t *testing.T) {
		if err := page.Navigate(ts.URL + "/table"); err != nil {
			t.Errorf("page.Navigate() failed: unexpected error %v", err)
			return
		}
		t.Run("<th> elements", func(t *testing.T) {
			th := page.FindByID("1").All("th")
			if th == nil {
				t.Fatalf("page.FindByID(\"1\").All(\"th\") returns unexpected nil")
			}
			count, err := th.Count()
			if err != nil {
				t.Fatalf("page.FindByID(\"1\").All(\"th\").Count() failed: unexpected error, %v", err)
			}
			if want := 4; count != want {
				t.Fatalf("want %d, but got %d", want, count)
			}
			for i, v := range []string{"国", "首都", "人口", "言語"} {
				if got, err := th.At(i).Text(); err != nil {
					t.Fatalf("page.FindByID(\"1\").All(\"th\").At(%d).Text() failed: unexpected error, %v", i, err)
				} else if want := v; got != want {
					t.Errorf("want %s, but got %s", want, got)
				}
			}
		})
		t.Run("<td> elements", func(t *testing.T) {
			th := page.FindByID("1").All("td")
			if th == nil {
				t.Fatalf("page.FindByID(\"1\").All(\"td\") returns unexpected nil")
			}
			count, err := th.Count()
			if err != nil {
				t.Fatalf("page.FindByID(\"1\").All(\"td\").Count() failed: unexpected error, %v", err)
			}
			if want := 8; count != want {
				t.Fatalf("want %d, but got %d", want, count)
			}
			for i, v := range []string{
				"アメリカ合衆国", "ワシントン D.C.", "3 億 9 百万人", "英語",
				"スウェーデン", "ストックホルム", "9 百万人", "スウェーデン語",
			} {
				if got, err := th.At(i).Text(); err != nil {
					t.Fatalf("page.FindByID(\"1\").All(\"th\").At(%d).Text() failed: unexpected error, %v", i, err)
				} else if want := v; got != want {
					t.Errorf("want %s, but got %s", want, got)
				}
			}
		})
	})

	if err := d.Stop(); err != nil {
		t.Errorf("d.Stop() failed: unexpected error %v", err)
	}
}
