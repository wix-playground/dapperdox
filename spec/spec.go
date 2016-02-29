package spec

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	//"github.com/companieshouse/swaggerly/logger"
	"github.com/companieshouse/go-swagger/spec"
	//"github.com/davecgh/go-spew/spew"
	"github.com/serenize/snaker"
	"github.com/shurcooL/github_flavored_markdown"
)

// APISet is a slice of API structs
type APISet []API

// APIs represents the parsed APIs
var APIs APISet
var SecurityDefinitions map[string]SecurityScheme

// GetByName returns an API by name
func (a APISet) GetByName(name string) *API {
	for _, a := range APIs {
		if a.Name == name {
			return &a
		}
	}
	return nil
}

// GetByID returns an API by ID
func (a APISet) GetByID(id string) *API {
	for _, a := range APIs {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

// API represents an API
// XXX Versisons will only contain other (not the current) version if the x-version (singular) extension is implemented XXX
type API struct {
	ID             string
	Name           string
	Versions       map[string][]Method // All versions, keyed by version string.
	Methods        []Method            // The current version
	URL            *url.URL
	CurrentVersion string
}

type Version struct {
	Version string
	Methods []Method
}

type OAuth2Scheme struct {
	OAuth2Flow       string
	AuthorizationUrl string
	TokenUrl         string
	Scopes           map[string]string
}

type SecurityScheme struct {
	IsApiKey      bool
	IsBasic       bool
	IsOAuth2      bool
	Type          string
	Description   string
	ParamName     string
	ParamLocation string
	OAuth2Scheme
}

type Security struct {
	Scheme *SecurityScheme
	Scopes map[string]string
}

// Method represents an API method
type Method struct {
	ID              string
	Name            string
	Description     string
	Method          string
	Path            string
	PathParams      []Parameter
	QueryParams     []Parameter
	HeaderParams    []Parameter
	BodyParam       *Parameter
	FormParams      []Parameter
	Responses       map[int]Response
	DefaultResponse *Response
	Resources       []*Resource
	Security        map[string]Security
	API             *API
}

// Parameter represents an API method parameter
type Parameter struct {
	Name        string
	In          string
	Required    bool
	Description string
	Type        string
	Enum        []string
}

// Response represents an API method response
type Response struct {
	Description string
	Schema      *Resource
}

// Resource represents an API resource
type Resource struct {
	ID          string
	FQNS        []string
	Title       string
	Description string
	Example     string
	Schema      string
	Type        []string
	Properties  map[string]*Resource
	Required    bool
	Methods     []Method
	Enum        []string
}

// Load loads API specs from the supplied host (usually local!)
func Load(host string) {
	swaggerdoc, err := loadSpec("http://" + host + "/spec/swagger.json")
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(swaggerdoc.Spec().Schemes[0] + "://" + swaggerdoc.Spec().Host)
	if err != nil {
		log.Fatal(err)
	}

	getSecurityDefinitions(swaggerdoc.Spec())

	// Use the top level TAGS to order the API resources/endpoints
	for _, tag := range swaggerdoc.Spec().Tags {
		api := API{
			ID:   titleToKebab(tag.Name),
			Name: tag.Name,
			URL:  u,
		}

		// Match up on tags:
		var ok bool
		var ver interface{}
		for p, o := range swaggerdoc.AllPaths() {
			if ver, ok = o.Extensions["x-version"]; !ok {
				ver = "latest"
			}
			api.CurrentVersion = ver.(string)

			getMethods(tag, &api, &api.Methods, o, p) // Current version
			getVersions(tag, &api, o.Versions, p)     // All versions
		}

		APIs = append(APIs, api)
	}
}

func getVersions(tag spec.Tag, api *API, versions map[string]spec.PathItem, path string) {
	if versions == nil {
		return
	}
	api.Versions = make(map[string][]Method)

	for v, pi := range versions {
		fmt.Printf("Process version %s\n", v)
		var method []Method
		getMethods(tag, api, &method, pi, path)
		api.Versions[v] = method
	}
}

func getMethods(tag spec.Tag, api *API, methods *[]Method, pi spec.PathItem, path string) {

	getMethod(tag, api, methods, pi.Get, path, "get")
	getMethod(tag, api, methods, pi.Post, path, "post")
	getMethod(tag, api, methods, pi.Put, path, "put")
	getMethod(tag, api, methods, pi.Delete, path, "delete")
	getMethod(tag, api, methods, pi.Head, path, "head")
	getMethod(tag, api, methods, pi.Options, path, "options")
	getMethod(tag, api, methods, pi.Patch, path, "patch")
}

func getMethod(tag spec.Tag, api *API, methods *[]Method, o *spec.Operation, path, methodname string) {
	if o == nil {
		return
	}
	for _, t := range o.Tags {
		if t == tag.Name {
			method := processMethod(api, o, path, methodname)
			*methods = append(*methods, *method)
		}
	}
}

func getSecurityDefinitions(spec *spec.Swagger) {

	if SecurityDefinitions == nil {
		SecurityDefinitions = make(map[string]SecurityScheme)
	}

	for n, d := range spec.SecurityDefinitions {
		stype := d.Type

		def := &SecurityScheme{
			Description:   d.Description,
			Type:          stype,  // basic, apiKey or oauth2
			ParamName:     d.Name, // name of header to be used if ParamLocation is 'header'
			ParamLocation: d.In,   // Either query or header
		}

		if stype == "apiKey" {
			def.IsApiKey = true
		}
		if stype == "basic" {
			def.IsBasic = true
		}
		if stype == "oauth2" {
			def.IsOAuth2 = true
			def.OAuth2Flow = d.Flow                   // implicit, password (explicit) application or accessCode
			def.AuthorizationUrl = d.AuthorizationURL // Only for implicit or accesscode flow
			def.TokenUrl = d.TokenURL                 // Only for implicit, accesscode or password flow
			def.Scopes = make(map[string]string)
			for s, n := range d.Scopes {
				def.Scopes[s] = n
			}
		}

		SecurityDefinitions[n] = *def
	}
}

func processMethod(api *API, o *spec.Operation, path, methodname string) *Method {

	method := &Method{
		ID:          camelToKebab(o.ID),
		Name:        o.Summary,
		Description: o.Description,
		Method:      methodname,
		Path:        path,
		Responses:   make(map[int]Response),
		API:         api,
	}

	resources := make(map[string]*Resource)

	for _, param := range o.Parameters {
		p := Parameter{
			Name:        param.Name,
			In:          param.In,
			Description: param.Description,
			Type:        param.Type,
			Required:    param.Required,
		}
		switch strings.ToLower(param.In) {
		case "form":
			method.FormParams = append(method.FormParams, p)
		case "path":
			method.PathParams = append(method.PathParams, p)
		case "body":
			method.BodyParam = &p
		case "header":
			method.HeaderParams = append(method.HeaderParams, p)
		case "query":
			method.QueryParams = append(method.QueryParams, p)
		}
		switch strings.ToLower(param.Type) {
		case "enum":
			for _, e := range param.Enum {
				p.Enum = append(p.Enum, fmt.Sprintf("%s", e))
			}
		}
	}

	for status, response := range o.Responses.StatusCodeResponses {

		r := resourceFromSchema(response.Schema, nil)

		if response.Schema != nil {
			r.Methods = append(r.Methods, *method)
			resources[r.ID] = r
		}

		method.Responses[status] = Response{
			Description: response.Description,
			Schema:      r,
		}
	}

	if o.Responses.Default != nil {
		r := resourceFromSchema(o.Responses.Default.Schema, nil)
		if r != nil {
			r.Methods = append(r.Methods, *method)
			resources[r.ID] = r
		}

		method.DefaultResponse = &Response{
			Description: o.Responses.Default.Description,
			Schema:      r,
		}
	}

	for _, r := range resources {
		method.Resources = append(method.Resources, r)
	}

	// Lookup security reference against SecurityDefinitions
	// TODO FIXME If no Security given from operation, then the global defaults are appled. CHECK THIS IS TRUE!

	method.Security = make(map[string]Security)

	for _, sec := range o.Security {
		for n, scopes := range sec {
			// Lookup security name in definitions
			if scheme, ok := SecurityDefinitions[n]; ok {

				// Add security to method
				method.Security[n] = Security{
					Scheme: &scheme,
					Scopes: make(map[string]string),
				}

				// Populate method specific scopes by cross referencing SecurityDefinitions
				for _, scope := range scopes {
					if scope_desc, ok := scheme.Scopes[scope]; ok {
						method.Security[n].Scopes[scope] = scope_desc
					}
				}
			}
		}
	}

	//fmt.Printf("DUMPING Method Security\n")
	//spew.Dump(method.Security)

	return method
}

// -----------------------------------------------------------------------------

func resourceFromSchema(s *spec.Schema, fqNS []string) *Resource {
	if s == nil {
		return nil
	}

	myFQNS := append([]string{}, fqNS...)
	id := titleToKebab(s.Title)

	var chopped bool
	if len(id) == 0 && len(myFQNS) > 0 {
		id = myFQNS[len(myFQNS)-1]
		myFQNS = append([]string{}, myFQNS[0:len(myFQNS)-1]...)
		chopped = true
	}

	r := &Resource{
		ID:          id,
		Title:       s.Title,
		Description: string(github_flavored_markdown.Markdown([]byte(s.Description))),
		Type:        s.Type,
		Properties:  make(map[string]*Resource),
		FQNS:        myFQNS,
	}

	if s.Example != nil {
		example, err := json.MarshalIndent(&s.Example, "", "    ")
		if err != nil {
			log.Printf("error encoding example json: %s", err)
		}
		r.Example = string(example)
	}

	if len(s.Enum) > 0 {
		for _, e := range s.Enum {
			r.Enum = append(r.Enum, fmt.Sprintf("%s", e))
		}
	}

	required := make(map[string]bool)
	for _, i := range s.Required {
		required[i] = true
	}

	json_representation := make(map[string]interface{})

	// log.Printf("expandSchema Type %s FQNS %s\n", s.Type, strings.Join(myFQNS, "."))

	for name, property := range s.Properties {
		//log.Printf("Process property name '%s'  Type %s\n", name, s.Properties[name].Type)
		newFQNS := append([]string{}, myFQNS...)
		if chopped && len(id) > 0 {
			newFQNS = append(newFQNS, id)
		}
		newFQNS = append(newFQNS, name)

		// log.Printf("Recurse into resourceFromSchema for property name '%s'\n", name)
		r.Properties[name] = resourceFromSchema(&property, newFQNS)

		if _, ok := required[name]; ok {
			r.Properties[name].Required = true
		}

		// FIXME this is as nasty as it looks...
		if strings.ToLower(r.Properties[name].Type[0]) != "object" {
			// log.Printf("Type is '%s'\n", r.Properties[name].Type[0])

			// Arrays of objects need to be handled as a special case
			if strings.ToLower(r.Properties[name].Type[0]) == "array" {
				if property.Items != nil {
					if property.Items.Schema != nil {

						// log.Printf("ARRAY PROCESS %s:\n", name)
						// spew.Dump(property.Items.Schema)

						// Add [] to end of fully qualified name space
						xFQNS := append([]string{}, newFQNS...)
						if len(xFQNS) > 0 {
							xFQNS = append(newFQNS[0:len(newFQNS)-1], newFQNS[len(newFQNS)-1]+"[]")
						}

						r.Properties[name] = resourceFromSchema(property.Items.Schema, xFQNS)

						// log.Printf("Generated Properties:\n")
						// spew.Dump(r.Properties[name])

						// Some outputs (example schema, member description) are generated differently
						// if the array member references an object or a primitive type
						var example_sch string
						if strings.ToLower(r.Properties[name].Type[0]) == "object" {
							example_sch = r.Properties[name].Schema
						} else {
							example_sch = "\"" + r.Properties[name].Type[0] + "\""
							r.Properties[name].Description = property.Description
						}

						var f interface{}
						_ = json.Unmarshal([]byte("["+example_sch+"]"), &f)
						json_representation[name] = f

						// Override type to reflect it is an array
						r.Properties[name].Type[0] = "array[" + r.Properties[name].Type[0] + "]"
					}
				}
			} else {
				json_representation[name] = r.Properties[name].Schema
			}
		} else {
			var f interface{}
			_ = json.Unmarshal([]byte(r.Properties[name].Schema), &f)
			json_representation[name] = f
		}
	}

	// Build element of resource schema example
	// FIXME also as nasty as it looks
	if strings.ToLower(r.Type[0]) != "object" {
		r.Schema = r.Type[0]
	} else {
		schema, err := json.MarshalIndent(json_representation, "", "    ")
		if err != nil {
			log.Printf("error encoding schema json: %s", err)
		}
		r.Schema = string(schema)
	}

	return r
}

func titleToKebab(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, " ", "-", -1)
	return s
}

func camelToKebab(s string) string {
	s = snaker.CamelToSnake(s)
	s = strings.Replace(s, "_", "-", -1)
	return s
}

func loadSpec(url string) (*spec.Document, error) {
	spec, err := spec.Load(url)
	if err != nil {
		return nil, err
	}

	spec, err = spec.Expanded()
	if err != nil {
		return nil, err
	}

	return spec, err
}
