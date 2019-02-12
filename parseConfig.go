// go build server.go frame.go robot.go
package main

import (
//	"flag"
//	"log"
//	"fmt"
//	"net"

//	"io"
	"os"
	"encoding/hex"
	"strings"
	"strconv"
	"bufio"

//	"strconv"
//	"time"
//	"sync"
)

func readData() (error) {
	lines, err := readFile(*userData)
	if err != nil {
		Vln(2, "[open]", err)
		return err
	}

	grid2 := NewGrid()

	idx := 0
	for _, line := range lines {

		fields := strings.Split(line, "\t")
		if fields[0] == "" {
			continue
		}
		if strings.HasPrefix(fields[0], "#") {
			continue
		}

		Vln(6, "[dbg]", len(fields), fields)
		if strings.HasPrefix(fields[0], "!!") {
			readUser(fields)
			continue
		}

		var rid uint16 = 0x4286
		var C4 []byte = []byte{0xFF, 0xFF, 0xFF, 0xFF}
		var wing uint8 = 0
		var wingLv []byte = []byte{0x00, 0x00, 0x00, 0x00}
		var Lv uint8 = 13
		var exp uint32 = 12345
		var sess uint32 = 23333
		var skill []byte = []byte{0x00, 0x00, 0x00, 0x00}
		var polish uint16 = 0
		var color []uint16 = []uint16{0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000}
		var coat []uint32 = []uint32{0x00000000, 0x00000000, 0x00000000}


		switch len(fields) {
		default:
			fallthrough
		case 18: // 紋章
			tmp, _ := strconv.ParseUint(fields[17], 16, 32)
			coat[2] = uint32(tmp)
			fallthrough
		case 17:
			tmp, _ := strconv.ParseUint(fields[16], 16, 32)
			coat[1] = uint32(tmp)
			fallthrough
		case 16:
			tmp, _ := strconv.ParseUint(fields[15], 16, 32)
			coat[0] = uint32(tmp)
			fallthrough

		case 15: // 塗裝
			tmp, _ := strconv.ParseUint(fields[14], 16, 16)
			color[5] = uint16(tmp)
			fallthrough
		case 14:
			tmp, _ := strconv.ParseUint(fields[13], 16, 16)
			color[4] = uint16(tmp)
			fallthrough
		case 13:
			tmp, _ := strconv.ParseUint(fields[12], 16, 16)
			color[3] = uint16(tmp)
			fallthrough
		case 12:
			tmp, _ := strconv.ParseUint(fields[11], 16, 16)
			color[2] = uint16(tmp)
			fallthrough
		case 11:
			tmp, _ := strconv.ParseUint(fields[10], 16, 16)
			color[1] = uint16(tmp)
			fallthrough
		case 10:
			tmp, _ := strconv.ParseUint(fields[9], 16, 16)
			color[0] = uint16(tmp)
			fallthrough

		case 9: // 拋光
			tmp, _ := strconv.ParseUint(fields[8], 16, 16)
			polish = uint16(tmp)
			fallthrough
		case 8:
			skill, _ = hex.DecodeString(fields[7])
			fallthrough
		case 7:
			tmp, _ := strconv.ParseUint(fields[6], 10, 32)
			sess = uint32(tmp)
			fallthrough
		case 6:
			tmp, _ := strconv.ParseUint(fields[5], 10, 32)
			exp = uint32(tmp)
			fallthrough
		case 5:
			tmp, _ := strconv.ParseUint(fields[4], 10, 8)
			if tmp >= 1 && tmp <= 13 {
				Lv = uint8(tmp)
			}
			fallthrough
		case 4:
			wingLv, _ = hex.DecodeString(fields[3])
			fallthrough
		case 3:
			tmp, _ := strconv.ParseUint(fields[2], 10, 8)
			wing = uint8(tmp)
			fallthrough
		case 2:
			C4, _ = hex.DecodeString(fields[1])
			fallthrough
		case 1:
			tmp, _ := hex.DecodeString(fields[0])
			if len(tmp) == 2 {
				rid = (uint16(tmp[1]) << 8) | uint16(tmp[0])
			}

		case 0:
			Vln(1, "[open]?!!")
			continue
		}

		bot := NewBot(rid)
		bot.C4 = C4
		bot.Lv = Lv
		bot.Exp = exp
		bot.Sess = sess
		bot.Wing = wing
		bot.WingLv = wingLv
		bot.Skill = skill
		bot.Polish = polish
		bot.Color = color
		bot.Coat = coat


		Vf(5, "[dbg][open]%04X, %04X, %d, %04X\n", rid, C4, wing, wingLv)
		//Vf(7, "[dbg][open]%v, %X", bot, bot.GetBytes(idx))
		idx += 1

		grid2.Add(bot)
	}

	grid2.BuildCached()
	grid2.BuildCachedAll()

	grid = grid2

	Vln(4, "[dbg][grid]", len(grid.Robot), len(grid.buf))
	Vln(4, "[dbg][user]", user)
	return nil
}

