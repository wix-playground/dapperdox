package spec

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
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
	Resource    *Resource // FIXME rename as Resource?
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

// -----------------------------------------------------------------------------

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
	// If Tags: [] is not defined, or empty, then no filtering or ordering takes place,#
	// and all API paths will be documented..
	for _, tag := range getTags(swaggerdoc.Spec()) {
		// Tag matching may not be as expected if multiple paths have the same TAG (which is technically permitted)
		var ok bool
		var ver interface{}

		for path, pathItem := range swaggerdoc.AllPaths() {

			var name string // Will only populate if Tagging used in spec. processMethod overrides if needed.
			name = tag.Description
			if name == "" {
				name = tag.Name
			}

			api := &API{
				ID:   TitleToKebab(name),
				Name: name,
				URL:  u,
				Info: &APIInfo,
			}

			if ver, ok = pathItem.Extensions["x-version"]; !ok {
				ver = "latest"
			}
			api.CurrentVersion = ver.(string)

			getMethods(tag, api, &api.Methods, &pathItem, path, ver.(string)) // Current version
			getVersions(tag, api, pathItem.Versions, path)                    // All versions

			// If API was populated, add to set
			if len(api.Methods) > 0 {
				APIs = append(APIs, *api) // All APIs (versioned within)
			}
		}
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

// -----------------------------------------------------------------------------

func getTags(specification *spec.Swagger) []spec.Tag {
	var tags []spec.Tag

	for _, tag := range specification.Tags {
		tags = append(tags, tag)
	}
	if len(tags) == 0 {
		tags = append(tags, spec.Tag{})
	}
	return tags
}

// -----------------------------------------------------------------------------

func getVersions(tag spec.Tag, api *API, versions map[string]spec.PathItem, path string) {
	if versions == nil {
		return
	}
	api.Versions = make(map[string][]Method)

	for v, pi := range versions {
		logger.Tracef(nil, "Process version %s\n", v)
		var method []Method
		getMethods(tag, api, &method, &pi, path, v)
		api.Versions[v] = method
	}
}

// -----------------------------------------------------------------------------

func getMethods(tag spec.Tag, api *API, methods *[]Method, pi *spec.PathItem, path string, version string) {

	getMethod(tag, api, methods, version, pi, pi.Get, path, "get")
	getMethod(tag, api, methods, version, pi, pi.Post, path, "post")
	getMethod(tag, api, methods, version, pi, pi.Put, path, "put")
	getMethod(tag, api, methods, version, pi, pi.Delete, path, "delete")
	getMethod(tag, api, methods, version, pi, pi.Head, path, "head")
	getMethod(tag, api, methods, version, pi, pi.Options, path, "options")
	getMethod(tag, api, methods, version, pi, pi.Patch, path, "patch")
}

// -----------------------------------------------------------------------------

func getMethod(tag spec.Tag, api *API, methods *[]Method, version string, pathitem *spec.PathItem, operation *spec.Operation, path, methodname string) {
	if operation == nil {
		return
	}
	// Filter and sort by matching current top-level tag with the operation tags.
	// If Tagging is not used by spec, then process each operation without filtering.
	taglen := len(operation.Tags)
	if taglen == 0 {
		if tag.Name != "" {
			logger.Tracef(nil, "Skipping %s - Operation does not contain a tag member, and tagging is in use.", operation.Summary)
			return
		}
		method := processMethod(api, pathitem, operation, path, methodname, version)
		*methods = append(*methods, *method)
	} else {
		for _, t := range operation.Tags {
			if tag.Name == "" || t == tag.Name {
				method := processMethod(api, pathitem, operation, path, methodname, version)
				*methods = append(*methods, *method)
			}
		}
	}
}

// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------

func processMethod(api *API, pathItem *spec.PathItem, o *spec.Operation, path, methodname string, version string) *Method {

	id := o.ID
	if id == "" {
		id = methodname
	}

	method := &Method{
		ID:          CamelToKebab(id),
		Name:        o.Summary,
		Description: o.Description,
		Method:      methodname,
		Path:        path,
		Responses:   make(map[int]Response),
		API:         api,
	}

	// If Tagging is not used by spec to select and order API paths to document, then
	// complete the missing names.
	// First try the vendor extension x-pathName, falling back to summary if not set.
	if pathname, ok := pathItem.Extensions["x-pathName"]; ok {
		api.Name = pathname.(string)
		api.ID = TitleToKebab(api.Name)
	}
	if api.Name == "" {
		name := o.Summary
		if name == "" {
			logger.Errorf(nil, "Error: Operation '%s' does not have an operationId or summary member.", id)
			os.Exit(1)
		}
		api.Name = name
		api.ID = TitleToKebab(name)
	}

	if ResourceList == nil {
		ResourceList = make(map[string]map[string]*Resource)
	}

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
		var vres *Resource

		// Discover if the resource is already declared, and pick it up
		// if it is (keyed on version number)
		if response.Schema != nil {
			if _, ok := ResourceList[version]; !ok {
				ResourceList[version] = make(map[string]*Resource)
			}
			var ok bool
			r, example_json := resourceFromSchema(response.Schema, method, nil) // May be thrown away

			r.Schema = jsonResourceToString(example_json, r.Type[0])

			// Look for a pre-declared resource with the response ID, and use that or create the first one...
			logger.Tracef(nil, "++ Resource version %s  ID %s\n", version, r.ID)
			if vres, ok = ResourceList[version][r.ID]; !ok {
				logger.Tracef(nil, "   - Creating new resource\n")
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
			Resource:    vres,
		}
	}

	if o.Responses.Default != nil {
		r, example_json := resourceFromSchema(o.Responses.Default.Schema, method, nil)
		if r != nil {

			r.Schema = jsonResourceToString(example_json, r.Type[0])

			logger.Tracef(nil, "++ Resource version %s  ID %s\n", version, r.ID)
			// Look for a pre-declared resource with the response ID, and use that or create the first one...
			var vres *Resource
			var ok bool
			if vres, ok = ResourceList[version][r.ID]; !ok {
				logger.Tracef(nil, "   - Creating new resource\n")
				vres = r
			}
			ResourceList[version][r.ID] = vres

			// Add to the compiled list of methods which use this resource
			vres.Methods = append(vres.Methods, *method)

			// Set the default response
			method.DefaultResponse = &Response{
				Description: o.Responses.Default.Description,
				Resource:    vres,
			}
		}
	}

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

	return method
}

// -----------------------------------------------------------------------------

func jsonResourceToString(jsonres map[string]interface{}, restype string) string {

	// If the resource is an array, then append json object to outer array, else serialise the object.
	var example []byte
	if strings.ToLower(restype) == "array" {
		var array_obj []map[string]interface{}
		array_obj = append(array_obj, jsonres)
		example, _ = json.MarshalIndent(array_obj, "", "    ")
	} else {
		example, _ = json.MarshalIndent(jsonres, "", "    ")
	}
	return string(example)
}

// -----------------------------------------------------------------------------

func checkPropertyType(s *spec.Schema) string {

	/*
	   (string) (len=12) "string_array": (spec.Schema) {
	    SchemaProps: (spec.SchemaProps) {
	     Description: (string) (len=16) "Array of strings",
	     Type: (spec.StringOrArray) (len=1 cap=1) { (string) (len=5) "array" },
	     Items: (*spec.SchemaOrArray)(0xc8205bb000)({
	      Schema: (*spec.Schema)(0xc820202480)({
	       SchemaProps: (spec.SchemaProps) {
	        Type: (spec.StringOrArray) (len=1 cap=1) { (string) (len=6) "string" },
	       },
	      }),
	     }),
	    },
	   }
	*/
	ptype := "primitive"

	if s.Type == nil {
		ptype = "object"
	}

	if s.Items != nil {
		ptype = "UNKNOWN"

		if s.Items.Schema != nil {
			s = s.Items.Schema
		} else {
			s = &s.Items.Schemas[0] // - Main schema [1] = Additional properties? See online swagger editior.
		}

		if s.Type == nil {
			ptype = "array of objects"
			if s.SchemaProps.Type != nil {
				ptype = "array of SOMETHING"
			}
		} else if s.Type.Contains("array") {
			ptype = "array of primitives"
		}
	}

	return ptype
}

// -----------------------------------------------------------------------------

func resourceFromSchema(s *spec.Schema, method *Method, fqNS []string) (*Resource, map[string]interface{}) {
	if s == nil {
		return nil, nil
	}

	stype := checkPropertyType(s)
	logger.Tracef(nil, "resourceFromSchema: Schema type: %s\n", stype)
	logger.Tracef(nil, "CHECK schema type and items\n")
	//spew.Dump(s)

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
	// TODO Check if this is valid?? TODO
	if s.Type == nil {
		s.Type = append(s.Type, "object")
	}

	if s.Items != nil {
		stringorarray := s.Type

		// Jump to nearest schema for items, depending on how it was declared
		if s.Items.Schema != nil { // items: { properties: {} }
			s = s.Items.Schema
			logger.Tracef(nil, "got s.Items.Schema for %s\n", s.Title)
		} else { // items: { $ref: "" }
			s = &s.Items.Schemas[0]
			logger.Tracef(nil, "got s.Items.Schemas[0] for %s\n", s.Title)
		}

		if s.Type == nil {
			logger.Tracef(nil, "Got array of objects? Name %s\n", s.Title)
			// Trying to fix issue/11
			s.Type = stringorarray
			// Trying to fix issue/11
		} else if s.Type.Contains("array") {
			logger.Tracef(nil, "Got array for %s\n", s.Title)
			s.Type = stringorarray
		} else if stringorarray.Contains("array") && len(s.Properties) == 0 {
			// if we get here then we can assume the type is supposed to be an array of primitives
			// Store the actual primitive type in the second element of the Type array.
			s.Type = spec.StringOrArray([]string{"array", s.Type[0]})
		}
		logger.Tracef(nil, "REMAP SCHEMA\n")
	}

	id := TitleToKebab(s.Title)

	if len(fqNS) > 0 && s.Type.Contains("array") {
		id = ""
	}

	if len(fqNS) == 0 && id == "" {
		logger.Errorf(nil, "Error: %s %s references a model definition that does not have a title member.", strings.ToUpper(method.Method), method.Path)
		os.Exit(1)
	}

	if strings.ToLower(s.Type[0]) != "object" {
		if strings.ToLower(s.Type[0]) == "array" {
			fqNSlen := len(fqNS)
			if fqNSlen > 0 {
				fqNS = append(fqNS[0:fqNSlen-1], fqNS[fqNSlen-1]+"[]")
			}
		}
	}

	myFQNS := append([]string{}, fqNS...)
	var chopped bool

	if len(id) == 0 && len(myFQNS) > 0 {
		id = myFQNS[len(myFQNS)-1]
		myFQNS = append([]string{}, myFQNS[0:len(myFQNS)-1]...)
		chopped = true
	}

	// If there is no description... the case where we have an array of objects. See issue/11
	var description string
	if s.Description != "" {
		description = string(github_flavored_markdown.Markdown([]byte(s.Description)))
	} else {
		description = s.Title
	}

	logger.Tracef(nil, "Create resource %s\n", id)
	r := &Resource{
		ID:          id,
		Title:       s.Title,
		Description: description,
		Type:        s.Type,
		Properties:  make(map[string]*Resource),
		FQNS:        myFQNS,
	}

	if s.Example != nil {
		example, err := json.MarshalIndent(&s.Example, "", "    ")
		if err != nil {
			logger.Errorf(nil, "Error encoding example json: %s", err)
		}
		r.Example = string(example)
	}

	if len(s.Enum) > 0 {
		for _, e := range s.Enum {
			r.Enum = append(r.Enum, fmt.Sprintf("%s", e))
		}
	}

	required := make(map[string]bool)
	json_representation := make(map[string]interface{})

	logger.Tracef(nil, "Call compileproperties...\n")
	compileproperties(s, r, method, id, required, json_representation, myFQNS, chopped)

	for allof := range s.AllOf {
		compileproperties(&s.AllOf[allof], r, method, id, required, json_representation, myFQNS, chopped)
	}

	logger.Tracef(nil, "resourceFromSchema done\n")

	return r, json_representation
}

// -----------------------------------------------------------------------------
// Takes a Schema object and adds properties to the Resource object.
// It uses the 'required' map to set when properties are required and builds a JSON
// representation of the resource.
//
func compileproperties(s *spec.Schema, r *Resource, method *Method, id string, required map[string]bool, json_rep map[string]interface{}, myFQNS []string, chopped bool) {

	// First, grab the required members
	for _, i := range s.Required {
		required[i] = true
	}

	// Now process the properties
	for name, property := range s.Properties {
		logger.Tracef(nil, "Process property name '%s'  Type %s\n", name, s.Properties[name].Type)
		newFQNS := append([]string{}, myFQNS...)

		if chopped && len(id) > 0 {
			newFQNS = append(newFQNS, id)
		}

		newFQNS = append(newFQNS, name)

		var json_resource map[string]interface{}

		logger.Tracef(nil, "A call resourceFromSchema for property %s\n", name)
		r.Properties[name], json_resource = resourceFromSchema(&property, method, newFQNS)

		if _, ok := required[name]; ok {
			r.Properties[name].Required = true
		}
		logger.Tracef(nil, "resource property %s type: %s\n", name, r.Properties[name].Type[0])

		if strings.ToLower(r.Properties[name].Type[0]) != "object" {
			// Arrays of objects need to be handled as a special case
			if strings.ToLower(r.Properties[name].Type[0]) == "array" {
				logger.Tracef(nil, "Processing an array property %s", name)
				if property.Items != nil {
					if property.Items.Schema != nil {

						logger.Tracef(nil, "ARRAY PROCESS %s:\n", name)

						// Some outputs (example schema, member description) are generated differently
						// if the array member references an object or a primitive type
						r.Properties[name].Description = property.Description

						// If here, we have no json_resource returned from resourceFromSchema, then the property
						// is an array of primitive, so construct either an array of string or array of object
						// as appropriate.
						if len(json_resource) > 0 {
							var array_obj []map[string]interface{}
							array_obj = append(array_obj, json_resource)
							json_rep[name] = array_obj
						} else {
							var array_obj []string
							// We stored the real type of the primitive in Type array index 1 (see the note in
							// resourceFromSchema).
							array_obj = append(array_obj, r.Properties[name].Type[1])
							json_rep[name] = array_obj
						}
					} else { // property.Items.Schema is NIL
						// Pretty sure this can never happen, due to the schema manipulation that
						// occurs in resourceFromSchema
					}
				} else {
					logger.Tracef(nil, "... and Items for %s are nil", name)
				}
			} else {
				json_rep[name] = r.Properties[name].Type[0]
			}
		} else {
			json_rep[name] = json_resource
		}
	}
}

// -----------------------------------------------------------------------------

func TitleToKebab(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, " ", "-", -1)
	return s
}

// -----------------------------------------------------------------------------

func CamelToKebab(s string) string {
	s = snaker.CamelToSnake(s)
	s = strings.Replace(s, "_", "-", -1)
	return s
}

// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
