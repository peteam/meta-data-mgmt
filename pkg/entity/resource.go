package entity

type Resource struct {
	Key          string `json:"key"`
	URN          string `json:"urn"`
	ID           string `json:"id"`
	ContentType  string `json:"contentType"`
	EntityStatus string `json:"entityStatus"`
	ProviderName string `json:"providerName"`
	Name         string `json:"name"`
	// catalogType - default value is: vodcatalog
	CatalogType string `json:"catalogType"`
	Body        string `json:"body"`
}
