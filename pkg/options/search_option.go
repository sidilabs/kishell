package options

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/sidilabs/kishell/pkg/config"
	"mime"
	"net/http"
	"text/template"
	"time"
)

// SearchParams represents attributes used to query for data.
type SearchParams struct {
	Index        string
	Size         int32
	WindowFilter string
	Zone         string
	Clause       string
	Older        int64
	Newer        int64
}

// ResponseData represents the response payload.
type ResponseData struct {
	Payload map[string]interface{}
}

const (
	kibanaVersionHeaderKey = "kbn-version"
	postContentType        = "application/x-ndjson"
	esSearchPath           = "/elasticsearch/_msearch"
	matchAllClause         = `{"match_all": {}}`
	queryClauseTemplate    = `{"query_string":{"query":"{{.Query}}","analyze_wildcard":true,"default_field":"*"}}`
	payloadTemplate        = `{"index":"{{.Index}}","ignore_unavailable":true,"preference":1569331617740}
{"version":true,"size":{{.Size}},"sort":[{"{{.WindowFilter}}":{"order":"desc","unmapped_type":"boolean"}}],"_source":{"excludes":[]},"aggs":{"2":{"date_histogram":{"field":"{{.WindowFilter}}","interval":"3h","time_zone":"{{.Zone}}","min_doc_count":1}}},"stored_fields":["*"],"script_fields":{},"query":{"bool":{"must":[{{.Clause}},{"range":{"{{.WindowFilter}}":{"gte":{{.Newer}},"lte":{{.Older}},"format":"epoch_millis"}}}],"filter":[],"should":[],"must_not":[]}},"highlight":{"pre_tags":["@kibana-highlighted-field@"],"post_tags":["@/kibana-highlighted-field@"],"fields":{"*":{}},"fragment_size":2147483647},"timeout":"30000ms"}
`
)

func (r *ResponseData) printAllSources() error {
	for _, element := range r.Payload {
		responses := element.([]interface{})
		for _, item := range responses {
			response := item.(map[string]interface{})
			hits := response["hits"].(map[string]interface{})["hits"].([]interface{})
			for _, hitItem := range hits {
				hit := hitItem.(map[string]interface{})
				source := hit["_source"]
				asJson, err := json.Marshal(source)
				if err != nil {
					return err
				}
				fmt.Println(string(asJson))
			}
		}
	}
	return nil
}

// Run the search option.
// Queries ES server for data.
// Prints the results in the stdout.
func (s *SearchCmd) Run(ctx *Context) error {
	err := ctx.Configuration.CheckEmpty()
	if err != nil {
		return err
	}
	clause := matchAllClause
	if len(s.Query) > 0 {
		out, err := buildFromTemplate("query", queryClauseTemplate, s)
		if err != nil {
			return err
		}
		clause = out.String()
	}

	server := ctx.Configuration.GetCurrentServer()
	if len(s.Server) > 0 {
		serverArg, ok := ctx.Configuration.FindServer(s.Server)
		if !ok {
			return fmt.Errorf("server '%s' is invalid", s.Server)
		}
		server = serverArg
	}
	role := ctx.Configuration.GetCurrentRole()

	currentTime := time.Now()
	olderTs, err := s.OlderAsTimestamp()
	if err != nil {
		return err
	}
	newerTs, err := s.NewerAsTimestamp()
	if err != nil {
		return err
	}
	searchParams := SearchParams{
		Index:        role.Index,
		Zone:         currentTime.Format("Z07:00"),
		WindowFilter: role.WindowFilter,
		Clause:       clause,
		Size:         s.Limit,
		Older:        olderTs,
		Newer:        newerTs,
	}
	payload, err := buildFromTemplate("payload", payloadTemplate, searchParams)
	if err != nil {
		return err
	}
	data, err := s.callApi(server, payload)
	if err != nil {
		return err
	}
	err = data.printAllSources()
	if err != nil {
		return err
	}
	return nil
}

func (s *SearchCmd) callApi(server config.Server, payload bytes.Buffer) (*ResponseData, error) {
	url := fmt.Sprintf("%s://%s:%s%s", server.Protocol, server.Hostname, server.GetPort(), esSearchPath)

	request, err := s.httpClient.NewRequest("POST", url, &payload)
	if err != nil {
		return nil, err
	}
	request.Header.Add(headers.ContentType, postContentType)
	request.Header.Add(kibanaVersionHeaderKey, server.KibanaVersion)
	if len(server.BasicAuth) > 0 {
		request.Header.Add(headers.Authorization, fmt.Sprintf("%s %s", "Basic", server.BasicAuth))
	}

	response, err := s.httpClient.Call(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return parseResponse(response)
}

func parseResponse(response *http.Response) (*ResponseData, error) {

	contentType, _, err := mime.ParseMediaType(response.Header.Get(headers.ContentType))
	if contentType == "application/json" {
		if response.StatusCode < 400 {
			var result map[string]interface{}
			err = json.NewDecoder(response.Body).Decode(&result)
			if err != nil {
				return nil, err
			}
			responseData := new(ResponseData)
			responseData.Payload = result
			return responseData, nil
		}
	}

	return nil, fmt.Errorf("invalid content type: %s", contentType)
}

func buildFromTemplate(name string, templateObj string, data interface{}) (bytes.Buffer, error) {
	queryTemplate, _ := template.New(name).Parse(templateObj)
	var out bytes.Buffer
	err := queryTemplate.Execute(&out, data)
	return out, err
}
