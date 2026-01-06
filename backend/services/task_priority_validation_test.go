package services

import "testing"

func TestNormalizeTaskPriority_EnglishValues(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Low", "Low"},
		{"low", "Low"},
		{"Medium", "Medium"},
		{"medium", "Medium"},
		{"High", "High"},
		{"high", "High"},
		{"", "Medium"},
		{"   ", "Medium"},
	}

	for _, tc := range cases {
		got, err := normalizeTaskPriority(tc.in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("normalize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestNormalizeTaskPriority_PersianAliases(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"کم", "Low"},
		{"متوسط", "Medium"},
		{"زیاد", "High"},
		{"بالا", "High"},
	}

	for _, tc := range cases {
		got, err := normalizeTaskPriority(tc.in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("normalize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestNormalizeTaskPriority_Invalid(t *testing.T) {
	_, err := normalizeTaskPriority("Urgent")
	if err == nil {
		t.Fatalf("expected error for invalid priority")
	}
}
