package types

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/QOSGroup/cassini/log"
)

// BytesInt64 Int64 转换
func BytesInt64(bs []byte) (x int64, err error) {
	buf := bytes.NewBuffer(bs)
	err = binary.Read(buf, binary.BigEndian, &x)
	return x, err
}

// Int64Bytes int64 与 byte 数组转换
func Int64Bytes(in int64) []byte {
	var ret = bytes.NewBuffer([]byte{})
	err := binary.Write(ret, binary.BigEndian, in)
	if err != nil {
		log.Infof("Int2Byte error:%s", err.Error())
		return nil
	}

	return ret.Bytes()
}


// ParseAddrs parse protocol and addrs
func ParseAddrs(address string) (protocol string, addrs []string) {
	addrs = strings.SplitN(address, "://", 2)
	if len(addrs) == 2 {
		protocol = addrs[0]
		protocol = strings.TrimSpace(protocol)
		a := addrs[1]
		addrs = strings.Split(a, ",")
	} else {
		protocol, addrs = "", []string{}
	}

	return
}
