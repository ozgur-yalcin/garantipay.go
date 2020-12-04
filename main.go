package main

import (
	garantipay "github.com/OzqurYalcin/garantipay/src"
)

func main() {
	api := garantipay.API{"test"} // "test","prod"
	request := garantipay.Request{}
	request.Terminal.ID = ""                // Kullanıcı adı
	request.Terminal.MerchantID = ""        // Müşteri No
	request.Terminal.UserID = "PROVAUT"     // Kullanıcı adı
	request.Terminal.ProvUserID = "PROVAUT" // Kullanıcı adı
	// Ödeme
	request.Transaction.Type = "Auth"
	request.Mode = "TEST"                                           // TEST : "TEST" - PRODUCTION "PROD"
	request.Customer.IPAddress = ""                                 // Müşteri IP adresi
	request.Card.Number = ""                                        // Kart numarası
	request.Card.ExpireDate = "xx/xx"                               // Kart son kullanma tarihi
	request.Card.CVV2 = "xxx"                                       // Kart Cvv2 Kodu
	request.Transaction.Amount = "0.00"                             // Satış tutarı
	request.Transaction.CurrencyCode = garantipay.Currencies["TRY"] // Para birimi
	request.Transaction.MotoInd = "H"
	// 3D (varsa)
	request.Transaction.CardholderPresentCode = nil
	request.Transaction.Secure3D.TxnID = nil
	request.Transaction.Secure3D.SecurityLevel = nil
	request.Transaction.Secure3D.AuthenticationCode = nil
	api.Transaction(request)
}
