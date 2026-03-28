package service

import "testing"

func TestClassifyLoadThresholds(t *testing.T) {
	cases := []struct {
		name string
		w    float64
		want int
	}{
		{name: "unknown", w: 0, want: 0},
		{name: "light low", w: 3.99, want: 1},
		{name: "medium low", w: 4.0, want: 2},
		{name: "medium high", w: 7.99, want: 2},
		{name: "heavy low", w: 8.0, want: 3},
		{name: "heavy high", w: 14.99, want: 3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyLoad(tc.w)
			if int(got) != tc.want {
				t.Fatalf("expected %d, got %d", tc.want, got)
			}
		})
	}
}

func TestProfileBucketLabelsReflectThresholds(t *testing.T) {
	want := []string{
		"idle (<4W)",
		"light (4-8W)",
		"medium (8-15W)",
		"heavy (15W+)",
	}
	for i, label := range want {
		if got := profileBucketLabel(i); got != label {
			t.Fatalf("bucket %d: expected %q, got %q", i, label, got)
		}
	}
}
