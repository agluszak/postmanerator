package postman

type Collection struct {
	Name        string
	Description string
	Requests    []Request
	Folders     []Folder
	Structures  []StructureDefinition
}

type OriginalRequest struct {
	Method        string
	URL           string
	PayloadType   string
	PayloadRaw    string
	PayloadParams []KeyValuePair
	Headers       []KeyValuePair
}

type Request struct {
	OriginalRequest
	ID            string
	Name          string
	Description   string
	PathVariables []KeyValuePair
	Responses     []Response
	Tests         string
}

type Response struct {
	ID              string
	Name            string
	Status          string
	StatusCode      int
	Body            string
	Headers         []KeyValuePair
	OriginalRequest OriginalRequest
}

type Folder struct {
	ID          string
	Name        string
	Description string
	Folders     []Folder
	Requests    []Request
}

type StructureDefinition struct {
	Name        string
	Description string
	Fields      []StructureFieldDefinition
}

type StructureFieldDefinition struct {
	Name        string
	Description string
	Type        string
}

type KeyValuePair struct {
	Name        string
	Key         string
	Value       interface{}
	Description string
}
