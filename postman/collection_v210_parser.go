package postman

import (
	"encoding/json"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type CollectionV210Parser struct{}

func (p *CollectionV210Parser) CanParse(contents []byte) bool {
	return true // TODO: compare json schema value
}

func (p *CollectionV210Parser) Parse(contents []byte, options BuilderOptions) (Collection, error) {
	src := collectionV210{}
	if err := json.Unmarshal(contents, &src); err != nil {
		return Collection{}, err
	}
	return p.buildCollection(src, options)
}

func (p *CollectionV210Parser) buildCollection(src collectionV210, options BuilderOptions) (Collection, error) {

	auth, err := p.parseAuth(src.Auth)
	if err != nil {
		return Collection{}, err
	}
	collection := Collection{
		Name:        src.Info.Name,
		Description: src.Info.Description,
		Requests:    make([]Request, 0),
		Folders:     make([]Folder, 0),
		Structures:  make([]StructureDefinition, 0),
		Auth:        auth,
	}

	rootItem := Folder{}
	if err := p.computeItem(&rootItem, src.Item, options); err != nil {
		return collection, fmt.Errorf("failed to build request: %v", err)
	}

	collection.Requests = rootItem.Requests
	collection.Folders = rootItem.Folders

	return collection, nil
}

func (p *CollectionV210Parser) parseAuth(auth *collectionV210Auth) (*Auth, error) {
	if auth == nil {
		return nil, nil
	}

	params := make([]KeyValuePair, 0)
	for _, pair := range auth.Bearer {
		params = append(params, KeyValuePair{
			Name:  pair.Key,
			Key:   pair.Key,
			Value: pair.Value,
		})
	}

	ret := Auth{
		Type:   auth.Type,
		Params: params,
	}

	switch auth.Type {
	case "bearer":
		if len(auth.Bearer) != 1 || auth.Bearer[0].Key != "token" {
			return nil, fmt.Errorf("incorrect auth structure for type: %v", auth.Type)
		}
	default:
		return nil, fmt.Errorf("unsupported auth type: %v", auth.Type)
	}

	return &ret, nil
}

func (p *CollectionV210Parser) parseOriginalRequest(request *collectionV210Request, options BuilderOptions, parentAuth *Auth) OriginalRequest {
	return OriginalRequest{
		Method:        request.Method,
		URL:           request.Url.Raw,
		PayloadType:   request.Body.Mode,
		PayloadRaw:    request.Body.Raw,
		PayloadParams: p.parseRequestPayloadParams(*request),
		Headers:       p.parseRequestHeaders(request.Header, options),
		Auth:          parentAuth,
	}
}

func (p *CollectionV210Parser) computeItem(parentFolder *Folder, items []collectionV210Item, options BuilderOptions) error {
	for _, item := range items {
		auth, err := p.parseAuth(item.Auth)
		if err != nil {
			return err
		}
		if item.Request == nil { // item is a folder
			folder := Folder{
				ID:          uuid.NewV4().String(),
				Description: item.Description,
				Name:        item.Name,
				Auth:        auth,
			}
			if err := p.computeItem(&folder, item.Item, options); err != nil {
				return err
			}
			parentFolder.Folders = append(parentFolder.Folders, folder)
		} else { // item is a request
			var thisAuth *Auth
			if auth != nil {
				thisAuth = auth
			} else {
				thisAuth = parentFolder.Auth
			}
			request := Request{
				OriginalRequest: OriginalRequest{
					Method:        item.Request.Method,
					URL:           item.Request.Url.Raw,
					PayloadType:   item.Request.Body.Mode,
					PayloadRaw:    item.Request.Body.Raw,
					PayloadParams: p.parseRequestPayloadParams(*item.Request),
					Headers:       p.parseRequestHeaders(item.Request.Header, options),
					Auth:          thisAuth,
				},
				ID:            uuid.NewV4().String(),
				Name:          item.Name,
				Description:   item.Request.Description,
				Tests:         p.parseRequestTests(item),
				PathVariables: p.parseRequestPathVariables(item),
				Responses:     p.parseRequestResponses(item.Response, options, thisAuth),
			}
			parentFolder.Requests = append(parentFolder.Requests, request)
		}
	}

	return nil
}

func (p *CollectionV210Parser) parseRequestTests(item collectionV210Item) string {
	for _, event := range item.Event {
		if event.Listen == "test" {
			return strings.Join(event.Script.Exec, "\n")
		}
	}
	return ""
}

func (p *CollectionV210Parser) parseRequestPathVariables(item collectionV210Item) []KeyValuePair {
	pathVariables := make([]KeyValuePair, 0)

	for _, variable := range item.Request.Url.Variable {
		pathVariables = append(pathVariables, KeyValuePair{
			Name:        variable.Key,
			Key:         variable.Key,
			Value:       variable.Value,
			Description: variable.Description,
		})
	}

	return pathVariables
}

func (p *CollectionV210Parser) parseRequestPayloadParams(request collectionV210Request) []KeyValuePair {
	payloadParams := make([]KeyValuePair, 0)

	keyValuePairCollection := make([]collectionV210KeyValuePair, 0)
	switch request.Body.Mode {
	case "urlencoded":
		keyValuePairCollection = request.Body.UrlEncoded
	case "formdata":
		keyValuePairCollection = request.Body.FormData
	}

	for _, pair := range keyValuePairCollection {
		payloadParams = append(payloadParams, KeyValuePair{
			Name:        pair.Key,
			Key:         pair.Key,
			Value:       pair.Value,
			Description: pair.Description,
		})
	}

	return payloadParams
}

func (p *CollectionV210Parser) parseRequestHeaders(headers []collectionV210KeyValuePair, options BuilderOptions) []KeyValuePair {
	parsedHeaders := make([]KeyValuePair, 0)

	for _, header := range headers {
		if containsString(options.IgnoredRequestHeaders, header.Key) {
			continue
		}
		parsedHeaders = append(parsedHeaders, KeyValuePair{
			Name:        header.Key,
			Key:         header.Key,
			Value:       header.Value,
			Description: header.Description,
		})
	}

	return parsedHeaders
}

func (p *CollectionV210Parser) parseRequestResponses(responses []collectionV210Response, options BuilderOptions, parentAuth *Auth) []Response {
	parsedResponses := make([]Response, 0)

	for _, resp := range responses {
		parsedResponses = append(parsedResponses, Response{
			ID:              uuid.NewV4().String(),
			Name:            resp.Name,
			Body:            resp.Body,
			Status:          resp.Status,
			StatusCode:      resp.Code,
			Headers:         p.parseResponseHeaders(resp.Header, options),
			OriginalRequest: p.parseOriginalRequest(resp.OriginalRequest, options, parentAuth),
		})
	}

	return parsedResponses
}

func (p *CollectionV210Parser) parseResponseHeaders(headers []collectionV210KeyValuePair, options BuilderOptions) []KeyValuePair {
	parsedHeaders := make([]KeyValuePair, 0)

	for _, header := range headers {
		if containsString(options.IgnoredResponseHeaders, header.Key) {
			continue
		}
		parsedHeaders = append(parsedHeaders, KeyValuePair{
			Name:        header.Key,
			Key:         header.Key,
			Value:       header.Value,
			Description: header.Description,
		})
	}
	return parsedHeaders
}

func containsString(target []string, symbol string) bool {
	for _, element := range target {
		if element == symbol {
			return true
		}
	}
	return false
}