func readUser(d []string) {
	if len(d) < 3 {
		return
	}

	val := d[2]
	switch d[1] {
	case "Name":
		user.SetName(val)

	case "GP":
		tmp, err := strconv.ParseUint(val, 10, 32)
		if err == nil {
			user.Mx.Lock()
			user.GP = uint32(tmp)
			user.Mx.Unlock()
		}

	case "GO":
		tmp, err := strconv.ParseUint(val, 10, 8)
		if err == nil {
			if tmp < 37 && tmp > 0 {
				user.Mx.Lock()
				user.GO = int(tmp)
				user.Mx.Unlock()
			}
		}

	case "SearchID":
		var rid uint16 = 0x4286
		tmp, _ := hex.DecodeString(val)
		if len(tmp) == 2 {
			rid = (uint16(tmp[1]) << 8) | uint16(tmp[0])
		}
		user.Mx.Lock()
		user.SearchID = rid
		user.Mx.Unlock()

	case "SearchExp":
		tmp, err := strconv.ParseUint(val, 10, 32)
		if err == nil {
			user.Mx.Lock()
			user.SearchExp = uint32(tmp)
			user.Mx.Unlock()
		}

	case "PageCount":
		tmp, err := strconv.ParseUint(val, 10, 32)
		if err == nil {
			user.Mx.Lock()
			user.PageCount = int(tmp)
			user.Mx.Unlock()
		}

	default:
		return
	}
}

func readEggPool() (error) {
	lines, err := readFile(*eggPoolData)
	if err != nil {
		Vln(2, "[open]", err)
		return err
	}

	eggPool2 := NewEggPool()
	for _, line := range lines {

		fields := strings.Split(line, "\t")
		if fields[0] == "" {
			continue
		}
		if strings.HasPrefix(fields[0], "#") {
			continue
		}

		var rid uint16 = 0x4286
		var C uint8 = 3
		var P int = 0

		switch len(fields) {
		default:
			fallthrough
		case 3:
			tmp, _ := strconv.ParseUint(fields[2], 10, 32)
			P = int(tmp)
			fallthrough
		case 2:
			tmp, _ := strconv.ParseUint(fields[1], 10, 8)
			C = uint8(tmp)
			fallthrough
		case 1:
			tmp, _ := hex.DecodeString(fields[0])
			if len(tmp) == 2 {
				rid = (uint16(tmp[1]) << 8) | uint16(tmp[0])
			}

		case 0:
			Vln(1, "[open]?!!")
			continue
		}

		Vf(5, "[dbg][egg][add]%04X, %d, %d\n", rid, C, P)
		eggPool2.Add(rid, C, P)
	}

	eggPool = eggPool2
	Vln(4, "[dbg][eggpool]", eggPool)
	return nil
}

func readFile(path string) ([]string, error) {
	af, err := os.Open(path)
	if err != nil {
		Vln(2, "[open]", err)
		return nil, err
	}
	defer af.Close()

	data := make([]string, 0)
	r := bufio.NewReader(af)
	b, err := r.Peek(3)
	if err != nil {
		return nil, err
	}
	if b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		r.Discard(3)
	}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.Trim(line, "\n\r\t")
		data = append(data, line)
	}

	Vln(7, "[dbg][file]", data)
	return data, nil
}

func readExtra() (error) {
	lines, err := readFile(*extraData)
	if err != nil {
		Vln(2, "[extra]", err)
		return err
	}

	tab := make(map[string][]byte)
	tmpName := ""
	tmpVal := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}

		Vln(6, "[dbg]", line)
		if strings.HasPrefix(line, "$") {
			tmpName = strings.Trim(line, "${ ")
			continue
		}
		if strings.HasPrefix(line, "}") {
			buf := Raw2Byte(tmpVal)
			tmpVal = ""
			if buf == nil {
				Vf(2, "[dbg][extra]%v decode error!!\n", tmpName)
				continue
			}
			tab[tmpName] = buf
			Vf(4, "[dbg][extra]%v = %v[% 02X]\n", tmpName, len(buf), buf)
			continue
		}

		tmpVal += line
	}

	Vln(5, "[dbg][tab]", len(tab), tab)

	// parse
	for k, v := range tab {
		switch k {
		case "UNIT1":
			if len(v) == len(WZC) {
				copy(WZC, v)
				Vf(3, "[extra]update %v[%d]\n", k, len(v))
			}
		case "UNIT2":
			if len(v) == len(IJ) {
				copy(IJ, v)
				Vf(3, "[extra]update %v[%d]\n", k, len(v))
			}
		case "UserInfo001":
			if len(v) == len(UserInfo001) {
				copy(UserInfo001, v)
				Vf(3, "[extra]update %v[%d]\n", k, len(v))
			}
		case "UserInfo002":
			if len(v) == len(UserInfo002) {
				copy(UserInfo002, v)
				Vf(3, "[extra]update %v[%d]\n", k, len(v))
			}
		case "PageHead":
			copy(PageHead, v)
			Vf(3, "[extra]update %v[%d]\n", k, len(v))
		}
	}

	return nil
}