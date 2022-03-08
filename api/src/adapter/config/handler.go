package config

import (
	"os"
	"ravxcheckout/src/internal/model"
	"strconv"
)

var config *model.Config

func init() {
	// err := godotenv.Load(filepath.Join("~/Documents/works/checkout-santuu/", ".env"))
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
}

type GetConfigFn func() *model.Config

func GetConfig() *model.Config {
	if config == nil {
		config = &model.Config{
			MongoDBConnString: os.Getenv("MONGODB_CONNSTRING"),
			GormDBConnString:  os.Getenv("GORM_CONNSTRING"),
			NatsConnString:    os.Getenv("NATS_CONNSTRING"),
			SentryDNS:         os.Getenv("SENTRY_DNS"),
			Environment:       os.Getenv("ENV"),
			GetnetURL:         os.Getenv("GETNET_URL"),
			DeviceURL:         os.Getenv("DEVICE_URL"),
			KeyDevice:         os.Getenv("KEY_DEVICE"),
			ClientID:          os.Getenv("CLIENT_ID"),
			ClientSecret:      os.Getenv("CLIENT_SECRET"),
			SellerID:          os.Getenv("SELLER_ID"),
			CheckoutURL:       os.Getenv("CHECKOUT_URL"),
			RedirectUrl:       os.Getenv("REDIRECT_URL"),
		}

		if config.Environment == "" {
			config.Environment = "PRODUCTION"
		}

		restPortString := os.Getenv("REST_PORT")

		if restPortString == "" {
			restPortString = "9090"
		}

		restPort, err := strconv.Atoi(restPortString)
		if err != nil {
			panic(err.Error())
		}

		config.RestPort = restPort
	}

	return config
}
