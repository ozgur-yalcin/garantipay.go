[![license](https://img.shields.io/:license-mit-blue.svg)](https://github.com/ozgur-soft/garantipay.go/blob/main/LICENSE.md)
[![documentation](https://pkg.go.dev/badge/github.com/ozgur-soft/garantipay.go)](https://pkg.go.dev/github.com/ozgur-soft/garantipay.go/src)

# Garantipay.go
Garanti Bankası Virtual POS API with golang

# Installation
```bash
go get github.com/ozgur-soft/garantipay.go
```

# Sanalpos satış işlemi
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	garantipay "github.com/ozgur-soft/garantipay.go/src"
)

func main() {
	api, req := garantipay.Api("terminal no", "işyeri no")
	req.SetMode("TEST")                   // TEST : "TEST" - PRODUCTION "PROD"
	req.SetCardNumber("4242424242424242") // Kart numarası
	req.SetCardExpiry("02", "20")         // Son kullanma tarihi (Ay ve Yılın son 2 hanesi) AAYY
	req.SetCardCode("123")                // Cvv2 Kodu (kartın arka yüzündeki 3 haneli numara)
	req.SetIPAddress("1.2.3.4")           // Müşteri IP adresi (zorunlu)
	req.SetAmount("1.00")                 // Satış tutarı (zorunlu)
	req.SetInstalment("")                 // Taksit sayısı (varsa)
	req.SetCurrency("TRY")                // Para birimi
	req.SetOrderId("")                    // Sipariş numarası (varsa)

	// Kişisel bilgiler
	req.Order.AddressList = new(garantipay.AddressList)
	req.Order.AddressList.Address = new(garantipay.Address)
	req.Order.AddressList.Address.Type = "B"
	req.Order.AddressList.Address.Name = ""        // İsim
	req.Order.AddressList.Address.LastName = ""    // Soyisim
	req.Order.AddressList.Address.PhoneNumber = "" // Telefon numarası

	// Hash
	password := "123qweASD" // PROVAUT kullanıcı şifresi
	hashpassword := strings.ToUpper(garantipay.SHA1(password + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := fmt.Sprintf("%v", req.Order.OrderID) + fmt.Sprintf("%v", req.Terminal.ID) + fmt.Sprintf("%v", req.Card.Number) + fmt.Sprintf("%v", req.Transaction.Amount) + hashpassword
	req.Terminal.HashData = strings.ToUpper(garantipay.SHA1(hashdata))

	// 3D (varsa)
	//req.Transaction.Secure3D = new(garantipay.Secure3D)
	//req.Transaction.Secure3D.Md = ""
	//req.Transaction.Secure3D.TxnID = ""
	//req.Transaction.Secure3D.SecurityLevel = ""
	//req.Transaction.Secure3D.AuthenticationCode = ""
	//req.Transaction.CardholderPresentCode = "13"

	// Satış
	ctx := context.Background()
	res := api.Pay(ctx, req)
	pretty, _ := xml.MarshalIndent(res, " ", " ")
	fmt.Println(string(pretty))
}
```

# Sanalpos iade işlemi
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	garantipay "github.com/ozgur-soft/garantipay.go/src"
)

func main() {
	api, req := garantipay.Api("terminal no", "işyeri no")
	// TEST : "TEST" - PRODUCTION "PROD"
	req.SetMode("TEST")
	// Sipariş numarası (zorunlu)
	req.SetOrderId("SISTxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	// İade tutarı (zorunlu)
	req.SetAmount("1.00")
	// Para birimi (zorunlu)
	req.SetCurrency("TRY")
	// IP adresi (zorunlu)
	req.SetIPAddress("1.2.3.4")

	// Hash
	password := "123qweASD" // PROVRFN kullanıcı şifresi
	hashpassword := strings.ToUpper(garantipay.SHA1(password + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := fmt.Sprintf("%v", req.Order.OrderID) + fmt.Sprintf("%v", req.Terminal.ID) + fmt.Sprintf("%v", req.Transaction.Amount) + hashpassword
	req.Terminal.HashData = strings.ToUpper(garantipay.SHA1(hashdata))

	// İade
	ctx := context.Background()
	res := api.Refund(ctx, req)
	pretty, _ := xml.MarshalIndent(res, " ", " ")
	fmt.Println(string(pretty))
}
```

# Sanalpos iptal işlemi
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	garantipay "github.com/ozgur-soft/garantipay.go/src"
)

func main() {
	api, req := garantipay.Api("terminal no", "işyeri no")
	// TEST : "TEST" - PRODUCTION "PROD"
	req.SetMode("TEST")
	// Sipariş numarası (zorunlu)
	req.SetOrderId("SISTxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	// İptal tutarı (zorunlu)
	req.SetAmount("1.00")
	// Para birimi (zorunlu)
	req.SetCurrency("TRY")
	// IP adresi (zorunlu)
	req.SetIPAddress("1.2.3.4")

	// Hash
	password := "123qweASD" // PROVRFN kullanıcı şifresi
	hashpassword := strings.ToUpper(garantipay.SHA1(password + fmt.Sprintf("%09v", req.Terminal.ID)))
	hashdata := fmt.Sprintf("%v", req.Order.OrderID) + fmt.Sprintf("%v", req.Terminal.ID) + fmt.Sprintf("%v", req.Transaction.Amount) + hashpassword
	req.Terminal.HashData = strings.ToUpper(garantipay.SHA1(hashdata))

	// İade
	ctx := context.Background()
	res := api.Cancel(ctx, req)
	pretty, _ := xml.MarshalIndent(res, " ", " ")
	fmt.Println(string(pretty))
}
```
