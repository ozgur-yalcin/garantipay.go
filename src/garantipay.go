package garantipay

import (
	"bytes"
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

type Request struct {
	XMLName     xml.Name `xml:"GVPSRequest,omitempty"`
	Mode        string   `xml:"Mode,omitempty"`
	Version     string   `xml:"Version,omitempty"`
	ChannelCode string   `xml:"ChannelCode,omitempty"`

	Terminal struct {
		MerchantID string `xml:"MerchantID,omitempty"`
		ProvUserID string `xml:"ProvUserID,omitempty"`
		UserID     string `xml:"UserID,omitempty"`
		ID         string `xml:"ID,omitempty"`
		HashData   string `xml:"HashData,omitempty"`
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
		OrderID     string `xml:"OrderID,omitempty"`
		GroupID     string `xml:"GroupID,omitempty"`
		AddressList struct {
			Address struct {
				Type        string `xml:"Type,omitempty"`
				Name        string `xml:"Name,omitempty"`
				LastName    string `xml:"LastName,omitempty"`
				Company     string `xml:"Company,omitempty"`
				Text        string `xml:"Text,omitempty"`
				City        string `xml:"City,omitempty"`
				District    string `xml:"District,omitempty"`
				Country     string `xml:"Country,omitempty"`
				PostalCode  string `xml:"PostalCode,omitempty"`
				PhoneNumber string `xml:"PhoneNumber,omitempty"`
				GsmNumber   string `xml:"GsmNumber,omitempty"`
				FaxNumber   string `xml:"FaxNumber,omitempty"`
			} `xml:"Address,omitempty"`
		} `xml:"AddressList,omitempty"`
	} `xml:"Order,omitempty"`

	Transaction struct {
		Type                  string `xml:"Type,omitempty"`
		SubType               string `xml:"SubType,omitempty"`
		FirmCardNo            string `xml:"FirmCardNo,omitempty"`
		InstallmentCnt        string `xml:"InstallmentCnt,omitempty"`
		Amount                string `xml:"Amount,omitempty"`
		CurrencyCode          string `xml:"CurrencyCode,omitempty"`
		CardholderPresentCode string `xml:"CardholderPresentCode,omitempty"`
		MotoInd               string `xml:"MotoInd,omitempty"`
		Description           string `xml:"Description,omitempty"`
		Secure3D              struct {
			AuthenticationCode string `xml:"AuthenticationCode,omitempty"`
			SecurityLevel      string `xml:"SecurityLevel,omitempty"`
			TxnID              string `xml:"TxnID,omitempty"`
			Md                 string `xml:"Md,omitempty"`
		} `xml:"Secure3D,omitempty"`
	} `xml:"Transaction,omitempty"`
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

func Transaction(request Request) (response Response) {
	postdata, _ := xml.Marshal(request)
	res, err := http.Post(EndPoints[request.Mode], "text/xml; charset=utf-8", bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := xml.NewDecoder(res.Body)
	decoder.Decode(&response)
	return response
}
