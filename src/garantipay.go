package garantipay

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"log"
	"net/http"
)

var EndPoints map[string]string = map[string]string{
	"TEST": "https://sanalposprovtest.garanti.com.tr/VPServlet",
	"PROD": "https://sanalposprov.garanti.com.tr/VPServlet",
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

type API struct{}

type Request struct {
	XMLName     xml.Name     `xml:"GVPSRequest,omitempty"`
	Mode        interface{}  `xml:"Mode,omitempty"`
	Version     interface{}  `xml:"Version,omitempty"`
	ChannelCode interface{}  `xml:"ChannelCode,omitempty"`
	Terminal    *Terminal    `xml:"Terminal,omitempty"`
	Customer    *Customer    `xml:"Customer,omitempty"`
	Card        *Card        `xml:"Card,omitempty"`
	Order       *Order       `xml:"Order,omitempty"`
	Transaction *Transaction `xml:"Transaction,omitempty"`
}

type Terminal struct {
	MerchantID interface{} `xml:"MerchantID,omitempty"`
	ProvUserID interface{} `xml:"ProvUserID,omitempty"`
	UserID     interface{} `xml:"UserID,omitempty"`
	ID         interface{} `xml:"ID,omitempty"`
	HashData   interface{} `xml:"HashData,omitempty"`
}

type Customer struct {
	IPAddress    interface{} `xml:"IPAddress,omitempty"`
	EmailAddress interface{} `xml:"EmailAddress,omitempty"`
}

type Card struct {
	Number     interface{} `xml:"Number,omitempty"`
	ExpireDate interface{} `xml:"ExpireDate,omitempty"`
	CVV2       interface{} `xml:"CVV2,omitempty"`
}

type Order struct {
	OrderID     interface{}  `xml:"OrderID,omitempty"`
	GroupID     interface{}  `xml:"GroupID,omitempty"`
	AddressList *AddressList `xml:"AddressList,omitempty"`
}

type AddressList struct {
	Address *Address `xml:"Address,omitempty"`
}

type Address struct {
	Type        interface{} `xml:"Type,omitempty"`
	Name        interface{} `xml:"Name,omitempty"`
	LastName    interface{} `xml:"LastName,omitempty"`
	Company     interface{} `xml:"Company,omitempty"`
	Text        interface{} `xml:"Text,omitempty"`
	City        interface{} `xml:"City,omitempty"`
	District    interface{} `xml:"District,omitempty"`
	Country     interface{} `xml:"Country,omitempty"`
	PostalCode  interface{} `xml:"PostalCode,omitempty"`
	PhoneNumber interface{} `xml:"PhoneNumber,omitempty"`
	GsmNumber   interface{} `xml:"GsmNumber,omitempty"`
	FaxNumber   interface{} `xml:"FaxNumber,omitempty"`
}

type Transaction struct {
	Type                  interface{} `xml:"Type,omitempty"`
	SubType               interface{} `xml:"SubType,omitempty"`
	FirmCardNo            interface{} `xml:"FirmCardNo,omitempty"`
	InstallmentCnt        interface{} `xml:"InstallmentCnt,omitempty"`
	Amount                interface{} `xml:"Amount,omitempty"`
	CurrencyCode          interface{} `xml:"CurrencyCode,omitempty"`
	CardholderPresentCode interface{} `xml:"CardholderPresentCode,omitempty"`
	MotoInd               interface{} `xml:"MotoInd,omitempty"`
	Description           interface{} `xml:"Description,omitempty"`
	Secure3D              *Secure3D   `xml:"Secure3D,omitempty"`
}

type Secure3D struct {
	AuthenticationCode interface{} `xml:"AuthenticationCode,omitempty"`
	SecurityLevel      interface{} `xml:"SecurityLevel,omitempty"`
	TxnID              interface{} `xml:"TxnID,omitempty"`
	Md                 interface{} `xml:"Md,omitempty"`
}

type Response struct {
	XMLName xml.Name `xml:"GVPSResponse,omitempty"`
	Mode    string   `xml:"Mode,omitempty"`

	Terminal struct {
		MerchantID string `xml:"MerchantID,omitempty"`
		ProvUserID string `xml:"ProvUserID,omitempty"`
		UserID     string `xml:"UserID,omitempty"`
		ID         string `xml:"ID,omitempty"`
	} `xml:"Terminal,omitempty"`

	Customer struct {
		IPAddress    string `xml:"IPAddress,omitempty"`
		EmailAddress string `xml:"EmailAddress,omitempty"`
	} `xml:"Customer,omitempty"`

	Card struct {
		Number     string `xml:"Number,omitempty"`
		ExpireDate string `xml:"ExpireDate,omitempty"`
		CVV2       string `xml:"CVV2,omitempty"`
	} `xml:"Card,omitempty"`

	Order struct {
		OrderID string `xml:"OrderID,omitempty"`
		GroupID string `xml:"GroupID,omitempty"`
	} `xml:"Order,omitempty"`

	Transaction struct {
		Response struct {
			Source     string `xml:"Source,omitempty"`
			Code       string `xml:"Code,omitempty"`
			ReasonCode string `xml:"ReasonCode,omitempty"`
			Message    string `xml:"Message,omitempty"`
			ErrorMsg   string `xml:"ErrorMsg,omitempty"`
			SysErrMsg  string `xml:"SysErrMsg,omitempty"`
		} `xml:"Response,omitempty"`
		RetrefNum        string `xml:"RetrefNum,omitempty"`
		AuthCode         string `xml:"AuthCode,omitempty"`
		BatchNum         string `xml:"BatchNum,omitempty"`
		SequenceNum      string `xml:"SequenceNum,omitempty"`
		ProvDate         string `xml:"ProvDate,omitempty"`
		CardNumberMasked string `xml:"CardNumberMasked,omitempty"`
		CardHolderName   string `xml:"CardHolderName,omitempty"`
		CardType         string `xml:"CardType,omitempty"`
		HashData         string `xml:"HashData,omitempty"`
		HostMsgList      string `xml:"HostMsgList,omitempty"`
		RewardInqResult  struct {
			RewardList string `xml:"RewardList,omitempty"`
			ChequeList string `xml:"ChequeList,omitempty"`
		} `xml:"RewardInqResult,omitempty"`
		GarantiCardInd string `xml:"GarantiCardInd,omitempty"`
	} `xml:"Transaction,omitempty"`
}

func SHA1(data string) (hash string) {
	h := sha1.New()
	h.Write([]byte(data))
	hash = hex.EncodeToString(h.Sum(nil))
	return hash
}

func (api *API) Transaction(ctx context.Context, req *Request) (res Response) {
	postdata, err := xml.Marshal(req)
	if err != nil {
		log.Println(err)
		return res
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[req.Mode.(string)], bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return res
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	decoder.Decode(&res)
	return res
}
