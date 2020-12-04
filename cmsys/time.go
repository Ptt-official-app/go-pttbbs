package cmsys

// 摩羯 水瓶 雙魚 牡羊 金牛 雙子 巨蟹 獅子 處女 天秤 天蠍 射手
// Reference: http://zh.wikipedia.org/wiki/%E6%98%9F%E5%BA%A7
// https://zh.wikipedia.org/wiki/%E9%BB%83%E9%81%93%E5%8D%81%E4%BA%8C%E5%AE%AE
// XXX looks different from wikipedia, but the value was defined 25 years ago.
var horoscopeFirstDay = [12]int{
	/* Jan. */ 20, 19, 21, 20, 21, 22, 23, 23, 24, 24, 22, 22,
}

func IsLeapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}

//GetHoroscope
//
//給日期求星座
//
//Return:
//	1..12
func GetHoroscope(m int, d int) int {
	if m > 12 || m < 1 {
		return 1
	}

	if d >= horoscopeFirstDay[m-1] {
		if m == 12 {
			return 1
		} else {
			return m + 1
		}
	} else {
		return m
	}
}
