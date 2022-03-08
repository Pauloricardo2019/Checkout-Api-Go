package model

type Config struct {
	MongoDBConnString string
	GormDBConnString  string
	NatsConnString    string
	RestPort          int
	SentryDNS         string
	Environment       string
	GetnetURL         string
	DeviceURL         string
	KeyDevice         string
	ClientID          string
	ClientSecret      string
	SellerID          string
	CheckoutURL       string
	RedirectUrl       string
}
