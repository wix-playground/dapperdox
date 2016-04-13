package spec

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	//"github.com/davecgh/go-spew/spew"
	"github.com/serenize/snaker"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/zxchris/go-swagger/spec"
	"github.com/zxchris/swaggerly/config"
	"github.com/zxchris/swaggerly/logger"
)

// APISet is a slice of API structs
type APISet []API

// APIs represents the parsed APIs
var APIs APISet
var APIInfo Info
var SecurityDefinitions map[string]SecurityScheme
var ResourceList map[string]map[string]*Resource // Version->ResourceName->Resource
var APIVersions map[string]APISet                // Version->APISet

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

type Info struct {
	Title       string
	Description string
}

// API represents an API
type API struct {
	ID             string
	Name           string
	URL            *url.URL
	Versions       map[string][]Method // All versions, keyed by version string.
	Methods        []Method            // The current version
	CurrentVersion string              // The latest version in operation for the API
	Info           *Info
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
	Schema      *Resource // FIXME rename as Resource?
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

	cfg, err := config.Get()
	if err != nil {
		logger.Errorf(nil, "error configuring app: %s", err)
	}

	fname := cfg.SpecFilename
	if !strings.HasPrefix(fname, "/") {
		fname = "/" + fname
	}

	swaggerdoc, err := loadSpec("http://" + host + fname)
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(swaggerdoc.Spec().Schemes[0] + "://" + swaggerdoc.Spec().Host)
	if err != nil {
		log.Fatal(err)
	}

	APIInfo.Title = swaggerdoc.Spec().Info.Title
	APIInfo.Description = swaggerdoc.Spec().Info.Description

	getSecurityDefinitions(swaggerdoc.Spec())

	// Use the top level TAGS to order the API resources/endpoints
	for _, tag := range swaggerdoc.Spec().Tags {
		api := API{
			ID:   titleToKebab(tag.Name),
			Name: tag.Name,
			URL:  u,
			Info: &APIInfo,
		}

		// Match up on tags: FIXME This does not work correctly if multiple paths have the same TAG (which is allowed)
		var ok bool
		var ver interface{}
		for p, o := range swaggerdoc.AllPaths() {
			if ver, ok = o.Extensions["x-version"]; !ok {
				ver = "latest"
			}
			api.CurrentVersion = ver.(string)

			getMethods(tag, &api, &api.Methods, o, p, ver.(string)) // Current version
			getVersions(tag, &api, o.Versions, p)                   // All versions
		}
		APIs = append(APIs, api) // All APIs (versioned within)
	}

	// Build a API map, grouping by version
	for _, api := range APIs {
		for v, _ := range api.Versions {
			if APIVersions == nil {
				APIVersions = make(map[string]APISet)
			}
			// Create copy of API and set Methods array to be correct for the version we are building
			napi := api
			napi.Methods = napi.Versions[v]
			napi.Versions = nil
			APIVersions[v] = append(APIVersions[v], napi) // Group APIs by version
		}
	}

}

func getVersions(tag spec.Tag, api *API, versions map[string]spec.PathItem, path string) {
	if versions == nil {
		return
	}
	api.Versions = make(map[string][]Method)

	for v, pi := range versions {
		logger.Tracef(nil, "Process version %s\n", v)
		var method []Method
		getMethods(tag, api, &method, pi, path, v)
		api.Versions[v] = method
	}
}

func getMethods(tag spec.Tag, api *API, methods *[]Method, pi spec.PathItem, path string, version string) {

	getMethod(tag, api, methods, version, pi.Get, path, "get")
	getMethod(tag, api, methods, version, pi.Post, path, "post")
	getMethod(tag, api, methods, version, pi.Put, path, "put")
	getMethod(tag, api, methods, version, pi.Delete, path, "delete")
	getMethod(tag, api, methods, version, pi.Head, path, "head")
	getMethod(tag, api, methods, version, pi.Options, path, "options")
	getMethod(tag, api, methods, version, pi.Patch, path, "patch")
}

