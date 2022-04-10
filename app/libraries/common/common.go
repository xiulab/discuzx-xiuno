package common

import (
	"encoding/binary"
	"github.com/gogf/gf/util/gconv"
	"math/rand"
	"net"
	"time"
)

// IP2Long IP 转整型
func IP2Long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

// Long2IP 整型转 IP
func Long2IP(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}

// GetRandomString 生成指定长度随机字符串
func GetRandomString(t string, l int) string {
	pool := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	switch t {

	case "alpha":
		pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	case "alnum":
		pool = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	case "numeric":
		pool = "0123456789"

	case "nozero":
		pool = "123456789"

	case "hex":
		pool = "0123456789abcdefABCDEF"
	}

	bytes := []byte(pool)
	result := []byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}

// FixGID 修正自用义用户组
func FixGID(val interface{}) int {
	gid := gconv.Int(val)

	// 系统组 对应关系
	// Dx  xn
	// 1 -> 1 管理员
	// 2 -> 2 超级版主
	// 3 -> 4 版主
	// 4 -> 7 (禁止发言->禁止)
	// 5 -> 7 (禁止访问->禁止)
	// 6 -> 7 (用户IP禁止->禁止)
	// 7 -> 0 游客
	// 8 -> 6 等待验证
	// 其余所有的组都 + 100

	switch gid {
	case 1:
	case 2:
	case 3:
		gid = 4
	case 4, 5, 6:
		gid = 7
	case 7:
		gid = 0
	case 8:
		gid = 6
	default:
		gid += 100
	}

	return gid
}
