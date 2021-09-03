package parser_test

import (
	"log"
	"parser"
	"testing"
)

func TestParser(t *testing.T) {
	sms := "[Web발신]\n삼성0163승인 홍*동\n39,000원 일시불\n09/02 20:16 (주)에스더포\n누적2,316,207원"
	pInfo := parser.Parse(sms)
	log.Print(pInfo.ToString())
	sms = "[Web발신]\n삼성체크5620승인 홍*동\n4,730원\n11/03 21:09\n이마트에브리데이"
	pInfo = parser.Parse(sms)
	log.Print(pInfo.ToString())
	sms = "[Web발신]\n삼성0163취소 홍*동\n-1,368,930원 일시불\n07/25 17:33 쿠팡\n누적3,732,118원"
	pInfo = parser.Parse(sms)
	log.Print(pInfo.ToString())
	sms = "[Web발신]\n신한카드(5688)승인 홍*동 1,350원(일시불)08/30 23:30 결제대행2_4 누적611,640원"
	pInfo = parser.Parse(sms)
	log.Print(pInfo.ToString())
	sms = "[Web발신]\n신한카드(5688)취소 홍*동 150,000원(일시불)08/07 14:47 망향주유소 누적994,030원"
	pInfo = parser.Parse(sms)
	log.Print(pInfo.ToString())
	// sms = "[Web발신]\n[삼성카드]0163\n자동결제 08/23접수\n한화손해보험(주)\n127,860원"
	// pInfo = parser.Parse(sms)
	// log.Print(pInfo.ToString())
}
