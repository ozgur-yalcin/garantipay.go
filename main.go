package main

import (
	garantipay "github.com/OzqurYalcin/garantipay/src"
)

func main() {
	request := garantipay.Request{}
	request.Terminal.ID = "111995"          // Terminal no
	request.Terminal.MerchantID = "600218"  // İşyeri No
	request.Terminal.UserID = "PROVAUT"     // Kullanıcı adı
	request.Terminal.ProvUserID = "PROVAUT" // Prov. Kullanıcı adı
	// Ödeme
	request.Mode = "TEST"                                           // TEST : "TEST" - PRODUCTION "PROD"
	request.Customer.IPAddress = ""                                 // Müşteri IP adresi
	request.Card.Number = ""                                        // Kart numarası
	request.Card.ExpireDate = "xx/xx"                               // Kart son kullanma tarihi
	request.Card.CVV2 = "xxx"                                       // Kart Cvv2 Kodu
	request.Transaction.Amount = "0.00"                             // Satış tutarı
	request.Transaction.CurrencyCode = garantipay.Currencies["TRY"] // Para birimi
	request.Transaction.MotoInd = "N"
	request.Transaction.Type = "Auth"
	request.Order.AddressList.Address.Type = "B"
	request.Order.AddressList.Address.Name = ""      // Kart sahibi
	request.Order.AddressList.Address.Company = ""   // Fatura unvanı
	request.Order.AddressList.Address.GsmNumber = "" // Cep telefonu
	// 3D (varsa)
	request.Transaction.CardholderPresentCode = nil
	request.Transaction.Secure3D.TxnID = nil
	request.Transaction.Secure3D.SecurityLevel = nil
	request.Transaction.Secure3D.AuthenticationCode = nil
	garantipay.Transaction(request)
}
