package dns

import (
	"encoding/hex"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"

	microlog "../../log"
)

// DnsCache <-- use map save history
type DnsCache struct {
	dns       sync.Map
	DNSServer string
}

// DNSv4 <-- dns for ipv4
func DNSv4(DNSServer, domain string) (DNS []string, success bool) {
	// +------------------------------+
	// |             id               |  16bit
	// +------------------------------+
	// |qr|opcpde|aa|tc|rd|ra|z|rcode |
	// +------------------------------+
	// |          QDCOUNT             |
	// +------------------------------+
	// |          ancount             |
	// +------------------------------+
	// |          nscount             |
	// +------------------------------+
	// |          arcount             |
	// +------------------------------+

	// • ID：这是由生成DNS查询的程序指定的16位的标志符。该标志符也被随后的应答报文所用，申请者利用这个标志将应答和原来的请求对应起来。

	// • QR：该字段占1位，用以指明DNS报文是请求（0）还是应答（1）。
	// • OPCODE：该字段占4位，用于指定查询的类型。值为0表示标准查询，值为1表示逆向查询，值为2表示查询服务器状态，值为3保留，值为4表示通知，值为5表示更新报文，值6～15的留为新增操作用。
	// • AA：该字段占1位，仅当应答时才设置。值为1，即意味着正应答的域名服务器是所查询域名的管理机构或者说是被授权的域名服务器。
	// • TC：该字段占1位，代表截断标志。如果报文长度比传输通道所允许的长而被分段，该位被设为1。
	// • RD：该字段占1位，是可选项，表示要求递归与否。如果为1，即意味 DNS解释器要求DNS服务器使用递归查询。

	// • RA：该字段占1位，代表正在应答的域名服务器可以执行递归查询，该字段与查询段无关。
	// • Z：该字段占3位，保留字段，其值在查询和应答时必须为0。
	// • RCODE：该字段占4位，该字段仅在DNS应答时才设置。用以指明是否发生了错误。
	// 允许取值范围及意义如下：
	// 0：无错误情况，DNS应答表现为无错误。
	// 1：格式错误，DNS服务器不能解释应答。
	// 2：严重失败，因为名字服务器上发生了一个错误，DNS服务器不能处理查询。
	// 3：名字错误，如果DNS应答来自于授权的域名服务器，意味着DNS请求中提到的名字不存在。
	// 4：没有实现。DNS服务器不支持这种DNS请求报文。
	// 5：拒绝，由于安全或策略上的设置问题，DNS名字服务器拒绝处理请求。
	// 6 ～15 ：留为后用。

	// • QDCOUNT：该字段占16位，指明DNS查询段中的查询问题的数量。
	// • ANCOUNT：该字段占16位，指明DNS应答段中返回的资源记录的数量，在查询段中该值为0。
	// • NSCOUNT：该字段占16位，指明DNS应答段中所包括的授权域名服务器的资源记录的数量，在查询段中该值为0。
	// • ARCOUNT：该字段占16位，指明附加段里所含资源记录的数量，在查询段中该值为0。
	// (2）DNS正文段
	// 在DNS报文中，其正文段封装在图7-42所示的DNS报文头内。DNS有四类正文段：查询段、应答段、授权段和附加段。

	// id -> 2byte
	// qr -> 1bit opcpde -> 4bit aa -> 1bit tc -> 1bit rd -> 1bit  => sum 1byte
	//  ra -> 1bit z -> 3bit rcode -> 4bit => sum 1byte
	// QDCOUNT -> 2byte
	// ANCOUNT -> 2byte
	// NSCOUNT -> 2byte
	// ARCOUNT -> 2byte
	header := []byte{0x01, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	var domainSet []byte
	for _, domain := range strings.Split(domain, ".") {
		domainSet = append(domainSet, byte(len(domain)))
		domainSet = append(domainSet, []byte(domain)...)
	}

	// append domain and qType And QClass
	// domain []byte(domainSet), 0x00
	// qType 0x00,0x01
	// qClass 0x00,0x01
	domainSetAndQTypeAndQClass := append(domainSet, 0x00, 0x00, 0x01, 0x00, 0x01)

	all := append(header, domainSetAndQTypeAndQClass...)

	conn, err := net.Dial("udp", DNSServer)
	if err != nil {
		microlog.Debug(err)
		return []string{}, false
	}
	defer conn.Close()

	// log.Println(all, len(all))

	var b [1024]byte
	_, _ = conn.Write(all)
	// log.Println("write")
	// var b [1024]byte
	n, _ := conn.Read(b[:])

	if b[3]&1 != 0 {
		// log.Println("no such name")
		return []string{}, false
	}
	// log.Println("header", b[0:12], "qr+opcode+aa+tc+rd:", b[2:3], "ra+z+rcode:", b[3], "rcode:", b[3]&1, "....", b[3]&2, b[3]&4, b[3]&8)

	// • QDCOUNT：该字段占16位，指明DNS查询段中的查询问题的数量。
	// • ANCOUNT：该字段占16位，指明DNS应答段中返回的资源记录的数量
	// • NSCOUNT：该字段占16位，指明DNS应答段中所包括的授权域名服务器的资源记录的数量
	// • ARCOUNT：该字段占16位，指明附加段里所含资源记录的数量
	// log.Println("QDCOUNT:", b[4], b[5], "ANCOUNT:", b[6], b[7], "NSCOUNT:", b[8], b[9], "ARCOUNT:", b[10], b[11])
	anCount := int(b[6])<<8 + int(b[7])

	bHeader := b[12:n]
	// bf, index := 1, int(b_header[0])+1
	index := int(bHeader[0]) + 1
	for {
		// log.Println("bf", bf, "index", index, string(b_header[bf:index]))
		if bHeader[index] == 0 {
			break
		}
		// bf = index
		index = index + int(bHeader[index]) + 1
	}
	// index = index + 1
	// // log.Println("type", b_header[index:index+2])
	// // log.Println("class", b_header[index+2:index+4])
	// index = index + 4
	answer := bHeader[index+5:]
	answerIndex := 0
	var dns []string
	for i := 0; i < anCount; i++ {
		if answer[answerIndex]&128 == 128 && answer[answerIndex]&64 == 64 {
			answerIndex += 2
		}
		if int16(answer[answerIndex])<<8+int16(answer[answerIndex+1]) != 0x01 {
			answerIndex += 8
			answerIndex += 2 + int(answer[answerIndex])<<8 + int(answer[answerIndex+1])
		} else {
			answerIndex += 8
			if int16(answer[answerIndex])<<8+int16(answer[answerIndex+1]) == 4 {
				answerIndex += 2
				dns = append(dns, strconv.Itoa(int(answer[answerIndex]))+"."+
					strconv.Itoa(int(answer[answerIndex+1]))+"."+
					strconv.Itoa(int(answer[answerIndex+2]))+"."+
					strconv.Itoa(int(answer[answerIndex+3])))
				answerIndex += 4
			} else if int16(answer[answerIndex])<<8+int16(answer[answerIndex+1]) == 6 {
				answerIndex += 2
				hexDNS := hex.EncodeToString(answer[answerIndex : answerIndex+16])
				microlog.Debug(hexDNS[0:4] + ":" + hexDNS[4:8] + ":" + hexDNS[8:12] + ":" + hexDNS[12:16] +
					":" + hexDNS[16:20] + ":" + hexDNS[20:24] + ":" + hexDNS[24:28] + ":" + hexDNS[28:32])
				dns = append(dns, hexDNS[0:4]+":"+hexDNS[4:8]+":"+hexDNS[8:12]+":"+hexDNS[12:16]+
					":"+hexDNS[16:20]+":"+hexDNS[20:24]+":"+hexDNS[24:28]+":"+hexDNS[28:32])

				answerIndex += 16
			}
			// log.Println(answer[answerIndex], answer[answerIndex+1], answer[answerIndex+2], answer[answerIndex+3])
			// log.Println(strconv.Itoa(int(answer[answerIndex])) + "." + strconv.Itoa(int(answer[answerIndex+1])) + "." + strconv.Itoa(int(answer[answerIndex+2])) + "." + strconv.Itoa(int(answer[answerIndex+3])))

		}
	}

	// log.Println(dns)

	// log.Println(b[bytes.Index(b[:n], []byte{192, 12})+2+2+2+4 : n])
	// ip := b[bytes.Index(b[:n], []byte{192, 12})+2+2+2+4 : n]
	// log.Println("ip:", strconv.Itoa(int(ip[2]))+"."+strconv.Itoa(int(ip[3]))+"."+strconv.Itoa(int(ip[4]))+"."+strconv.Itoa(int(ip[5])))
	// log.Println(strconv.Itoa(int(b[n-4])) + "." + strconv.Itoa(int(b[n-3])) + "." + strconv.Itoa(int(b[n-2])) + "." + strconv.Itoa(int(b[n-1])))
	if len(dns) != 0 {
		//microlog.Debug(domain, dns)
		return dns, true
	}
	//microlog.Debug(domain, dns)
	return dns, false
}

// Match <--
func (dnscache *DnsCache) Match(host, hostTemplate string, cidrmatch func(string) bool) (isMatched bool) {
	var isMatch bool
	if _, exist := dnscache.dns.Load(host); exist == false {
		if hostTemplate != "ip" {
			// ip, err := net.LookupHost(host)
			ip, isSuccess := DNSv4(dnscache.DNSServer, host)
			// microlog.Debug(host, ip, isSuccess)
			if isSuccess == true {
				isMatch = cidrmatch(ip[0])
				// microlog.Debug(isMatch, ip[0])
			} else {
				isMatch = false
			}
		} else {
			isMatch = cidrmatch(host)
		}
		// if len(socks5Server.dns) > 10000 {
		// 	i := 0
		// 	for key := range socks5Server.dns {
		// 		delete(socks5Server.dns, key)
		// 		i++
		// 		if i > 0 {
		// 			break
		// 		}
		// 	}
		// }
		dnscache.dns.Store(host, isMatch)
		// fmt.Println(runtime.NumGoroutine(), host, isMatch)
		microlog.Debug(runtime.NumGoroutine(), host, isMatch)
	} else {
		isMatchTemp, _ := dnscache.dns.Load(host)
		isMatch = isMatchTemp.(bool)
		// fmt.Println(runtime.NumGoroutine(), "use cache", host, isMatch)
		microlog.Debug(runtime.NumGoroutine(), "use cache", host, isMatch)
	}
	return isMatch
}

// -------------------------------old---------------------
// func DNSv4(DNSServer, domain string) (domainIP string, success bool) {
// 	// +------------------------------+
// 	// |             id               |  16bit
// 	// +------------------------------+
// 	// |qr|opcpde|aa|tc|rd|ra|z|rcode |
// 	// +------------------------------+
// 	// |          QDCOUNT             |
// 	// +------------------------------+
// 	// |          ancount             |
// 	// +------------------------------+
// 	// |          nscount             |
// 	// +------------------------------+
// 	// |          arcount             |
// 	// +------------------------------+

// 	// • ID：这是由生成DNS查询的程序指定的16位的标志符。该标志符也被随后的应答报文所用，申请者利用这个标志将应答和原来的请求对应起来。

// 	// • QR：该字段占1位，用以指明DNS报文是请求（0）还是应答（1）。
// 	// • OPCODE：该字段占4位，用于指定查询的类型。值为0表示标准查询，值为1表示逆向查询，值为2表示查询服务器状态，值为3保留，值为4表示通知，值为5表示更新报文，值6～15的留为新增操作用。
// 	// • AA：该字段占1位，仅当应答时才设置。值为1，即意味着正应答的域名服务器是所查询域名的管理机构或者说是被授权的域名服务器。
// 	// • TC：该字段占1位，代表截断标志。如果报文长度比传输通道所允许的长而被分段，该位被设为1。
// 	// • RD：该字段占1位，是可选项，表示要求递归与否。如果为1，即意味 DNS解释器要求DNS服务器使用递归查询。

// 	// • RA：该字段占1位，代表正在应答的域名服务器可以执行递归查询，该字段与查询段无关。
// 	// • Z：该字段占3位，保留字段，其值在查询和应答时必须为0。
// 	// • RCODE：该字段占4位，该字段仅在DNS应答时才设置。用以指明是否发生了错误。
// 	// 允许取值范围及意义如下：
// 	// 0：无错误情况，DNS应答表现为无错误。
// 	// 1：格式错误，DNS服务器不能解释应答。
// 	// 2：严重失败，因为名字服务器上发生了一个错误，DNS服务器不能处理查询。
// 	// 3：名字错误，如果DNS应答来自于授权的域名服务器，意味着DNS请求中提到的名字不存在。
// 	// 4：没有实现。DNS服务器不支持这种DNS请求报文。
// 	// 5：拒绝，由于安全或策略上的设置问题，DNS名字服务器拒绝处理请求。
// 	// 6 ～15 ：留为后用。

// 	// • QDCOUNT：该字段占16位，指明DNS查询段中的查询问题的数量。
// 	// • ANCOUNT：该字段占16位，指明DNS应答段中返回的资源记录的数量，在查询段中该值为0。
// 	// • NSCOUNT：该字段占16位，指明DNS应答段中所包括的授权域名服务器的资源记录的数量，在查询段中该值为0。
// 	// • ARCOUNT：该字段占16位，指明附加段里所含资源记录的数量，在查询段中该值为0。
// 	// (2）DNS正文段
// 	// 在DNS报文中，其正文段封装在图7-42所示的DNS报文头内。DNS有四类正文段：查询段、应答段、授权段和附加段。

// 	// id -> 2byte
// 	// qr -> 1bit opcpde -> 4bit aa -> 1bit tc -> 1bit rd -> 1bit  => sum 1byte
// 	//  ra -> 1bit z -> 3bit rcode -> 4bit => sum 1byte
// 	// QDCOUNT -> 2byte
// 	// ANCOUNT -> 2byte
// 	// NSCOUNT -> 2byte
// 	// ARCOUNT -> 2byte
// 	header := []byte{0x01, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

// 	var domainSet []byte
// 	for _, domain := range strings.Split(domain, ".") {
// 		domainSet = append(domainSet, byte(len(domain)))
// 		domainSet = append(domainSet, []byte(domain)...)
// 	}

// 	// append domain and qType And QClass
// 	// domain []byte(domainSet), 0x00
// 	// qType 0x00,0x01
// 	// qClass 0x00,0x01
// 	domainSetAndQTypeAndQClass := append([]byte(domainSet), 0x00, 0x00, 0x01, 0x00, 0x01)

// 	all := append(header, domainSetAndQTypeAndQClass...)

// 	conn, err := net.Dial("udp", DNSServer)
// 	if err != nil {
// 		log.Println(err)
// 		return "", false
// 	}
// 	defer conn.Close()

// 	// log.Println(all, len(all))

// 	var b [1024]byte
// 	conn.Write(all)
// 	// log.Println("write")
// 	// var b [1024]byte
// 	n, _ := conn.Read(b[:])
// 	// log.Println("header", b[0:12], "qr+opcode+aa+tc+rd:", b[2:3], "ra+z+rcode:", b[3], "rcode:", b[3]&1, "....", b[3]&2, b[3]&4, b[3]&8)
// 	if b[3]&1 != 0 {
// 		// log.Println("no such name")
// 		return "", false
// 	}
// 	// log.Println(b[bytes.Index(b[:n], []byte{192, 12})+2+2+2+4 : n])
// 	// ip := b[bytes.Index(b[:n], []byte{192, 12})+2+2+2+4 : n]
// 	// log.Println("ip:", strconv.Itoa(int(ip[2]))+"."+strconv.Itoa(int(ip[3]))+"."+strconv.Itoa(int(ip[4]))+"."+strconv.Itoa(int(ip[5])))
// 	// log.Println(strconv.Itoa(int(b[n-4])) + "." + strconv.Itoa(int(b[n-3])) + "." + strconv.Itoa(int(b[n-2])) + "." + strconv.Itoa(int(b[n-1])))
// 	return strconv.Itoa(int(b[n-4])) + "." + strconv.Itoa(int(b[n-3])) + "." + strconv.Itoa(int(b[n-2])) + "." + strconv.Itoa(int(b[n-1])), true
// }
