package creditsmsparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Vendor string
const (
	SamsungCheck = Vendor("삼성체크")
	Samsung = Vendor("삼성")
	Shinhan = Vendor("신한")
)
var vendorArray = [...]Vendor{SamsungCheck, Samsung, Shinhan}

type Approval string
const (
	Approve = Approval("승인")
	Cancel = Approval("취소")
)
var approvalArray = [...]Approval{Approve, Cancel}

type PaymentInfo struct {
	vendor       Vendor
	last4digit   string
	approval     Approval
	price        int
	installments int
	time         time.Time
	shop         string
	cumulative   int
}

func (p PaymentInfo) ToString() string {
	return fmt.Sprintf("[Vendor]%s, [Number]%s, [Approval]%s, [Price]%d, [Installments]%d, [Time]%s, [Shop]%s, [Cumulative]%d",
		p.vendor, p.last4digit, p.approval, p.price, p.installments, p.time.Format(time.RFC822), p.shop, p.cumulative)
}

func Parse(sms string) PaymentInfo {
	sms = strings.ReplaceAll(sms, "\n", " ")
	return PaymentInfo{
		parseVendor(sms),
		parseLast4Digit(sms),
		parseApproval(sms),
		parsePrice(sms),
		parseInstallments(sms),
		parseTimestamp(sms),
		parseShop(sms),
		parseCumulative(sms)}
}

func parseVendor(sms string) Vendor{
	for _, v := range vendorArray {
		if strings.Contains(sms, string(v)) {
			return v
		}
	}
	panic("Not support vendor")
}

func parseLast4Digit(sms string) string {
	re := regexp.MustCompile(`(?:\d{4})`)
	return re.FindString(sms)
}

func parseApproval(sms string) Approval {
	for _, t := range approvalArray {
		if strings.Contains(sms, string(t)) {
			return t
		}
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

func parseTimestamp(sms string) time.Time {
	dateString := regexp.MustCompile(`\d{2}/\d{2}`).FindString(sms)
	timeString := regexp.MustCompile(`\d{2}:\d{2}`).FindString(sms)
	currentTime := fmt.Sprintf("%d-%sT%s:00+09:00", time.Now().Local().Year(), strings.ReplaceAll(dateString, "/", "-"), timeString)
	t, e := time.Parse(time.RFC3339, currentTime)
	if e != nil {
		panic("parse time error")
	}
	return t
}

func parseShop(sms string) string {
	re := regexp.MustCompile(`(:\d{2}[ |\n])(.+)`)
	sms = re.FindString(sms)
	size := len(sms)
	if strings.Contains(sms, "누적") {
		size = strings.Index(sms, "누적")
	}
	return strings.TrimSpace(sms[4:size])
}

func parseCumulative(sms string) int {
	re := regexp.MustCompile(`누적([\d,\-]+원)`)
	priceStr := re.FindAllString(sms, -1)
	if len(priceStr) <= 0 {
		return 0
	}
	return priceStrToInt(priceStr[0])
}
