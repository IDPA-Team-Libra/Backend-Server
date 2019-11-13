package stockapi

type Api struct {
	ApiKey       string
	DefaultRoute string
	Headers      map[string]string
	Format       string
}

func NewApiInstance(apikey string, defaultRoutes string, headers map[string]string,
	format string) Api {
	api := Api{
		ApiKey:       apikey,
		DefaultRoute: defaultRoutes,
		Headers:      headers,
		Format:       format,
	}
	return api
}
