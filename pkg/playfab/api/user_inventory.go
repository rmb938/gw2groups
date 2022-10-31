package api

import "time"

type VirtualCurrencyRechargeTimeResponse struct {
	RechargeMax       int       `json:"RechargeMax"`
	RechargeTime      time.Time `json:"RechargeTime"`
	SecondsToRecharge int       `json:"SecondsToRecharge"`
}

type ItemInstanceResponse struct {
	Annotation        string            `json:"Annotation"`
	BundleContents    []string          `json:"BundleContents"`
	BundleParent      string            `json:"BundleParent"`
	CatalogVersion    string            `json:"CatalogVersion"`
	CustomData        map[string]string `json:"CustomData"`
	DisplayName       string            `json:"DisplayName"`
	Expiration        time.Time         `json:"Expiration"`
	ItemClass         string            `json:"ItemClass"`
	ItemId            string            `json:"ItemId"`
	ItemInstanceId    string            `json:"ItemInstanceId"`
	PurchaseDate      time.Time         `json:"PurchaseDate"`
	RemainingUses     int               `json:"RemainingUses"`
	UnitCurrency      string            `json:"UnitCurrency"`
	UnitPrice         int               `json:"UnitPrice"`
	UsesIncrementedBy int               `json:"UsesIncrementedBy"`
}
