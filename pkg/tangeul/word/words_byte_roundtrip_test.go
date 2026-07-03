//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestWordsByteRoundTrip — ReadWords(PutWords(ws)) reproduces every word (BE)
package word

import "testing"

func TestWordsByteRoundTrip(t *testing.T) {
	want := []Word{0x0001, 0xD8F0, 0x2710, 0x0000, 0xFFFF}
	got, err := ReadWords(PutWords(want))
	if err != nil || len(got) != len(want) {
		t.Fatalf("len got %d want %d err=%v", len(got), len(want), err)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("word %d: got 0x%04X want 0x%04X", i, got[i], want[i])
		}
	}
}
