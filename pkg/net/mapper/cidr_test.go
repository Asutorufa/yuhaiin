package mapper

import (
	"math"
	"net"
	"testing"
)

func TestCidrMatch_Inset(t *testing.T) {
	cidrMatch := NewCidrMapper()
	if err := cidrMatch.Insert("10.2.2.1/18", "testIPv4"); err != nil {
		t.Error(err)
	}
	if err := cidrMatch.Insert("2001:0db8:0000:0000:1234:0000:0000:9abc/32", "testIPv6"); err != nil {
		t.Error(err)
	}
	cidrMatch.Insert("127.0.0.1/8", "testlocal")
	cidrMatch.Insert("ff::/16", "testV6local")
	cidrMatch.v6CidrTrie.Print()
	cidrMatch.v4CidrTrie.Print()
	testIPv4 := "10.2.2.1"
	testIPv4b := "100.2.2.1"
	testIPv6 := "2001:0db8:0000:0000:1234:0000:0000:9abc"
	testIPv6b := "3001:0db8:0000:0000:1234:0000:0000:9abc"
	t.Log(cidrMatch.Search(testIPv4))    // true
	t.Log(cidrMatch.Search(testIPv6))    // true
	t.Log(cidrMatch.Search("127.1.1.1")) // true
	t.Log(cidrMatch.Search("129.1.1.1")) // false
	t.Log(cidrMatch.Search("ff:ff::"))   // true
	t.Log(cidrMatch.Search(testIPv4b))   // false
	t.Log(cidrMatch.Search(testIPv6b))   // false
	t.Log()
	t.Log(cidrMatch.Search(testIPv4))    // true
	t.Log(cidrMatch.Search(testIPv6))    // true
	t.Log(cidrMatch.Search("127.1.1.1")) // true
	t.Log(cidrMatch.Search("129.1.1.1")) // false
	t.Log(cidrMatch.Search("ff:ff::"))   // true
	t.Log(cidrMatch.Search(testIPv4b))   // false
	t.Log(cidrMatch.Search(testIPv6b))   // false
}

// BenchmarkCidrMatch_Search-4 9119133 130.6 ns/op 16 B/op 1 allocs/op
func BenchmarkCidrMatch_Search(b *testing.B) {
	b.StopTimer()

	cidrMatch := NewCidrMapper()
	if err := cidrMatch.Insert("10.2.2.1/18", "testIPv4"); err != nil {
		b.Error(err)
	}
	if err := cidrMatch.Insert("2001:0db8:0000:0000:1234:0000:0000:9abc/32", "testIPv6"); err != nil {
		b.Error(err)
	}
	testIPv4 := "10.2.2.1"
	//testIPv4b := "100.2.2.1"
	//testIPv6 := "2001:0db8:0000:0000:1234:0000:0000:9abc"
	testIPv6b := "3001:0db8:0000:0000:1234:0000:0000:9abc"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			cidrMatch.Search(testIPv4)
		} else {
			cidrMatch.Search(testIPv6b)
		}
	}
}

func TestIpToInt(t *testing.T) {
	t.Log([]byte(net.ParseIP("127.0.0.1").To4()))
	t.Log(ipv4toInt(net.ParseIP("127.0.0.1")))
	t.Log(ipv4toInt2(net.ParseIP("127.0.0.1").To4()))
	t.Log(ipv4toInt(net.ParseIP("0.0.0.1")))
	t.Log(ipv4toInt2(net.ParseIP("0.0.0.1").To4()))
	t.Log(ipv4toInt(net.ParseIP("255.255.255.255")))
	t.Log(ipv4toInt2(net.ParseIP("255.255.255.255").To4()))
	t.Log(ipv6toInt(net.ParseIP("ff::ff")))
	t.Log(ipv6toInt(net.ParseIP("::ff")))
}

// 297 293
func BenchmarkIpv4toInt(b *testing.B) {
	str := net.ParseIP("0.0.255.255")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ipv4toInt(str)
	}
}

// 827 821
// 729 752
func BenchmarkIpv6toInt(b *testing.B) {
	str := net.ParseIP("ffff::ffff")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ipv6toInt2(str)
	}
}

func TestIpToCidr(t *testing.T) {
	ip, ipNet, err := net.ParseCIDR("127.0.0.1/28")
	if err != nil {
		t.Error(err)
	}
	t.Log(ip)
	t.Log(ipNet.Mask.Size())
	ip, ipNet, err = net.ParseCIDR("ff::ff/64")
	if err != nil {
		t.Error(err)
	}
	t.Log(ip.To4())
	t.Log(ipNet.Mask.Size())
}

func TestTo6(t *testing.T) {
	// Addresses in this group consist of an 80-bit prefix of zeros,
	// the next 16 bits are ones, and the remaining,
	// least-significant 32 bits contain the IPv4 address.
	// For example,
	// ::ffff:192.0.2.128 represents the IPv4 address 192.0.2.128.
	// Another format, called "IPv4-compatible IPv6 address",
	// is ::192.0.2.128; however, this method is deprecated.
	t.Log(ipv6toInt(net.ParseIP("127.0.0.1")))
	t.Log(ipv6toInt2(net.ParseIP("127.0.0.1")))
	t.Log(ipv6toInt(net.ParseIP("::127.0.0.1")))     //deprecated
	t.Log([]byte(net.ParseIP("::127.0.0.1").To16())) //deprecated
	t.Log([]byte(net.ParseIP("127.0.0.1").To16()))
	t.Log(net.ParseIP("::ffff:a:b"))
	//01111111 00000000 00000000 00000001

	_, ips, err := net.ParseCIDR("127.0.0.1/28")
	if err != nil {
		t.Error(err)
	}
	t.Log(ips.IP.To16())
	if len(ips.IP) == net.IPv4len {
		size, _ := ips.Mask.Size()
		t.Log(size + 96)
		t.Log(ipv6toInt(ips.IP)[:size+96])
		t.Log([]byte(ips.IP.Mask(ips.Mask).To16()))
	}
}

// 2102 ns/op,2106 ns/op
// BenchmarkSingleTrie-4 9823910 128.0 ns/op 16 B/op  1 allocs/op
func BenchmarkSingleTrie(b *testing.B) {
	m := NewCidrMapper()
	if err := m.singleInsert("127.0.0.1/28", "localhost"); err != nil {
		b.Error(err)
	}
	m.singleInsert("ff::/56", "")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			m.singleSearch("127.0.0.0")
		} else {
			m.singleSearch("ff::")
		}
	}
}

func TestLog(t *testing.T) {
	t.Log(math.Log2(128))
}
