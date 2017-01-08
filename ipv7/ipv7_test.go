package ipv7

import (
	"testing"
)

func TestFirstABBA(t *testing.T) {
	cases := [...]struct {
		text string
		abba string
	}{
		{"abba", "abba"},
		{"ioxxoj", "oxxo"},
		{"oxxoj", "oxxo"},
		{"ioxxo", "oxxo"},
		{"abbaioxxo", "abba"},
	}

	for ncase, item := range cases {
		abba := firstABBA(item.text)
		if abba != item.abba {
			t.Errorf("[Case %d] Found '%s' in '%s'.  ('%s')", ncase,
				abba, item.text, item.abba)
		}
	}
}

func TestAllABA(t *testing.T) {
	cases := [...]struct {
		text string
		abas []string
	}{
		{"aba", []string{"aba"}},
		{"ioxoj", []string{"oxo"}},
		{"oxoj", []string{"oxo"}},
		{"ioxo", []string{"oxo"}},
		{"zazbzb", []string{"zaz", "zbz", "bzb"}},
	}

	for ncase, item := range cases {
		abas := allABA(item.text)
		if len(abas) != len(item.abas) {
			t.Errorf("[Case %d] Found %d ABAs in '%s' (expected %d)", ncase,
				len(abas), item.text, len(item.abas))
		}
		if abas[0] != item.abas[0] {
			t.Errorf("[Case %d] Found '%s' in '%s'.  ('%s')", ncase,
				abas[0], item.text, item.abas[0])
		}
	}
}

func TestParser(t *testing.T) {
	cases := [...]struct {
		addr      string
		is_tls    bool
		supernets []string
		hypernets []string
	}{
		{"abba[mnop]qrst", true, []string{"abba", "qrst"}, []string{"mnop"}},
		{"abcd[bddb]xyyx", false, []string{"abcd", "xyyx"}, []string{"bddb"}},
		{"aaaa[qwer]tyui", false, []string{"aaaa", "tyui"}, []string{"qwer"}},
		{"ioxxoj[asdfgh]zxcvbn", true, []string{"ioxxoj", "zxcvbn"},
			[]string{"asdfgh"}},
		{"oxxo[asdf]cvbn[mnop]qrst", true,
			[]string{"oxxo", "cvbn", "qrst"},
			[]string{"asdf", "mnop"}},
	}

	for ncase, item := range cases {
		ip, err := New(item.addr)
		if err != nil {
			t.Errorf("[Case %d] Error parsing '%s'", ncase, item.addr)
		}
		if len(item.supernets) != len(ip.Supernets) {
			t.Errorf("[Case %d] Found %d supernets in '%s' (expected %d)",
				ncase, len(ip.Supernets), item.addr, len(item.supernets))
		}
		if len(item.hypernets) != len(ip.Hypernets) {
			t.Errorf("[Case %d] Found %d hypernets in '%s' (expected %d)",
				ncase, len(ip.Hypernets), item.addr, len(item.hypernets))
		}
		if item.is_tls && !ip.IsTLS() {
			t.Errorf("[Case %d] Not identified as TLS, '%s'",
				ncase, item.addr)
		} else if !item.is_tls && ip.IsTLS() {
			t.Errorf("[Case %d] Incorrectly identified as TLS, '%s'",
				ncase, item.addr)
		}

	}

}

func TestSSL(t *testing.T) {
	cases := [...]struct {
		addr   string
		is_ssl bool
	}{
		{"aba[bab]xyz", true},
		{"xyx[xyx]xyx", false},
		{"aaa[kek]eke", true},
		{"zazbz[bzb]cdb", true},
	}

	for ncase, item := range cases {
		ip, err := New(item.addr)
		if err != nil {
			t.Errorf("[Case %d] Error parsing '%s'", ncase, item.addr)
		}
		if item.is_ssl && !ip.IsSSL() {
			t.Errorf("[Case %d] Not identified as SSL, '%s'",
				ncase, item.addr)
		} else if !item.is_ssl && ip.IsSSL() {
			t.Errorf("[Case %d] Incorrectly identified as SSL, '%s'",
				ncase, item.addr)
		}

	}

}
