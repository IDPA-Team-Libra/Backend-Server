package stock

type Loader struct {
	apiKey string
}

func NewLoader(apiKey string) Loader {
	loader := Loader{
		apiKey: apiKey,
	}
	return loader
}
