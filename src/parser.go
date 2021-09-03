package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type PaymentInfo struct {
	vendor       string
	last4digit   string
	kind         string
	price        int
	installments int
	date         string
	time         string
	shop         string
	cumulative   int
}

func (p PaymentInfo) ToString() string {
	return fmt.Sprintf("[Vendor]%s, [Number]%s, [Kind]%s, [Price]%d, [Installments]%d, [Date]%s, [Time]%s, [Shop]%s, [Cumulative]%d",
		p.vendor, p.last4digit, p.kind, p.price, p.installments, p.date, p.time, p.shop, p.cumulative)
}

func Parse(sms string) PaymentInfo {
	return PaymentInfo{
		parseVendor(sms),
		parseLast4Digit(sms),
		parseKind(sms),
		parsePrice(sms),
		parseInstallments(sms),
		parseDate(sms),
		parseTime(sms),
		parseShop(sms),
		parseCumulative(sms)}
}

func parseVendor(sms string) string {
	var vendorArray = [...]string{`삼성체크`, `삼성`, `신한`}
	for _, s := range vendorArray {
		if strings.Contains(sms, s) {
			return s
		}
	}
	panic("Not support vendor")
}

func parseLast4Digit(sms string) string {
	re := regexp.MustCompile(`(?:\d{4})`)
	return re.FindString(sms)
}

func parseKind(sms string) string {
	switch {
	case strings.Contains(sms, "승인"):
		return "승인"
	case strings.Contains(sms, "취소"):
		return "취소"
	}
	panic("Not support kind")
}

func priceStrToInt(sms string) int {
	var removeStr = [...]string{`누적`, `,`, `원`}
	for _, r := range removeStr {
		sms = strings.ReplaceAll(sms, r, ``)
	}
	p, e := strconv.Atoi(sms)
	if e != nil {
		panic(e)
	}
	return p
}

func parsePrice(sms string) int {
	re := regexp.MustCompile(`([\d,\-]+원)`)
	priceStr := re.FindAllString(sms, -1)
	if len(priceStr) <= 0 {
		panic("too few prices")
	}
	return priceStrToInt(priceStr[0])
}

func parseInstallments(str string) int {
	if !strings.Contains(str, "개월") {
		return 1
	}
	re := regexp.MustCompile(`\d*개월`)
	str = re.FindString(str)
	str = strings.Replace(str, `개월`, ``, 1)
	p, e := strconv.Atoi(str)
	if e != nil {
		panic(e)
	}
	return p
}

func parseDate(sms string) string {
	re := regexp.MustCompile(`\d{2}/\d{2}`)
	return re.FindString(sms)
}

func parseTime(sms string) string {
	re := regexp.MustCompile(`\d{2}:\d{2}`)
	return re.FindString(sms)
}

func parseShop(sms string) string {
	re := regexp.MustCompile(`(:\d{2}[ |\n])(.+)`)
	sms = re.FindString(sms)
	size := len(sms)
	if strings.Contains(sms, " 누적") {
		size = strings.Index(sms, " 누적")
	}
	return sms[4:size]
}

func parseCumulative(sms string) int {
	re := regexp.MustCompile(`누적([\d,\-]+원)`)
	priceStr := re.FindAllString(sms, -1)
	if len(priceStr) <= 0 {
		return 0
	}
	return priceStrToInt(priceStr[0])
}
