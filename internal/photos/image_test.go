package photos

import (
	"testing"
)

func TestBase58Encode(t *testing.T) {
	got := base58Encode([]byte("Hello World!"))
	want := "2NEpo7TZRRrLZSi2U"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}

	got = base58Encode([]byte("The quick brown fox jumps over the lazy dog."))
	want = "USm3fpXnKG5EUBx2ndxBDMPVciP5hGey2Jh4NDv6gmeo1LkMeiKrLJUUBk6Z"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}

	got = base58Encode([]byte{0x00, 0x00, 0x28, 0x7f, 0xb4, 0xcd})
	want = "11233QC4"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
