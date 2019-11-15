package postman

type collectionV210 struct {
	Info struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Schema      string `json:"schema"`
	} `json:"info"`
	Item []collectionV210Item `json:"item"`
	Auth *collectionV210Auth  `json:"auth,omitempty"`
}

type collectionV210Auth struct {
	Type   string                       `json:"type"`
	Bearer []collectionV210KeyValuePair `json:"bearer"`
}

type collectionV210Request struct {
	Method string                       `json:"method"`
	Header []collectionV210KeyValuePair `json:"header"`
	Body   struct {
		Mode       string                       `json:"mode"`
		Raw        string                       `json:"raw"`
		FormData   []collectionV210KeyValuePair `json:"formdata,omitempty"`
		UrlEncoded []collectionV210KeyValuePair `json:"urlencoded,omitempty"`
	} `json:"body"`
	Url struct {
		Raw      string                       `json:"raw"`
		Variable []collectionV210KeyValuePair `json:"variable"`
	} `json:"url"`
	Description string `json:"description"`
	Auth        *Auth  `json:"auth,omitempty"`
}

type collectionV210Response struct {
	Name            string                       `json:"name"`
	Status          string                       `json:"status"`
	Code            int                          `json:"code"`
	Header          []collectionV210KeyValuePair `json:"header"`
	Body            string                       `json:"body"`
	OriginalRequest *collectionV210Request       `json:"originalRequest,omitempty"`
}

type collectionV210Item struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Event       []collectionV210Event    `json:"event"`
	Item        []collectionV210Item     `json:"item"`
	Request     *collectionV210Request   `json:"request,omitempty"`
	Response    []collectionV210Response `json:"response"`
	Auth        *collectionV210Auth      `json:"auth,omitempty"`
}

type collectionV210Event struct {
	Listen string `json:"listen"`
	Script struct {
		Type string   `json:"type"`
		Exec []string `json:"exec"`
	} `json:"script"`
}

type collectionV210KeyValuePair struct {
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
}
