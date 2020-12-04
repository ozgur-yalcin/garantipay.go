package garantipay

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"
)

var EndPoints map[string]string = map[string]string{
	"test": "https://sanalposprovtest.garanti.com.tr/VPServlet",
	"prod": "https://sanalposprov.garanti.com.tr/VPServlet",
}

var Currencies map[string]string = map[string]string{
	"TRY": "949",
	"YTL": "949",
	"TRL": "949",
	"TL":  "949",
	"USD": "840",
	"EUR": "978",
	"GBP": "826",
	"JPY": "392",
}

type API struct {
	Bank string
}

type Request struct {
	XMLName xml.Name    `xml:"GVPSRequest,omitempty"`
	Mode    interface{} `xml:"Mode,omitempty"`
	Version interface{} `xml:"Version,omitempty"`

	Terminal struct {
		MerchantID interface{} `xml:"MerchantID,omitempty"`
		ProvUserID interface{} `xml:"ProvUserID,omitempty"`
		UserID     interface{} `xml:"UserID,omitempty"`
		ID         interface{} `xml:"ID,omitempty"`
		HashData   interface{} `xml:"HashData,omitempty"`
	} `xml:"Terminal,omitempty"`

	Customer struct {
		IPAddress    interface{} `xml:"IPAddress,omitempty"`
		EmailAddress interface{} `xml:"EmailAddress,omitempty"`
	} `xml:"Customer,omitempty"`

	Card struct {
		Number     interface{} `xml:"Number,omitempty"`
		ExpireDate interface{} `xml:"ExpireDate,omitempty"`
		CVV2       interface{} `xml:"CVV2,omitempty"`
	} `xml:"Card,omitempty"`

	Order struct {
		OrderID interface{} `xml:"OrderID,omitempty"`
		GroupID interface{} `xml:"GroupID,omitempty"`
	} `xml:"Order,omitempty"`

	Transaction struct {
		Type                  interface{} `xml:"Type,omitempty"`
		InstallmentCnt        interface{} `xml:"InstallmentCnt,omitempty"`
		Amount                interface{} `xml:"Amount,omitempty"`
		CurrencyCode          interface{} `xml:"CurrencyCode,omitempty"`
		CardholderPresentCode interface{} `xml:"CardholderPresentCode,omitempty"`
		MotoInd               interface{} `xml:"MotoInd,omitempty"`
		Description           interface{} `xml:"Description,omitempty"`
		Secure3D              struct {
			AuthenticationCode interface{} `xml:"AuthenticationCode,omitempty"`
			SecurityLevel      interface{} `xml:"SecurityLevel,omitempty"`
			TxnID              interface{} `xml:"TxnID,omitempty"`
			Md                 interface{} `xml:"Md,omitempty"`
		} `xml:"Secure3D,omitempty"`
	} `xml:"Transaction,omitempty"`
}

type Response struct {
	XMLName        xml.Name `xml:"GVPSResponse,omitempty"`
	AuthCode       string   `xml:"AuthCode,omitempty"`
	HostRefNum     string   `xml:"HostRefNum,omitempty"`
	ProcReturnCode string   `xml:"ProcReturnCode,omitempty"`
	ErrMsg         string   `xml:"ErrMsg,omitempty"`
}

func (api API) Transaction(request Request) (response Response) {
	response = Response{}
	postdata, _ := xml.Marshal(request)
	res, err := http.Post(EndPoints[api.Bank], "text/xml; charset=utf-8", strings.NewReader(strings.ToLower(xml.Header)+string(postdata)))
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := xml.NewDecoder(res.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	decoder.Decode(&response)
	return response
}
