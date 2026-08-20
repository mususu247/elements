package elements

import (
	"log"
	"slices"
	"strconv"
	"strings"
	"time"
)

func any2str(value any, tag bool) string {
	switch x := value.(type) {
	case string:
		xx := x
		if tag {
			xx = str2tag(xx)
		}
		xx = strings.Trim(xx, " ")
		return xx
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	case time.Time:
		return x.Format(time.RFC3339)
	default:
		return ""
	}
}

func str2tag(tag string) string {
	tag = strings.ReplaceAll(tag, "\n", "")
	tag = strings.ReplaceAll(tag, "\t", "")
	tag = strings.ReplaceAll(tag, "\r", "")
	tag = strings.TrimSpace(tag)

	return tag
}

func checkTag(tag string) (string, bool) {
	if len(tag) == 0 {
		return "", false
	}

	tag = str2tag(tag)

	for _, v := range errHead {
		if tag[:1] == v {
			log.Printf("(Error) checkTag: %v \n", tag)
			return v, false
		}
	}

	for _, v := range errName {
		if tag == v {
			log.Printf("(Error) checkTag: %v \n", tag)
			return v, false
		}
	}

	return tag, true
}

func checkAttr(key string) (string, bool) {
	key, ok := checkTag(key)
	if !ok {
		return key, false
	}

	lowKey := strings.ToLower(key)
	if len(lowKey) >= 3 {
		if lowKey[:3] == "xml" {
			return key, false
		}
	}

	return key, true
}

func removeDupStr(list []string) []string {
	var backup []string
	backup = append(backup, list...)
	slices.Sort(backup)
	list = nil

	v := ""
	for i := range backup {
		if v < backup[i] {
			v = backup[i]
			list = append(list, v)
		}
	}
	return list
}

func anyComp(src any, dst any) bool {
	var err error

	if dst != nil {
		switch xd := dst.(type) {
		case string:
			if xd == "*" {
				return true
			}
		}
	}

	//return true:equal, false:not equal
	switch xs := src.(type) {
	case string:
		var xxd string

		switch xd := dst.(type) {
		case string:
			xxd = xd
		case int:
			xxd = strconv.Itoa(xd)
		case int64:
			xxd = strconv.FormatInt(xd, 10)
		case float64:
			xxd = strconv.FormatFloat(xd, 'f', -1, 64)
		case bool:
			xxd = strconv.FormatBool(xd)
		case time.Time:
			xxd = xd.Format(time.RFC3339)
		default:
			log.Printf("(Error) not support dst:%v dst:%v \n", src, dst)
			return false
		}

		if xs == xxd {
			return true
		}
	case int:
		var xxd int

		switch xd := dst.(type) {
		case string:
			xxd, err = strconv.Atoi(xd)
		case int:
			xxd = xd
		case int64:
			xxd = int(xd)
		case float64:
			xxd = int(xd)
		case bool:
			xxd = int(0)
		case time.Time:
			xxd = int(xd.Unix())
		default:
			log.Printf("(Error) not support dst:%v dst:%v \n", src, dst)
			return false
		}

		if err != nil {
			return true
		}
		if xs == xxd {
			return true
		}
	case int64:
		var xxd int64

		switch xd := dst.(type) {
		case string:
			xxd, err = strconv.ParseInt(xd, 10, 64)
		case int:
			xxd = int64(xd)
		case int64:
			xxd = xd
		case float64:
			xxd = int64(xd)
		case bool:
			xxd = int64(0)
		case time.Time:
			xxd = xd.Unix()
		default:
			log.Printf("(Error) not support dst:%v dst:%v \n", src, dst)
			return false
		}

		if err != nil {
			return true
		}
		if xs == xxd {
			return true
		}
	case float64:
		var xxd float64

		switch xd := dst.(type) {
		case string:
			xxd, err = strconv.ParseFloat(xd, 64)
		case int:
			xxd = float64(xd)
		case int64:
			xxd = float64(xd)
		case float64:
			xxd = xd
		case bool:
			if xd {
				xxd = 1
			}
		case time.Time:
			xxd = float64(xd.Unix())
		default:
			log.Printf("(Error) not support dst:%v dst:%v \n", src, dst)
			return false
		}

		if err != nil {
			return true
		}
		if xs == xxd {
			return true
		}
	case bool:
		var xxd bool

		switch xd := dst.(type) {
		case string:
			xxd, err = strconv.ParseBool(xd)
		case int:
			if xd != 0 {
				xxd = true
			}
		case int64:
			if xd != 0 {
				xxd = true
			}
		case float64:
			if xd != 0 {
				xxd = true
			}
		case bool:
			xxd = xd
		case time.Time:
			return false
		default:
			log.Printf("(Error) not support dst:%v dst:%v \n", src, dst)
			return false
		}

		if err != nil {
			return true
		}
		if xs == xxd {
			return true
		}
	case time.Time:
		var xxd time.Time

		switch xd := dst.(type) {
		case string:
			xxd, err = time.Parse(time.RFC3339, xd)
		case int:
			d := int64(xd)
			xxd = time.Unix(d, 0)
		case int64:
			xxd = time.Unix(xd, 0)
		case float64:
			d := int64(xd)
			xxd = time.Unix(d, 0)
		case bool:
			return false
		case time.Time:
			xxd = xd
		default:
			log.Printf("(Error) not support dst:%v dst:%v \n", src, dst)
			return false
		}

		if err != nil {
			return true
		}
		if timeComp(xs, xxd) == 0 {
			return true
		}
	default:
		log.Printf("(Error) not support src:%v dst:%v \n", src, dst)
		return false
	}
	return false
}

func A1toR1C1(address string) (int64, int64) {
	addr := strings.ToUpper(address)
	adds := []byte(addr)

	var row int64
	var col int64
	for i := range adds {
		x := int64(adds[i])

		switch {
		case (x >= 65) && (x <= 90):
			xx := x - 64
			col = (col * 26) + xx
		case (x >= 48) && (x <= 57):
			xx := x - 48
			row = (row * 10) + xx
		}
	}
	return col, row
}

func R1C1toA1(row int64, col int64) string {
	var addr string

	r := row
	for r > 26 {
		w := r % 26
		r = (r - w) / 26
		char := byte(w + 64)
		s := string(char)
		addr = s + addr
	}
	char := byte(r + 64)
	s := string(char)
	addr = s + addr

	addr = addr + strconv.FormatInt(col, 10)
	return addr
}
