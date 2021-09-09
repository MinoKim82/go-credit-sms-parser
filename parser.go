package creditsmsparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Approval string
const (
	Approve = Approval("승인")
	Cancel = Approval("취소")
)
var approvalArray = [...]Approval{Approve, Cancel}

type PaymentInfo struct {
	id			 string
	approval     Approval
	price        int
	installments int
	time         time.Time
	shop         string
	cumulative   int
}

func (p PaymentInfo) ToString() string {
	return fmt.Sprintf("[Id]%s, [Approval]%s, [Price]%d, [Installments]%d, [Time]%s, [Shop]%s, [Cumulative]%d",
		p.id, p.approval, p.price, p.installments, p.time.Format(time.RFC822), p.shop, p.cumulative)
}

func Parse(sms string) PaymentInfo {
	sms = strings.ReplaceAll(sms, "\n", " ")
	return PaymentInfo{
		parseIdentifier(sms),
		parseApproval(sms),
		parsePrice(sms),
		parseInstallments(sms),
		parseTimestamp(sms),
		parseShop(sms),
		parseCumulative(sms)}
}

func parseIdentifier(sms string) string {
	return regexp.MustCompile(`\[Web발신].(.*?)\s*(승인|취소)`).FindStringSubmatch(sms)[1]
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
	str = regexp.MustCompile(`\d*개월`).FindString(str)
	str = strings.Replace(str, `개월`, ``, 1)
	p, e := strconv.Atoi(str)
	if e != nil {
		panic("Can not parse installments" + e.Error())
	}
	return p
}

func parseTimestamp(sms string) time.Time {
	dateString := regexp.MustCompile(`\d{2}/\d{2}`).FindString(sms)
	if dateString == "" {
		panic("Can not parse date : " + sms)
	}
	timeString := regexp.MustCompile(`\d{2}:\d{2}`).FindString(sms)
	if timeString == "" {
		panic("Can not parse time : " + sms)	
	}		
	currentTime := fmt.Sprintf("%d-%sT%s:00+09:00", time.Now().Local().Year(), strings.ReplaceAll(dateString, "/", "-"), timeString)
	t, e := time.Parse(time.RFC3339, currentTime)
	if e != nil {
		panic("parse time error")
	}
	return t
}

func parseShop(sms string) string {
	shopstr := regexp.MustCompile(`(:\d{2}[ |\n])(.+)`).FindString(sms)
	size := len(shopstr)
	if strings.Contains(shopstr, "누적") {
		size = strings.Index(shopstr, "누적")
	}
	shop:= strings.TrimSpace(shopstr[4:size])
	if shop == "" {
		panic("Can not parse shop : " + sms)
	}
	return shop
}

func parseCumulative(sms string) int {
	re := regexp.MustCompile(`누적([\d,\-]+원)`)
	priceStr := re.FindAllString(sms, -1)
	if len(priceStr) <= 0 {
		return 0
	}
	return priceStrToInt(priceStr[0])
}
