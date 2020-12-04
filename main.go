package main

import (
	"encoding/json"
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
	request.Customer.IPAddress = "176.88.97.59"                     // Müşteri IP adresi
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
	hashdata := fmt.Sprintf("%v", request.Order.OrderID) + fmt.Sprintf("%v", request.Terminal.ID) + fmt.Sprintf("%v", request.Card.Number) + fmt.Sprintf("%v", request.Transaction.Amount) + hashpassword
	fmt.Println(hashdata)
	request.Terminal.HashData = strings.ToUpper(garantipay.SHA1(hashdata))
	// 3D (varsa)
	request.Transaction.CardholderPresentCode = 0
	request.Transaction.Secure3D.TxnID = nil
	request.Transaction.Secure3D.SecurityLevel = nil
	request.Transaction.Secure3D.AuthenticationCode = nil
	response := garantipay.Transaction(request)
	pretty, _ := json.MarshalIndent(response, " ", " ")
	fmt.Println(string(pretty))
}
