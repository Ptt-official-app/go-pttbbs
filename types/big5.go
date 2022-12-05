package types

import (
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	utf8ToBig5 = make(map[string][]byte)
	big5ToUTF8 = make(map[string][]byte)
)

func initBig5() (err error) {
	logrus.Infof("initBig5: start")
	err = initB2U()
	if err != nil {
		return err
	}

	err = initU2B()
	if err != nil {
		return err
	}

	return nil
}

func initB2U() error {
	if len(big5ToUTF8) > 0 { // already loaded
		return nil
	}

	file, err := os.Open(BIG5_TO_UTF8)
	if err != nil {
		return err
	}
	defer file.Close()

	r := io.Reader(file)
	content, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	contentStr := string(content)

	lines := strings.Split(contentStr, "\n")
	lines = lines[1:]

	for _, line := range lines {
		lineList := strings.Split(line, " ")
		if len(lineList) != 2 {
			log.Warnf("initB2U: unable to split line: line: %v", line)
			continue
		}

		big5, err := initToBig5(lineList[0][2:])
		if err != nil {
			log.Warnf("initB2U: unable initToBig5: line: %v", lineList[0])
			continue
		}
		utf8, err := initToUtf8(lineList[1][2:])
		if err != nil {
			log.Warnf("initB2U: unable initToUtf8: line: %v", lineList[1])
			continue
		}
		big5ToUTF8[string(big5)] = utf8
	}

	log.Infof("initB2U: after map: big5ToUTF8: %v", len(big5ToUTF8))

	return nil
}

func initU2B() error {
	if len(utf8ToBig5) > 0 { // already loaded
		return nil
	}

	file, err := os.Open(UTF8_TO_BIG5)
	if err != nil {
		return err
	}
	defer file.Close()

	r := io.Reader(file)
	content, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	contentStr := string(content)

	lines := strings.Split(contentStr, "\n")
	lines = lines[1:]

	for _, line := range lines {
		lineList := strings.Split(line, " ")
		if len(lineList) != 2 {
			log.Warnf("initU2B: unable to split line: line: %v", line)
			continue
		}

		big5, err := initToBig5(lineList[0][2:])
		if err != nil {
			log.Warnf("initU2B: unable to initToBig5: line: %v", lineList[0])
			continue
		}
		utf8, err := initToUtf8(lineList[1][2:])
		if err != nil {
			log.Warnf("initU2B: unable to initToUtf8: line: %v", lineList[1])
			continue
		}
		utf8ToBig5[string(utf8)] = big5
	}

	log.Infof("initU2B: after map: utf8ToBig5: %v", len(utf8ToBig5))

	return nil
}

func initToBig5(big5Code string) (theBytes []byte, err error) {
	big5Code = strings.TrimSpace(big5Code)
	theBytes = make([]byte, 2)

	_, err = hex.Decode(theBytes, []byte(big5Code))
	if err != nil {
		return nil, err
	}

	return theBytes, nil
}

func initToUtf8(ucsCode string) (theBytes []byte, err error) {
	ucsCode = strings.TrimSpace(ucsCode)
	ucsBytes := make([]byte, 2)

	_, err = hex.Decode(ucsBytes, []byte(ucsCode))
	if err != nil {
		return nil, err
	}

	ucs2 := int(ucsBytes[0])*256 + int(ucsBytes[1])

	if (ucs2 & (^0x7f)) == 0 {
		theBytes = make([]byte, 1)
		return theBytes, nil
	}

	if (ucs2 & 0xF800) == 0 {
		// (2) 00000yyy yyxxxxxx -> 110yyyyy 10xxxxxx
		theBytes = make([]byte, 2)
		theBytes[0] = byte(0xc0 | (ucs2 >> 6))
		theBytes[1] = byte(0x80 | (ucs2 & 0x3f))
		return theBytes, nil

	} else {
		// (3) zzzzyyyy yyxxxxxx -> 1110zzzz 10yyyyyy 10xxxxxx
		theBytes = make([]byte, 3)

		theBytes[0] = byte(0xE0 | (ucs2 >> 12))
		theBytes[1] = byte(0x80 | ((ucs2 >> 6) & 0x3F))
		theBytes[2] = byte(0x80 | ((ucs2) & 0x3F))
		return theBytes, nil
	}
}

func Big5ToUtf8(big5 []byte) (utf8 string) {
	estimatedUtf8Sz := len(big5) * 3 / 2
	utf8Bytes := make([]byte, 0, estimatedUtf8Sz)
	for p_big5 := big5; len(p_big5) > 0; {
		if p_big5[0] < 0x80 {
			utf8Bytes = append(utf8Bytes, p_big5[0])
			p_big5 = p_big5[1:]
		} else {
			if len(p_big5) < 2 {
				// log.Warningf("Big5ToUtf8: unable to parse big5: p_big5: %v", p_big5)
				break
			}
			eachUtf8 := big5ToUTF8[string(p_big5[:2])]
			utf8Bytes = append(utf8Bytes, eachUtf8...)
			p_big5 = p_big5[2:]
		}
	}
	utf8 = string(utf8Bytes)

	return utf8
}

func Utf8ToBig5(utf8 string) (big5 []byte) {
	utf8Bytes := []byte(utf8)
	estimatedBig5Sz := len(utf8Bytes)
	big5 = make([]byte, 0, estimatedBig5Sz)

	for p_utf8 := utf8Bytes; len(p_utf8) > 0; {
		if p_utf8[0] < 0x80 {
			big5 = append(big5, p_utf8[0])
			p_utf8 = p_utf8[1:]
		} else if len(p_utf8) >= 2 && (p_utf8[0]&0xe0) == 0xc0 {
			eachBig5, ok := utf8ToBig5[string(p_utf8[:2])]
			if !ok {
				eachBig5 = []byte{0xff, 0xfd}
			}
			big5 = append(big5, eachBig5...)
			p_utf8 = p_utf8[2:]
		} else if len(p_utf8) >= 3 && (p_utf8[0]&0xf0) == 0xe0 {
			eachBig5, ok := utf8ToBig5[string(p_utf8[:3])]
			if !ok {
				eachBig5 = []byte{0xff, 0xfd}
			}
			big5 = append(big5, eachBig5...)
			p_utf8 = p_utf8[3:]
		}
	}

	return big5
}

func TrimDBCS(theCstr Cstr) (theBytes []byte) {
	theBytes = CstrToBytes(theCstr)
	if theBytes[len(theBytes)-1] >= 0x80 {
		theBytes[len(theBytes)-1] = 0
		theBytes = theBytes[:len(theBytes)-1]
	}

	return theBytes
}
