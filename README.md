[![license](https://img.shields.io/:license-mit-blue.svg)](https://github.com/ozgur-yalcin/garantipay.go/blob/main/LICENSE.md)
[![documentation](https://pkg.go.dev/badge/github.com/ozgur-yalcin/garantipay.go)](https://pkg.go.dev/github.com/ozgur-yalcin/garantipay.go/src)

# Garantipay.go
Garanti Bankası POS API with golang

# Installation
```bash
go get github.com/ozgur-yalcin/garantipay.go
```

# Satış
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	garantipay "github.com/ozgur-yalcin/garantipay.go/src"
)

// Pos bilgileri
const (
	envmode  = "PROD"    // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	merchant = ""        // İşyeri numarası
	terminal = ""        // Terminal numarası
	username = "PROVAUT" // PROVAUT
	password = ""        // PROVAUT şifresi
)

func main() {
	api, req := garantipay.Api(merchant, terminal, username, password)
	req.SetMode(envmode)

	req.SetIPAddress("1.2.3.4")           // Müşteri IPv4 adresi (zorunlu)
	req.SetCardHolder("AD SOYAD")         // Kart sahibi (zorunlu)
	req.SetCardNumber("4242424242424242") // Kart numarası (zorunlu)
	req.SetCardExpiry("02", "20")         // Son kullanma tarihi - AA,YY (zorunlu)
	req.SetCardCode("123")                // Kart arkasındaki 3 haneli numara (zorunlu)
	req.SetPhoneNumber("05554443322")     // Müşteri telefon numarası (zorunlu)
	req.SetAmount("1.00", "TRY")          // Satış tutarı ve para birimi (zorunlu)
	req.SetInstallment("0")               // Taksit sayısı (varsa)
	req.SetOrderId("")                    // Sipariş numarası (varsa)

	// Satış
	ctx := context.Background()
	if res, err := api.Auth(ctx, req); err == nil {
		pretty, _ := xml.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```

# İade
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	garantipay "github.com/ozgur-yalcin/garantipay.go/src"
)

// Pos bilgileri
const (
	envmode  = "PROD"    // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	merchant = ""        // İşyeri numarası
	terminal = ""        // Terminal numarası
	username = "PROVRFN" // PROVRFN
	password = ""        // PROVRFN şifresi
)

func main() {
	api, req := garantipay.Api(merchant, terminal, username, password)
	req.SetMode(envmode)

	req.SetIPAddress("1.2.3.4")            // Müşteri IPv4 adresi (zorunlu)
	req.SetAmount("1.00", "TRY")           // İade tutarı ve para birimi (zorunlu)
	req.SetOrderId("SISTxxxxxxxxxxxxxxxx") // Sipariş numarası (zorunlu)

	// İade
	ctx := context.Background()
	if res, err := api.Refund(ctx, req); err == nil {
		pretty, _ := xml.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```

# İptal
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	garantipay "github.com/ozgur-yalcin/garantipay.go/src"
)

// Pos bilgileri
const (
	envmode  = "PROD"    // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	merchant = ""        // İşyeri numarası
	terminal = ""        // Terminal numarası
	username = "PROVRFN" // PROVRFN
	password = ""        // PROVRFN şifresi
)

func main() {
	api, req := garantipay.Api(merchant, terminal, username, password)
	req.SetMode(envmode)

	req.SetIPAddress("1.2.3.4")            // Müşteri IPv4 adresi (zorunlu)
	req.SetAmount("1.00", "TRY")           // İptal tutarı ve para birimi (zorunlu)
	req.SetOrderId("SISTxxxxxxxxxxxxxxxx") // Sipariş numarası (zorunlu)

	// İptal
	ctx := context.Background()
	if res, err := api.Cancel(ctx, req); err == nil {
		pretty, _ := xml.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```