func getMethod(tag spec.Tag, api *API, methods *[]Method, version string, o *spec.Operation, path, methodname string) {
	if o == nil {
		return
	}
	// Filter by tags or, if no Tags, build all APIs
	taglen := len(o.Tags)
	for _, t := range o.Tags {
		if taglen == 0 || t == tag.Name {
			method := processMethod(api, o, path, methodname, version)
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

func processMethod(api *API, o *spec.Operation, path, methodname string, version string) *Method {

	id := o.ID
	if id == "" {
		id = methodname
	}

	method := &Method{
		ID:          camelToKebab(id),
		Name:        o.Summary,
		Description: o.Description,
		Method:      methodname,
		Path:        path,
		Responses:   make(map[int]Response),
		API:         api,
	}

	if ResourceList == nil {
		ResourceList = make(map[string]map[string]*Resource)
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

	// Compile resources from response declaration
	for status, response := range o.Responses.StatusCodeResponses {
		//log.Printf("Got response schema (status %s):\n", status)
		//spew.Dump(response.Schema)

		var vres *Resource

		// Discover if the resource is already declared, and pick it up
		// if it is (keyed on version number)
		if response.Schema != nil {
			if _, ok := ResourceList[version]; !ok {
				ResourceList[version] = make(map[string]*Resource)
			}
			var ok bool
			r := resourceFromSchema(response.Schema, nil) // May be thrown away

			// Look for a pre-declared resource with the response ID, and use that or create the first one...
			log.Printf("++ Resource version %s  ID %s\n", version, r.ID)
			if vres, ok = ResourceList[version][r.ID]; !ok {
				log.Printf("   - Creating new resource\n")
				vres = r
			}
			ResourceList[version][r.ID] = vres

			// Compile a list of the methods which use this resource
			vres.Methods = append(vres.Methods, *method)

			// Add the resource to the method which uses it
			method.Resources = append(method.Resources, vres)

		}

		method.Responses[status] = Response{
			Description: response.Description,
			Schema:      vres,
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

	//api.Resources = make(map[string]*Resource) // List resources against API

	//	// List resources against method
	//	for _, r := range resources {
	//		method.Resources = append(method.Resources, r)
	//		//api.Resources[r.ID] = r
	//	}

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

	// XXX This is a bit of a hack, as it is possible for a response to be an array of
	//     objects, and it it possible to declare this in several ways:
	// 1. As :
	//      "schema": {
	//        "$ref": "model"
	//      }
	//      Where the model declares itself of type array (of objects)
	// 2. Or :
	//    "schema": {
	//        "type": "array",
	//        "items": {
	//            "$ref": "model"
	//        }
	//    }
	//
	//  In the second version, "items" points to a schema. So what we have done to align these
	//  two cases is to keep the top level "type" in the second case, and apply it to items.schema.Type,
	//  reseting our schema variable to items.schema.
	//

	//fmt.Printf("CHECK schema type and items\n")
	//spew.Dump(s)

	if s.Type == nil {
		s.Type = append(s.Type, "object")
	}

	if s.Items != nil {
		stringorarray := s.Type

		// EEK This is officially icky! See the Activities model in petstore. It declares "items": [ { } ] !!
		//     with an ARRAY
		if s.Items.Schema != nil {
			s = s.Items.Schema
		} else {
			s = &s.Items.Schemas[0]
		}
		if s.Type == nil {
			s.Type = stringorarray
		} else if s.Type.Contains("array") {
			s.Type = stringorarray
		}
		//fmt.Printf("REMAP SCHEMA\n")
		//spew.Dump(s)
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

	//log.Printf("expandSchema Type %s FQNS '%s'\n", s.Type, strings.Join(myFQNS, "."))
	//fmt.Printf("DUMP s.Properties\n")
	//spew.Dump(s.Properties)

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
			// Arrays of objects need to be handled as a special case
			if strings.ToLower(r.Properties[name].Type[0]) == "array" {
				if property.Items != nil {
					if property.Items.Schema != nil {

						//log.Printf("ARRAY PROCESS %s:\n", name)
						//spew.Dump(property.Items.Schema)

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
	// FIXME This explodes if there is no "type" member in the actual model definition - which is probably right
	//       for setting type of array in the model is a bit restrictive - better if set in the response decl. to
	//       say that the response for a status code is { "type":"array", "schema" : { "$ref": model } }
	//

	//fmt.Printf("DUMP s.Type\n")
	//spew.Dump(s.Type)
	if strings.ToLower(r.Type[0]) != "object" {
		if strings.ToLower(r.Type[0]) == "array" {
			var array_obj []map[string]interface{}
			array_obj = append(array_obj, json_representation)
			schema, _ := json.MarshalIndent(array_obj, "", "    ")
			r.Schema = string(schema)
		} else {
			r.Schema = r.Type[0]
		}
	} else {
		schema, err := json.MarshalIndent(json_representation, "", "    ")
		if err != nil {
			log.Printf("error encoding schema json: %s", err)
		}
		r.Schema = string(schema)
	}

	return r
}

// -----------------------------------------------------------------------------
// Take all the resources used by the method, and add them to the global resource
// list, merging the methods:w

func mergeResources(method *Method, version string) {

	//var ResourceList map[string]map[string]Resource // Version->ResourceName->Resource

	for _, r := range method.Resources {
		//ResourceList[version]
		method.Resources = append(method.Resources, r)
	}
}

// -----------------------------------------------------------------------------

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
