package garantipay

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var EndPoints = map[string]string{
	"TEST": "https://sanalposprovtest.garanti.com.tr/VPServlet",
	"PROD": "https://sanalposprov.garanti.com.tr/VPServlet",

	"TEST3D": "https://sanalposprovtest.garanti.com.tr/servlet/gt3dengine",
	"PROD3D": "https://sanalposprov.garanti.com.tr/servlet/gt3dengine",
}

var CurrencyCode = map[string]string{
	"TRY": "949",
	"YTL": "949",
	"TRL": "949",
	"TL":  "949",
	"USD": "840",
	"EUR": "978",
	"GBP": "826",
	"JPY": "392",
}

var CurrencyISO = map[string]string{
	"949": "TRY",
	"840": "USD",
	"978": "EUR",
	"826": "GBP",
	"392": "JPY",
}

type API struct {
	ProvUser string
	ProvPass string
	Key      string
}

type Request struct {
	XMLName     xml.Name     `xml:"GVPSRequest,omitempty"`
	Mode        string       `xml:"Mode,omitempty" form:"mode,omitempty"`
	Version     string       `xml:"Version,omitempty" form:"apiversion,omitempty"`
	Company     string       `xml:",omitempty" form:"companyname,omitempty"`
	RefreshTime string       `xml:",omitempty" form:"refreshtime,omitempty"`
	Lang        string       `xml:",omitempty" form:"lang,omitempty"`
	ChannelCode string       `xml:"ChannelCode,omitempty"`
	Terminal    *Terminal    `xml:"Terminal,omitempty"`
	Customer    *Customer    `xml:"Customer,omitempty"`
	Card        *Card        `xml:"Card,omitempty"`
	Order       *Order       `xml:"Order,omitempty"`
	Transaction *Transaction `xml:"Transaction,omitempty"`
}

type Terminal struct {
	MerchantID string `xml:"MerchantID,omitempty" form:"terminalmerchantid,omitempty"`
	ProvUserID string `xml:"ProvUserID,omitempty" form:"terminalprovuserid,omitempty"`
	UserID     string `xml:"UserID,omitempty" form:"terminaluserid,omitempty"`
	ID         string `xml:"ID,omitempty" form:"terminalid,omitempty"`
	Hash       string `xml:"HashData,omitempty" form:"secure3dhash,omitempty"`
	Level      string `xml:",omitempty" form:"secure3dsecuritylevel,omitempty"`
}

type Customer struct {
	IPAddress    string `xml:"IPAddress,omitempty" form:"customeripaddress,omitempty"`
	EmailAddress string `xml:"EmailAddress,omitempty" form:"customeremailaddress,omitempty"`
}

type Card struct {
	Number string `xml:"Number,omitempty" form:"cardnumber,omitempty"`
	Expiry string `xml:"ExpireDate,omitempty"`
	Month  string `xml:",omitempty" form:"cardexpiredatemonth,omitempty"`
	Year   string `xml:",omitempty" form:"cardexpiredateyear,omitempty"`
	Code   string `xml:"CVV2,omitempty" form:"cardcvv2,omitempty"`
}

type Order struct {
	OrderId     string       `xml:"OrderID,omitempty" form:"orderid,omitempty"`
	GroupId     string       `xml:"GroupID,omitempty"`
	AddressList *AddressList `xml:"AddressList,omitempty"`
}

type AddressList struct {
	Address *Address `xml:"Address,omitempty"`
}

type Address struct {
	Type       string `xml:"Type,omitempty"`
	Name       string `xml:"Name,omitempty"`
	LastName   string `xml:"LastName,omitempty"`
	Company    string `xml:"Company,omitempty" form:"cardholder,omitempty"`
	Text       string `xml:"Text,omitempty"`
	City       string `xml:"City,omitempty"`
	District   string `xml:"District,omitempty"`
	Country    string `xml:"Country,omitempty"`
	PostalCode string `xml:"PostalCode,omitempty"`
	Phone      string `xml:"PhoneNumber,omitempty" form:"phone,omitempty"`
	Gsm        string `xml:"GsmNumber,omitempty"`
	Fax        string `xml:"FaxNumber,omitempty"`
}

