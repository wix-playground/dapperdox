package spec

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/serenize/snaker"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/zxchris/go-swagger/spec"
	"github.com/zxchris/swaggerly/config"
	"github.com/zxchris/swaggerly/logger"
)

type APISpecification struct {
	ID      string
	APIs    APISet // APIs represents the parsed APIs
	APIInfo Info

	SecurityDefinitions map[string]SecurityScheme
	ResourceList        map[string]map[string]*Resource // Version->ResourceName->Resource
	APIVersions         map[string]APISet               // Version->APISet
}

var APISuite map[string]*APISpecification

// GetByName returns an API by name
func (c *APISpecification) GetByName(name string) *API {
	for _, a := range c.APIs {
		if a.Name == name {
			return &a
		}
	}
	return nil
}

// GetByID returns an API by ID
func (c *APISpecification) GetByID(id string) *API {
	for _, a := range c.APIs {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

type APISet []API

type Info struct {
	Title       string
	Description string
}

// API represents an API
type API struct {
	ID                     string
	Name                   string
	URL                    *url.URL
	MethodNavigationByType bool
	Versions               map[string][]Method // All versions, keyed by version string.
	Methods                []Method            // The current version
	CurrentVersion         string              // The latest version in operation for the API
	Info                   *Info
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
	NavigationName  string
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
	Description string
	In          string
	Required    bool
	Type        string
	Enum        []string
	Resource    *Resource // For "in body" parameters
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

// -----------------------------------------------------------------------------

func LoadSpecifications(host string, collapse bool) error {

	if APISuite == nil {
		APISuite = make(map[string]*APISpecification)
	}

	cfg, err := config.Get()
	if err != nil {
		logger.Errorf(nil, "error configuring app: %s", err)
		return err
	}

	for _, specFilename := range cfg.SpecFilename {

		var ok bool
		var specification *APISpecification

		if specification, ok = APISuite[""]; !ok || !collapse {
			specification = &APISpecification{}
		}

		err = specification.Load(specFilename, host)
		if err != nil {
			return err
		}

		if collapse {
			//specification.ID = "api"
		}

		APISuite[specification.ID] = specification
	}

	return nil
}

// -----------------------------------------------------------------------------
// Load loads API specs from the supplied host (usually local!)
func (c *APISpecification) Load(specFilename string, host string) error {

	if !strings.HasPrefix(specFilename, "/") {
		specFilename = "/" + specFilename
	}

	swaggerdoc, err := loadSpec("http://" + host + specFilename) // XXX Is there a confusion here between SpecDir and SpecFilename
	if err != nil {
		return err
	}

	basePath := swaggerdoc.Spec().BasePath
	basePathLen := len(basePath)
	// Ignore basepath if it is a single '/'
	if basePathLen == 1 && basePath[0] == '/' {
		basePathLen = 0
	}

	u, err := url.Parse(swaggerdoc.Spec().Schemes[0] + "://" + swaggerdoc.Spec().Host)
	if err != nil {
		return err
	}

	c.APIInfo.Description = string(github_flavored_markdown.Markdown([]byte(swaggerdoc.Spec().Info.Description)))
	c.APIInfo.Title = swaggerdoc.Spec().Info.Title

	logger.Tracef(nil, "Parse OpenAPI specification '%s'\n", c.APIInfo.Title)

	c.ID = TitleToKebab(c.APIInfo.Title)

	c.getSecurityDefinitions(swaggerdoc.Spec())

	methodNavByType := false // Should methods in the navigation be presented by type (GET, POST) or name (string)?
	if byname, ok := swaggerdoc.Spec().Extensions["x-methods-by-type"]; ok {
		if byname.(bool) {
			methodNavByType = true
		}
	}

	// Use the top level TAGS to order the API resources/endpoints
	// If Tags: [] is not defined, or empty, then no filtering or ordering takes place,#
	// and all API paths will be documented..
	for _, tag := range getTags(swaggerdoc.Spec()) {
		logger.Tracef(nil, "  In tag loop...\n")
		// Tag matching may not be as expected if multiple paths have the same TAG (which is technically permitted)
		var ok bool
		var ver interface{}

		//logger.Printf(nil, "DUMP OF ENTIRE SWAGGER SPEC\n")
		//spew.Dump(swaggerdoc)

		var api *API
		groupingByTag := false

		if tag.Name != "" {
			groupingByTag = true
		}

		var name string // Will only populate if Tagging used in spec. processMethod overrides if needed.
		name = tag.Description
		if name == "" {
			name = tag.Name
		}
		logger.Tracef(nil, "    - %s\n", name)

		// If we're grouping by TAGs, then build the API at the tag level
		if groupingByTag {
			api = &API{
				ID:   TitleToKebab(name),
				Name: name,
				URL:  u,
				Info: &c.APIInfo,
				MethodNavigationByType: methodNavByType,
			}
		}

		for path, pathItem := range swaggerdoc.AllPaths() {
			logger.Tracef(nil, "    In path loop...\n")

			if basePathLen > 0 {
				path = basePath + path
			}

			if !groupingByTag {
				api = &API{
					ID:   TitleToKebab(name),
					Name: name,
					URL:  u,
					Info: &c.APIInfo,
					MethodNavigationByType: methodNavByType,
				}
			}

			if ver, ok = pathItem.Extensions["x-version"]; !ok {
				ver = "latest"
			}
			api.CurrentVersion = ver.(string)

			c.getMethods(tag, api, &api.Methods, &pathItem, path, ver.(string)) // Current version
			c.getVersions(tag, api, pathItem.Versions, path)                    // All versions

			// If API was populated (will not be if tags do not match), add to set
			if !groupingByTag && len(api.Methods) > 0 {
				logger.Tracef(nil, "    + Adding %s\n", name)
				c.APIs = append(c.APIs, *api) // All APIs (versioned within)
			}
		}

		if groupingByTag && len(api.Methods) > 0 {
			logger.Tracef(nil, "    + Adding %s\n", name)
			c.APIs = append(c.APIs, *api) // All APIs (versioned within)
		}
	}

	// Build a API map, grouping by version
	for _, api := range c.APIs {
		for v, _ := range api.Versions {
			if c.APIVersions == nil {
				c.APIVersions = make(map[string]APISet)
			}
			// Create copy of API and set Methods array to be correct for the version we are building
			napi := api
			napi.Methods = napi.Versions[v]
			napi.Versions = nil
			c.APIVersions[v] = append(c.APIVersions[v], napi) // Group APIs by version
		}
	}

	return nil
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

func (c *APISpecification) getVersions(tag spec.Tag, api *API, versions map[string]spec.PathItem, path string) {
	if versions == nil {
		return
	}
	api.Versions = make(map[string][]Method)

	for v, pi := range versions {
		logger.Tracef(nil, "Process version %s\n", v)
		var method []Method
		c.getMethods(tag, api, &method, &pi, path, v)
		api.Versions[v] = method
	}
}

// -----------------------------------------------------------------------------

func (c *APISpecification) getMethods(tag spec.Tag, api *API, methods *[]Method, pi *spec.PathItem, path string, version string) {

	c.getMethod(tag, api, methods, version, pi, pi.Get, path, "get")
	c.getMethod(tag, api, methods, version, pi, pi.Post, path, "post")
	c.getMethod(tag, api, methods, version, pi, pi.Put, path, "put")
	c.getMethod(tag, api, methods, version, pi, pi.Delete, path, "delete")
	c.getMethod(tag, api, methods, version, pi, pi.Head, path, "head")
	c.getMethod(tag, api, methods, version, pi, pi.Options, path, "options")
	c.getMethod(tag, api, methods, version, pi, pi.Patch, path, "patch")
}

// -----------------------------------------------------------------------------

func (c *APISpecification) getMethod(tag spec.Tag, api *API, methods *[]Method, version string, pathitem *spec.PathItem, operation *spec.Operation, path, methodname string) {
	if operation == nil {
		return
	}
	// Filter and sort by matching current top-level tag with the operation tags.
	// If Tagging is not used by spec, then process each operation without filtering.
	taglen := len(operation.Tags)
	logger.Tracef(nil, "  Operation tag length: %d", taglen)
	if taglen == 0 {
		if tag.Name != "" {
			logger.Tracef(nil, "Skipping %s - Operation does not contain a tag member, and tagging is in use.", operation.Summary)
			return
		}
		method := c.processMethod(api, pathitem, operation, path, methodname, version)
		*methods = append(*methods, *method)
	} else {
		logger.Tracef(nil, "    > Check tags")
		for _, t := range operation.Tags {
			logger.Tracef(nil, "      - Compare tag '%s' with '%s'\n", tag.Name, t)
			if tag.Name == "" || t == tag.Name {
				method := c.processMethod(api, pathitem, operation, path, methodname, version)
				*methods = append(*methods, *method)
			}
		}
	}
}

// -----------------------------------------------------------------------------

func (c *APISpecification) getSecurityDefinitions(spec *spec.Swagger) {

	if c.SecurityDefinitions == nil {
		c.SecurityDefinitions = make(map[string]SecurityScheme)
	}

	for n, d := range spec.SecurityDefinitions {
		stype := d.Type

		def := &SecurityScheme{
			Description:   string(github_flavored_markdown.Markdown([]byte(d.Description))),
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

		c.SecurityDefinitions[n] = *def
	}
}

// -----------------------------------------------------------------------------

func (c *APISpecification) processMethod(api *API, pathItem *spec.PathItem, o *spec.Operation, path, methodname string, version string) *Method {

	id := o.ID
	if id == "" {
		id = methodname
	}

	method := &Method{
		ID:          CamelToKebab(id),
		Name:        o.Summary,
		Description: string(github_flavored_markdown.Markdown([]byte(o.Description))),
		Method:      methodname,
		Path:        path,
		Responses:   make(map[int]Response),
		API:         api,
	}

	if navname, ok := o.Extensions["x-navigation-name"]; ok {
		method.NavigationName = navname.(string)
	} else {
		if api.MethodNavigationByType {
			method.NavigationName = method.Method
		} else {
			method.NavigationName = method.Name
		}
	}

	// If Tagging is not used by spec to select, group and order API paths to document, then
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

	if c.ResourceList == nil {
		c.ResourceList = make(map[string]map[string]*Resource)
	}

	for _, param := range o.Parameters {
		p := Parameter{
			Name:        param.Name,
			In:          param.In,
			Description: string(github_flavored_markdown.Markdown([]byte(param.Description))),
			Type:        param.Type,
			Required:    param.Required,
		}
		switch strings.ToLower(param.In) {
		case "form":
			method.FormParams = append(method.FormParams, p)
		case "path":
			method.PathParams = append(method.PathParams, p)
		case "body":
			p.Resource = c.resourceFromSchema(param.Schema, method, nil)
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
			if _, ok := c.ResourceList[version]; !ok {
				c.ResourceList[version] = make(map[string]*Resource)
			}
			var ok bool
			r := c.resourceFromSchema(response.Schema, method, nil) // May be thrown away

			// Look for a pre-declared resource with the response ID, and use that or create the first one...
			logger.Tracef(nil, "++ Resource version %s  ID %s\n", version, r.ID)
			if vres, ok = c.ResourceList[version][r.ID]; !ok {
				logger.Tracef(nil, "   - Creating new resource\n")
				vres = r
			}
			c.ResourceList[version][r.ID] = vres

			// Compile a list of the methods which use this resource
			vres.Methods = append(vres.Methods, *method)

			// Add the resource to the method which uses it
			method.Resources = append(method.Resources, vres)

		}

		method.Responses[status] = Response{
			Description: string(github_flavored_markdown.Markdown([]byte(response.Description))),
			Schema:      vres,
		}
	}

	if o.Responses.Default != nil {
		r := c.resourceFromSchema(o.Responses.Default.Schema, method, nil)
		if r != nil {
			logger.Tracef(nil, "++ Resource version %s  ID %s\n", version, r.ID)
			// Look for a pre-declared resource with the response ID, and use that or create the first one...
			var vres *Resource
			var ok bool
			if vres, ok = c.ResourceList[version][r.ID]; !ok {
				logger.Tracef(nil, "   - Creating new resource\n")
				vres = r
			}
			c.ResourceList[version][r.ID] = vres

			// Add to the compiled list of methods which use this resource
			vres.Methods = append(vres.Methods, *method)

			// Set the default response
			method.DefaultResponse = &Response{
				Description: string(github_flavored_markdown.Markdown([]byte(o.Responses.Default.Description))),
				Schema:      vres,
			}
		}
	}

	// TODO FIXME If no Security given from operation, then the global defaults are appled. CHECK THIS IS TRUE!

	method.Security = make(map[string]Security)

	for _, sec := range o.Security {
		for n, scopes := range sec {
			// Lookup security name in definitions
			if scheme, ok := c.SecurityDefinitions[n]; ok {

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

func (c *APISpecification) resourceFromSchema(s *spec.Schema, method *Method, fqNS []string) *Resource {
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

	//logger.Printf(nil, "CHECK schema type and items\n")
	//spew.Dump(s)

	if s.Type == nil {
		s.Type = append(s.Type, "object")
	}

	if s.Items != nil {
		stringorarray := s.Type

		// EEK This is officially icky! See the Activities model in the uber spec. It declares "items": [ { } ] !!
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

	id := TitleToKebab(s.Title)

	if len(fqNS) == 0 && id == "" {
		logger.Errorf(nil, "Error: %s %s references a model definition that does not have a title memeber.", strings.ToUpper(method.Method), method.Path)
		spew.Dump(s)
		os.Exit(1)
	}

	myFQNS := append([]string{}, fqNS...)
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
			logger.Errorf(nil, "Error encoding example json: %s", err)
		}
		r.Example = string(example)
	}

	if len(s.Enum) > 0 {
		for _, e := range s.Enum {
			r.Enum = append(r.Enum, fmt.Sprintf("%s", e))
		}
	}

	//logger.Tracef(nil, "expandSchema Type %s FQNS '%s'\n", s.Type, strings.Join(myFQNS, "."))

	required := make(map[string]bool)
	json_representation := make(map[string]interface{})

	c.compileproperties(s, r, method, id, required, json_representation, myFQNS, chopped)

	for allof := range s.AllOf {
		c.compileproperties(&s.AllOf[allof], r, method, id, required, json_representation, myFQNS, chopped)
	}

	// Build element of resource schema example
	// FIXME This *explodes* if there is no "type" member in the actual model definition - which is probably right
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
			logger.Errorf(nil, "Error encoding schema json: %s", err)
		}
		r.Schema = string(schema)
	}

	return r
}

// -----------------------------------------------------------------------------
// Takes a Schema object and adds properties to the Resource object.
// It uses the 'required' map to set when properties are required and builds a JSON
// representation of the resource.
//
func (c *APISpecification) compileproperties(s *spec.Schema, r *Resource, method *Method, id string, required map[string]bool, json_rep map[string]interface{}, myFQNS []string, chopped bool) {

	// First, grab the required members
	for _, i := range s.Required {
		required[i] = true
	}

	// Now process the properties
	for name, property := range s.Properties {
		//log.Printf("Process property name '%s'  Type %s\n", name, s.Properties[name].Type)
		newFQNS := append([]string{}, myFQNS...)
		if chopped && len(id) > 0 {
			newFQNS = append(newFQNS, id)
		}
		newFQNS = append(newFQNS, name)

		r.Properties[name] = c.resourceFromSchema(&property, method, newFQNS)

		if _, ok := required[name]; ok {
			r.Properties[name].Required = true
		}

		// XXX This really is quite a juggle!
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

						r.Properties[name] = c.resourceFromSchema(property.Items.Schema, method, xFQNS)

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
						json_rep[name] = f

						// Override type to reflect it is an array
						r.Properties[name].Type[0] = "array[" + r.Properties[name].Type[0] + "]"
					}
				}
			} else {
				json_rep[name] = r.Properties[name].Schema
			}
		} else {
			var f interface{}
			_ = json.Unmarshal([]byte(r.Properties[name].Schema), &f)
			json_rep[name] = f
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
