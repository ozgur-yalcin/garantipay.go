package garantipay

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"
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

type Request struct {
	XMLName     xml.Name    `xml:"GVPSRequest,omitempty"`
	Mode        interface{} `xml:"Mode,omitempty"`
	Version     interface{} `xml:"Version,omitempty"`
	ChannelCode interface{} `xml:"ChannelCode,omitempty"`

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
		OrderID     interface{} `xml:"OrderID,omitempty"`
		GroupID     interface{} `xml:"GroupID,omitempty"`
		AddressList struct {
			Address struct {
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
			} `xml:"Address,omitempty"`
		} `xml:"AddressList,omitempty"`
	} `xml:"Order,omitempty"`

	Transaction struct {
		Type                  interface{} `xml:"Type,omitempty"`
		SubType               interface{} `xml:"SubType,omitempty"`
		FirmCardNo            interface{} `xml:"FirmCardNo,omitempty"`
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
	XMLName xml.Name    `xml:"GVPSResponse,omitempty"`
	Mode    interface{} `xml:"Mode,omitempty"`

	Terminal struct {
		MerchantID interface{} `xml:"MerchantID,omitempty"`
		ProvUserID interface{} `xml:"ProvUserID,omitempty"`
		UserID     interface{} `xml:"UserID,omitempty"`
		ID         interface{} `xml:"ID,omitempty"`
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
		RetrefNum        interface{} `xml:"RetrefNum,omitempty"`
		AuthCode         interface{} `xml:"AuthCode,omitempty"`
		BatchNum         interface{} `xml:"BatchNum,omitempty"`
		SequenceNum      interface{} `xml:"SequenceNum,omitempty"`
		ProvDate         interface{} `xml:"ProvDate,omitempty"`
		CardNumberMasked interface{} `xml:"CardNumberMasked,omitempty"`
		CardHolderName   interface{} `xml:"CardHolderName,omitempty"`
		CardType         interface{} `xml:"CardType,omitempty"`
		HashData         interface{} `xml:"HashData,omitempty"`
		HostMsgList      interface{} `xml:"HostMsgList,omitempty"`
		Response         struct {
			Source     interface{} `xml:"Source,omitempty"`
			Code       interface{} `xml:"Code,omitempty"`
			ReasonCode interface{} `xml:"ReasonCode,omitempty"`
			Message    interface{} `xml:"Message,omitempty"`
			ErrorMsg   interface{} `xml:"ErrorMsg,omitempty"`
			SysErrMsg  interface{} `xml:"SysErrMsg,omitempty"`
		} `xml:"Response,omitempty"`
	} `xml:"Transaction,omitempty"`
}

func Transaction(request Request) (response Response) {
	postdata, _ := xml.Marshal(request)
	res, err := http.Post(EndPoints[request.Mode.(string)], "text/xml; charset=utf-8", strings.NewReader(strings.ToLower(xml.Header)+string(postdata)))
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