type Transaction struct {
	Type              string    `xml:"Type,omitempty" form:"txntype,omitempty"`
	SubType           string    `xml:"SubType,omitempty"`
	FirmCardNo        string    `xml:"FirmCardNo,omitempty"`
	Installment       string    `xml:"InstallmentCnt,omitempty" form:"txninstallmentcount,omitempty"`
	Amount            string    `xml:"Amount,omitempty" form:"txnamount,omitempty"`
	Currency          string    `xml:"CurrencyCode,omitempty" form:"txncurrencycode,omitempty"`
	MotoInd           string    `xml:"MotoInd,omitempty" form:"txnmotoind,omitempty"`
	PresentCode       string    `xml:"CardholderPresentCode,omitempty"`
	Description       string    `xml:"Description,omitempty"`
	OriginalRetrefNum string    `xml:"OriginalRetrefNum,omitempty"`
	Secure3D          *Secure3D `xml:"Secure3D,omitempty"`
	Timestamp         string    `xml:",omitempty" form:"txntimestamp,omitempty"`
	SuccessUrl        string    `xml:",omitempty" form:"successurl,omitempty"`
	ErrorUrl          string    `xml:",omitempty" form:"errorurl,omitempty"`
}

type Secure3D struct {
	CAVV string `xml:"AuthenticationCode,omitempty"`
	ECI  string `xml:"SecurityLevel,omitempty"`
	XID  string `xml:"TxnID,omitempty"`
	MD   string `xml:"Md,omitempty"`
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
		Number string `xml:"Number,omitempty"`
		Expiry string `xml:"ExpireDate,omitempty"`
		Code   string `xml:"CVV2,omitempty"`
	} `xml:"Card,omitempty"`

	Order struct {
		OrderId string `xml:"OrderID,omitempty"`
		GroupId string `xml:"GroupID,omitempty"`
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

func IPv4(r *http.Request) (ip string) {
	ipv4 := []string{
		r.Header.Get("X-Real-Ip"),
		r.Header.Get("X-Forwarded-For"),
		r.RemoteAddr,
	}
	for _, ipaddress := range ipv4 {
		if ipaddress != "" {
			ip = ipaddress
			break
		}
	}
	return strings.Split(ip, ":")[0]
}

func Random(n int) string {
	const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var bytes = make([]byte, n)
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func HEX(data string) (hash string) {
	b, err := hex.DecodeString(data)
	if err != nil {
		log.Println(err)
		return hash
	}
	hash = string(b)
	return hash
}

func SHA1(data string) (hash string) {
	h := sha1.New()
	h.Write([]byte(data))
	hash = hex.EncodeToString(h.Sum(nil))
	return hash
}

func B64(data string) (hash string) {
	hash = base64.StdEncoding.EncodeToString([]byte(data))
	return hash
}

func D64(data string) []byte {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println(err)
		return nil
	}
	return b
}

func Hash(data string) string {
	return B64(HEX(SHA1(data)))
}

func Api(merchant, terminal, provuser, provpass string) (*API, *Request) {
	api := new(API)
	api.ProvUser = provuser
	api.ProvPass = provpass
	request := new(Request)
	request.Terminal = new(Terminal)
	request.Customer = new(Customer)
	request.Card = new(Card)
	request.Order = new(Order)
	request.Order.AddressList = new(AddressList)
	request.Order.AddressList.Address = new(Address)
	request.Order.AddressList.Address.Type = "B"
	request.Transaction = new(Transaction)
	request.Version = "4.0"
	request.Terminal.ID = terminal
	request.Terminal.MerchantID = merchant
	request.Terminal.UserID = provuser
	request.Terminal.ProvUserID = provuser
	return api, request
}

func (api *API) SetStoreKey(key string) {
	api.Key = key
}

func (request *Request) SetMode(mode string) {
	request.Mode = mode
}

func (request *Request) SetIPAddress(ip string) {
	request.Customer.IPAddress = ip
}

func (request *Request) SetPhoneNumber(phone string) {
	request.Order.AddressList.Address.Phone = phone
}

func (request *Request) SetCardHolder(holder string) {
	request.Order.AddressList.Address.Company = holder
}

func (request *Request) SetCardNumber(number string) {
	request.Card.Number = number
}

func (request *Request) SetCardExpiry(month, year string) {
	request.Card.Expiry = month + year
	request.Card.Month = month
	request.Card.Year = year
}

func (request *Request) SetCardCode(code string) {
	request.Card.Code = code
}

func (request *Request) SetAmount(total string, currency string) {
	request.Transaction.Amount = strings.ReplaceAll(total, ".", "")
	request.Transaction.Currency = CurrencyCode[currency]
}

func (request *Request) SetInstallment(ins string) {
	request.Transaction.Installment = ins
}

func (request *Request) SetOrderId(oid string) {
	request.Order.OrderId = oid
}

func (request *Request) SetLang(lang string) {
	request.Lang = lang
}

func (api *API) PreAuth(ctx context.Context, req *Request) (Response, error) {
	req.Transaction.Type = "preauth"
	req.Transaction.MotoInd = "N"
	hashpassword := strings.ToUpper(SHA1(api.ProvPass + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := req.Order.OrderId + req.Terminal.ID + req.Card.Number + req.Transaction.Amount + hashpassword
	req.Terminal.Hash = strings.ToUpper(SHA1(hashdata))
	return api.Transaction(ctx, req)
}

func (api *API) Auth(ctx context.Context, req *Request) (Response, error) {
	req.Transaction.Type = "sales"
	req.Transaction.MotoInd = "N"
	hashpassword := strings.ToUpper(SHA1(api.ProvPass + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := req.Order.OrderId + req.Terminal.ID + req.Card.Number + req.Transaction.Amount + hashpassword
	req.Terminal.Hash = strings.ToUpper(SHA1(hashdata))
	return api.Transaction(ctx, req)
}

func (api *API) PreAuth3D(ctx context.Context, req *Request) (Response, error) {
	req.Transaction.Type = "preauth"
	req.Transaction.MotoInd = "N"
	hashpassword := strings.ToUpper(SHA1(api.ProvPass + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := req.Order.OrderId + req.Terminal.ID + req.Card.Number + req.Transaction.Amount + hashpassword
	req.Terminal.Hash = strings.ToUpper(SHA1(hashdata))
	return api.Transaction(ctx, req)
}

func (api *API) Auth3D(ctx context.Context, req *Request) (Response, error) {
	req.Transaction.Type = "sales"
	req.Transaction.MotoInd = "N"
	hashpassword := strings.ToUpper(SHA1(api.ProvPass + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := req.Order.OrderId + req.Terminal.ID + req.Card.Number + req.Transaction.Amount + hashpassword
	req.Terminal.Hash = strings.ToUpper(SHA1(hashdata))
	return api.Transaction(ctx, req)
}

func (api *API) PreAuth3Dhtml(ctx context.Context, req *Request) (string, error) {
	req.RefreshTime = "0"
	req.Terminal.Level = "3D"
	req.Transaction.Type = "preauth"
	req.Transaction.MotoInd = "N"
	req.Transaction.Timestamp = fmt.Sprintf("%v", time.Now().Unix())
	hashpassword := strings.ToUpper(SHA1(api.ProvPass + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := req.Terminal.ID + req.Order.OrderId + req.Transaction.Amount + req.Transaction.SuccessUrl + req.Transaction.ErrorUrl + req.Transaction.Type + req.Transaction.Installment + fmt.Sprintf("%x", api.Key) + hashpassword
	req.Terminal.Hash = strings.ToUpper(SHA1(hashdata))
	return api.Transaction3D(ctx, req)
}

func (api *API) Auth3Dhtml(ctx context.Context, req *Request) (string, error) {
	req.RefreshTime = "0"
	req.Terminal.Level = "3D"
	req.Transaction.Type = "sales"
	req.Transaction.MotoInd = "N"
	req.Transaction.Timestamp = fmt.Sprintf("%v", time.Now().Unix())
	hashpassword := strings.ToUpper(SHA1(api.ProvPass + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := req.Terminal.ID + req.Order.OrderId + req.Transaction.Amount + req.Transaction.SuccessUrl + req.Transaction.ErrorUrl + req.Transaction.Type + req.Transaction.Installment + fmt.Sprintf("%x", api.Key) + hashpassword
	req.Terminal.Hash = strings.ToUpper(SHA1(hashdata))
	return api.Transaction3D(ctx, req)
}

func (api *API) PostAuth(ctx context.Context, req *Request) (Response, error) {
	req.Transaction.Type = "postauth"
	req.Transaction.MotoInd = "N"
	hashpassword := strings.ToUpper(SHA1(api.ProvPass + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := req.Order.OrderId + req.Terminal.ID + req.Transaction.Amount + hashpassword
	req.Terminal.Hash = strings.ToUpper(SHA1(hashdata))
	return api.Transaction(ctx, req)
}

func (api *API) Refund(ctx context.Context, req *Request) (Response, error) {
	req.Transaction.Type = "refund"
	req.Transaction.MotoInd = "N"
	hashpassword := strings.ToUpper(SHA1(api.ProvPass + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := req.Order.OrderId + req.Terminal.ID + req.Transaction.Amount + hashpassword
	req.Terminal.Hash = strings.ToUpper(SHA1(hashdata))
	return api.Transaction(ctx, req)
}

func (api *API) Cancel(ctx context.Context, req *Request) (Response, error) {
	req.Transaction.Type = "void"
	req.Transaction.MotoInd = "N"
	hashpassword := strings.ToUpper(SHA1(api.ProvPass + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := req.Order.OrderId + req.Terminal.ID + req.Transaction.Amount + hashpassword
	req.Terminal.Hash = strings.ToUpper(SHA1(hashdata))
	return api.Transaction(ctx, req)
}

func (api *API) Transaction(ctx context.Context, req *Request) (res Response, err error) {
	payload, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[req.Mode], strings.NewReader(xml.Header+string(payload)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	if code, err := strconv.Atoi(res.Transaction.Response.Code); err == nil {
		switch code {
		case 0:
			return res, nil
		default:
			return res, errors.New(res.Transaction.Response.ErrorMsg)
		}
	} else {
		return res, errors.New(res.Transaction.Response.ErrorMsg)
	}
}

func (api *API) Transaction3D(ctx context.Context, req *Request) (res string, err error) {
	payload, err := QueryString(req)
	if err != nil {
		return res, err
	}
	html := []string{}
	html = append(html, `<!DOCTYPE html>`)
	html = append(html, `<html>`)
	html = append(html, `<head>`)
	html = append(html, `<meta http-equiv="Content-Type" content="text/html; charset=utf-8">`)
	html = append(html, `<script type="text/javascript">function submitonload() {document.payment.submit();document.getElementById('button').remove();document.getElementById('body').insertAdjacentHTML("beforeend", "Lütfen bekleyiniz...");}</script>`)
	html = append(html, `</head>`)
	html = append(html, `<body onload="javascript:submitonload();" id="body" style="text-align:center;margin:10px;font-family:Arial;font-weight:bold;">`)
	html = append(html, `<form action="`+EndPoints[req.Mode+"3D"]+`" method="post" name="payment">`)
	for k := range payload {
		html = append(html, `<input type="hidden" name="`+k+`" value="`+payload.Get(k)+`">`)
	}
	html = append(html, `<input type="submit" value="Gönder" id="button">`)
	html = append(html, `</form>`)
	html = append(html, `</body>`)
	html = append(html, `</html>`)
	res = B64(strings.Join(html, "\n"))
	return res, err
}
