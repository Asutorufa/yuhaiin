// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config/config.proto

package config

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//	BlackIcon          bool   `json:"black_icon"`
//	IsDNSOverHTTPS     bool   `json:"is_dns_over_https"`
//	DNSAcrossProxy     bool   `json:"dns_across_proxy"`
//	DnsServer          string `json:"dnsServer"`
//	DnsSubNet          string `json:"dns_sub_net"`
//	Bypass             bool   `json:"bypass"`
//	HttpProxyAddress   string `json:"httpProxyAddress"`
//	Socks5ProxyAddress string `json:"socks5ProxyAddress"`
//	RedirProxyAddress  string `json:"redir_proxy_address"`
//	BypassFile         string `json:"bypassFile"`
//	SsrPath            string `json:"ssrPath"`
//	LocalAddress       string `json:"localAddress"`
//	LocalPort          string `json:"localPort"`
type Setting struct {
	BlackIcon            bool     `protobuf:"varint,1,opt,name=BlackIcon,json=black_icon,proto3" json:"BlackIcon,omitempty"`
	IsDNSOverHTTPS       bool     `protobuf:"varint,2,opt,name=IsDNSOverHTTPS,json=is_dns_over_https,proto3" json:"IsDNSOverHTTPS,omitempty"`
	DNSAcrossProxy       bool     `protobuf:"varint,3,opt,name=DNSAcrossProxy,json=dns_across_proxy,proto3" json:"DNSAcrossProxy,omitempty"`
	DnsServer            string   `protobuf:"bytes,4,opt,name=DnsServer,json=dnsServer,proto3" json:"DnsServer,omitempty"`
	DnsSubNet            string   `protobuf:"bytes,5,opt,name=DnsSubNet,json=dns_sub_net,proto3" json:"DnsSubNet,omitempty"`
	Bypass               bool     `protobuf:"varint,6,opt,name=Bypass,json=bypass,proto3" json:"Bypass,omitempty"`
	HttpProxyAddress     string   `protobuf:"bytes,7,opt,name=HttpProxyAddress,json=httpProxyAddress,proto3" json:"HttpProxyAddress,omitempty"`
	Socks5ProxyAddress   string   `protobuf:"bytes,8,opt,name=Socks5ProxyAddress,json=socks5ProxyAddress,proto3" json:"Socks5ProxyAddress,omitempty"`
	RedirProxyAddress    string   `protobuf:"bytes,9,opt,name=RedirProxyAddress,json=redir_proxy_address,proto3" json:"RedirProxyAddress,omitempty"`
	BypassFile           string   `protobuf:"bytes,10,opt,name=BypassFile,json=bypassFile,proto3" json:"BypassFile,omitempty"`
	SsrPath              string   `protobuf:"bytes,11,opt,name=SsrPath,json=ssrPath,proto3" json:"SsrPath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Setting) Reset()         { *m = Setting{} }
func (m *Setting) String() string { return proto.CompactTextString(m) }
func (*Setting) ProtoMessage()    {}
func (*Setting) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc332a44e926b360, []int{0}
}

func (m *Setting) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Setting.Unmarshal(m, b)
}
func (m *Setting) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Setting.Marshal(b, m, deterministic)
}
func (m *Setting) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Setting.Merge(m, src)
}
func (m *Setting) XXX_Size() int {
	return xxx_messageInfo_Setting.Size(m)
}
func (m *Setting) XXX_DiscardUnknown() {
	xxx_messageInfo_Setting.DiscardUnknown(m)
}

var xxx_messageInfo_Setting proto.InternalMessageInfo

func (m *Setting) GetBlackIcon() bool {
	if m != nil {
		return m.BlackIcon
	}
	return false
}

func (m *Setting) GetIsDNSOverHTTPS() bool {
	if m != nil {
		return m.IsDNSOverHTTPS
	}
	return false
}

func (m *Setting) GetDNSAcrossProxy() bool {
	if m != nil {
		return m.DNSAcrossProxy
	}
	return false
}

func (m *Setting) GetDnsServer() string {
	if m != nil {
		return m.DnsServer
	}
	return ""
}

func (m *Setting) GetDnsSubNet() string {
	if m != nil {
		return m.DnsSubNet
	}
	return ""
}

func (m *Setting) GetBypass() bool {
	if m != nil {
		return m.Bypass
	}
	return false
}

func (m *Setting) GetHttpProxyAddress() string {
	if m != nil {
		return m.HttpProxyAddress
	}
	return ""
}

func (m *Setting) GetSocks5ProxyAddress() string {
	if m != nil {
		return m.Socks5ProxyAddress
	}
	return ""
}

func (m *Setting) GetRedirProxyAddress() string {
	if m != nil {
		return m.RedirProxyAddress
	}
	return ""
}

func (m *Setting) GetBypassFile() string {
	if m != nil {
		return m.BypassFile
	}
	return ""
}

func (m *Setting) GetSsrPath() string {
	if m != nil {
		return m.SsrPath
	}
	return ""
}

func init() {
	proto.RegisterType((*Setting)(nil), "yuhaiin.config.Setting")
}

func init() {
	proto.RegisterFile("config/config.proto", fileDescriptor_cc332a44e926b360)
}

var fileDescriptor_cc332a44e926b360 = []byte{
	// 340 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0x51, 0x6b, 0xea, 0x30,
	0x14, 0xc7, 0xf1, 0x7a, 0x6f, 0x6b, 0x8f, 0x20, 0x1a, 0xe1, 0x92, 0x87, 0x4d, 0x64, 0x63, 0xe0,
	0xf6, 0x50, 0x1f, 0xc6, 0x3e, 0x80, 0x22, 0x43, 0x5f, 0x9c, 0x58, 0x9f, 0xf6, 0x52, 0x92, 0x34,
	0xda, 0xa0, 0x4b, 0x4a, 0x4e, 0x2a, 0xf3, 0xa3, 0xec, 0xdb, 0x8e, 0xa6, 0x6e, 0xe0, 0xf6, 0x54,
	0xce, 0xef, 0xff, 0xeb, 0x49, 0xc2, 0x1f, 0xfa, 0xc2, 0xe8, 0xad, 0xda, 0x8d, 0xeb, 0x4f, 0x5c,
	0x58, 0xe3, 0x0c, 0xe9, 0x9c, 0xca, 0x9c, 0x29, 0xa5, 0xe3, 0x9a, 0xde, 0x7c, 0x34, 0x21, 0x4c,
	0xa4, 0x73, 0x4a, 0xef, 0xc8, 0x35, 0x44, 0xd3, 0x03, 0x13, 0xfb, 0x85, 0x30, 0x9a, 0x36, 0x86,
	0x8d, 0x51, 0x6b, 0x0d, 0xbc, 0x02, 0xa9, 0x12, 0x46, 0x93, 0x7b, 0xe8, 0x2c, 0x70, 0xb6, 0x4c,
	0x5e, 0x8e, 0xd2, 0xce, 0x37, 0x9b, 0x55, 0x42, 0xff, 0x78, 0xa7, 0xa7, 0x30, 0xcd, 0x34, 0xa6,
	0xe6, 0x28, 0x6d, 0x9a, 0x3b, 0x57, 0x20, 0x19, 0x41, 0x67, 0xb6, 0x4c, 0x26, 0xc2, 0x1a, 0xc4,
	0x95, 0x35, 0xef, 0x27, 0xda, 0xf4, 0x6a, 0xb7, 0xf2, 0x98, 0xc7, 0x69, 0x51, 0x71, 0x72, 0x05,
	0xd1, 0x4c, 0x63, 0x22, 0xed, 0x51, 0x5a, 0xfa, 0x77, 0xd8, 0x18, 0x45, 0xeb, 0x28, 0xfb, 0x02,
	0x64, 0x50, 0xa7, 0x25, 0x5f, 0x4a, 0x47, 0xff, 0xf9, 0xb4, 0x5d, 0xad, 0xc0, 0x92, 0xa7, 0x5a,
	0x3a, 0xf2, 0x1f, 0x82, 0xe9, 0xa9, 0x60, 0x88, 0x34, 0xf0, 0xfb, 0x03, 0xee, 0x27, 0xf2, 0x00,
	0xdd, 0xb9, 0x73, 0x85, 0x3f, 0x7a, 0x92, 0x65, 0x56, 0x22, 0xd2, 0xd0, 0xff, 0xde, 0xcd, 0x7f,
	0x70, 0x12, 0x03, 0x49, 0x8c, 0xd8, 0xe3, 0xd3, 0x85, 0xdd, 0xf2, 0x36, 0xc1, 0x5f, 0x09, 0x89,
	0xa1, 0xb7, 0x96, 0x99, 0xb2, 0x17, 0x7a, 0xe4, 0xf5, 0xbe, 0xad, 0x82, 0xfa, 0x65, 0x29, 0x3b,
	0xfb, 0x03, 0x80, 0xfa, 0x8e, 0xcf, 0xea, 0x20, 0x29, 0x78, 0x11, 0xf8, 0x37, 0x21, 0x14, 0xc2,
	0x04, 0xed, 0x8a, 0xb9, 0x9c, 0xb6, 0x7d, 0x18, 0x62, 0x3d, 0x4e, 0xef, 0x5e, 0x6f, 0x77, 0xca,
	0xe5, 0x25, 0x8f, 0x85, 0x79, 0x1b, 0x4f, 0xb0, 0x74, 0xc6, 0x96, 0x5b, 0x36, 0x3e, 0x57, 0x78,
	0x2e, 0x96, 0x07, 0xbe, 0xd9, 0xc7, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x85, 0x3e, 0x83, 0x90,
	0xf0, 0x01, 0x00, 0x00,
}
