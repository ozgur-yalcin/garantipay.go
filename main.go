package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	garantipay "github.com/ozgur-soft/garantipay/src"
)

func main() {
	api := new(garantipay.API)
	request := new(garantipay.Request)
	request.Mode = "TEST" // TEST : "TEST" - PRODUCTION "PROD"
	request.Version = "v1.0"

	request.Terminal = new(garantipay.Terminal)
	request.Terminal.ID = "111995"          // Terminal no
	request.Terminal.MerchantID = "600218"  // İşyeri No
	request.Terminal.UserID = "PROVAUT"     // Kullanıcı adı (satış için)
	request.Terminal.ProvUserID = "PROVAUT" // Prov. Kullanıcı adı (satış için)

	request.Card = new(garantipay.Card)
	request.Card.Number = "4242424242424242" // Kart numarası
	request.Card.ExpireDate = "0220"         // Son kullanma tarihi (Ay ve Yılın son 2 hanesi) AAYY
	request.Card.CVV2 = "123"                // Cvv2 Kodu (kartın arka yüzündeki 3 haneli numara)

	request.Customer = new(garantipay.Customer)
	request.Customer.IPAddress = "1.2.3.4" // Müşteri IP adresi (zorunlu)

	request.Transaction = new(garantipay.Transaction)
	request.Transaction.Amount = "100"                              // Satış tutarı (1,00 -> 100) Son 2 hane kuruş
	request.Transaction.InstallmentCnt = ""                         // Taksit sayısı
	request.Transaction.CurrencyCode = garantipay.Currencies["TRY"] // Para birimi
	request.Transaction.MotoInd = "N"
	request.Transaction.Type = "sales"

	request.Order = new(garantipay.Order)
	request.Order.OrderID = "" // Sipariş numarası (boş bırakılacak)

	request.Order.AddressList = new(garantipay.AddressList)
	request.Order.AddressList.Address = new(garantipay.Address)
	request.Order.AddressList.Address.Type = "B"
	request.Order.AddressList.Address.Name = ""        // İsim
	request.Order.AddressList.Address.LastName = ""    // Soyisim
	request.Order.AddressList.Address.PhoneNumber = "" // Telefon numarası

	password := "123qweASD" // PROVAUT kullanıcı şifresi
	hashpassword := strings.ToUpper(garantipay.SHA1(password + fmt.Sprintf("%09v", request.Terminal.ID)))
	hashdata := fmt.Sprintf("%v", request.Order.OrderID) + fmt.Sprintf("%v", request.Terminal.ID) + fmt.Sprintf("%v", request.Card.Number) + fmt.Sprintf("%v", request.Transaction.Amount) + hashpassword
	request.Terminal.HashData = strings.ToUpper(garantipay.SHA1(hashdata))
	// 3D (varsa)
	//request.Transaction.Secure3D = new(garantipay.Secure3D)
	//request.Transaction.Secure3D.Md = ""
	//request.Transaction.Secure3D.TxnID = ""
	//request.Transaction.Secure3D.SecurityLevel = ""
	//request.Transaction.Secure3D.AuthenticationCode = ""
	//request.Transaction.CardholderPresentCode = "13"
	response := api.Transaction(request)
	pretty, _ := xml.MarshalIndent(response, " ", " ")
	fmt.Println(string(pretty))
}
