package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	garantipay "github.com/OzqurYalcin/garantipay/src"
)

func main() {
	request := garantipay.Request{}
	request.Mode = "TEST" // TEST : "TEST" - PRODUCTION "PROD"
	request.Version = "v0.00"
	request.Terminal.ID = "111995"          // Terminal no
	request.Terminal.MerchantID = "600218"  // İşyeri No
	request.Terminal.UserID = "PROVAUT"     // Kullanıcı adı
	request.Terminal.ProvUserID = "PROVAUT" // Prov. Kullanıcı adı
	// Ödeme
	request.Order.OrderID = ""                                      // Sipariş numarası
	request.Customer.IPAddress = "127.0.0.1"                        // Müşteri IP adresi
	request.Card.Number = "4242424242424242"                        // Kart numarası
	request.Card.ExpireDate = "1110"                                // Kart son kullanma tarihi
	request.Card.CVV2 = ""                                          // Kart Cvv2 Kodu
	request.Transaction.Amount = "1"                                // Satış tutarı
	request.Transaction.CurrencyCode = garantipay.Currencies["TRY"] // Para birimi
	request.Transaction.MotoInd = "H"
	request.Transaction.Type = "sales"
	request.Order.AddressList.Address.Type = "B"
	request.Order.AddressList.Address.Name = ""      // Kart sahibi
	request.Order.AddressList.Address.Company = ""   // Fatura unvanı
	request.Order.AddressList.Address.GsmNumber = "" // Cep telefonu
	password := "123qweASD"
	hashpassword := strings.ToUpper(garantipay.SHA1(password + fmt.Sprintf("%09v", request.Terminal.ID)))
	hashdata := request.Order.OrderID + request.Terminal.ID + request.Card.Number + request.Transaction.Amount + hashpassword
	request.Terminal.HashData = strings.ToUpper(garantipay.SHA1(hashdata))
	// 3D (varsa)
	//request.Transaction.CardholderPresentCode = "0"
	//request.Transaction.Secure3D.TxnID = ""
	//request.Transaction.Secure3D.SecurityLevel = ""
	//request.Transaction.Secure3D.AuthenticationCode = ""
	response := garantipay.Transaction(request)
	pretty, _ := xml.MarshalIndent(response, " ", " ")
	fmt.Println(string(pretty))
}
