package entity

import "time"

type CreateOrderRequest struct {
	Order struct {
		OrderID     string `json:"OrderId"`
		ItemID      string `json:"ItemId"`
		Quantity    int    `json:"Quantity"`
		ItemName    string `json:"ItemName"`
		OrderStatus string `json:"OrderStatus"`
		OrderTotal  struct {
			CurrencyCode string `json:"CurrencyCode"`
			Amount       string `json:"Amount"`
		} `json:"OrderTotal"`
		OrderType    string    `json:"OrderType"`
		PurchaseDate time.Time `json:"PurchaseDate"`
	} `json:"order"`
	Customer struct {
		CustomerID      string `json:"CustomerId"`
		CustomerName    string `json:"CustomerName"`
		CustomerEmail   string `json:"CustomerEmail"`
		CustomerAddress string `json:"CustomerAddress"`
		IsPrime         bool   `json:"IsPrime"`
	} `json:"customer"`
	Payment struct {
		PaymentID             string `json:"PaymentId"`
		PaymentStatus         string `json:"PaymentStatus"`
		PaymentMethod         string `json:"PaymentMethod"`
		CardNumber            string `json:"CardNumber"`
		CardVerificationValue string `json:"CardVerificationValue"`
		BillingAddress        struct {
			Name          string `json:"Name"`
			AddressLine1  string `json:"AddressLine1"`
			City          string `json:"City"`
			StateOrRegion string `json:"StateOrRegion"`
			PostalCode    string `json:"PostalCode"`
			CountryCode   string `json:"CountryCode"`
		} `json:"BillingAddress"`
		ChargeCustomerTimestamp string `json:"ChargeCustomerTimestamp"`
	} `json:"payment"`
	DeliveryDetails struct {
		DeliveryID             string    `json:"DeliveryId"`
		StartShipmentTimestamp string    `json:"StartShipmentTimestamp"`
		DeliverierInfo         string    `json:"DeliverierInfo"`
		ShipmentService        string    `json:"ShipmentService"`
		EarliestShipDate       time.Time `json:"EarliestShipDate"`
		LatestShipDate         time.Time `json:"LatestShipDate"`
		ShippingAddress        struct {
			AddressLine1  string `json:"AddressLine1"`
			City          string `json:"City"`
			StateOrRegion string `json:"StateOrRegion"`
			PostalCode    string `json:"PostalCode"`
			CountryCode   string `json:"CountryCode"`
		} `json:"ShippingAddress"`
	} `json:"deliveryDetails"`
}
