package configuration

type Application struct {
	OAuth map[string]*struct {
		ClientId string
		SecretId string
	}
	Mysql map[string]*struct {
		Dns string
	}
	Domain map[string]*struct {
		Url string
	}
}
