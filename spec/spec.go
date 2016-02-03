package spec

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/serenize/snaker"
	"github.com/shurcooL/github_flavored_markdown"
)

// APISet is a slice of API structs
type APISet []API

// APIs represents the parsed APIs
var APIs APISet

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
type API struct {
	ID      string
	Name    string
	Methods []Method
	//Resources []*Resource
	URL *url.URL
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
	API             API
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
	Description    string
	Schema         *Resource
	Versions       map[string]*Resource
	DefaultVersion string
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
	spec, err := loadSpec("http://" + host + "/spec/swagger.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", spec)

	u, err := url.Parse(spec.Spec().Schemes[0] + "://" + spec.Spec().Host)
	if err != nil {
		log.Fatal(err)
	}

	for _, tag := range spec.Spec().Tags {
		api := API{
			ID:   titleToKebab(tag.Name),
			Name: tag.Name,
			URL:  u,
		}

		for p, o := range spec.AllPaths() {
			getMethod(tag, &api, o.Get, p, "get")
			getMethod(tag, &api, o.Post, p, "post")
			getMethod(tag, &api, o.Put, p, "put")
			getMethod(tag, &api, o.Delete, p, "delete")
			getMethod(tag, &api, o.Head, p, "head")
			getMethod(tag, &api, o.Options, p, "options")
			getMethod(tag, &api, o.Patch, p, "patch")
		}

		APIs = append(APIs, api)
	}
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

func getMethod(tag spec.Tag, api *API, o *spec.Operation, path, method string) {
	if o != nil {
		for _, t := range o.Tags {
			if t == tag.Name {
				method := &Method{
					ID:          camelToKebab(o.ID),
					Name:        o.Summary,
					Description: o.Description,
					Method:      method,
					Path:        path,
					Responses:   make(map[int]Response),
					API:         *api,
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

					alts := make(map[string]*Resource)
					if versions, ok := o.Extensions["x-versions"]; ok {
						if versions, ok := versions.(map[string]interface{}); ok {
							for version, obj := range versions {
								if obj, ok := obj.(map[string]interface{}); ok {
									if ref, ok := obj["$ref"]; ok {
										if url, ok := ref.(string); ok {
											ref := spec.RefProperty(url)
											err := spec.ExpandSchema(ref, nil, nil)
											if err != nil {
												// FIXME?
												log.Fatalf("error loading inner spec %s: %s", url, err)
											}
											alts[version] = resourceFromSchema(ref, []string{})
										}
									}
								}
							}
						}
					}

					var defaultVersion string
					if dV, ok := o.Extensions["x-default-version"]; ok {
						if sdV, ok := dV.(string); ok {
							defaultVersion = sdV
						}
					}
					if len(defaultVersion) > 0 {
						if _, ok := alts[defaultVersion]; !ok {
							// FIXME?
							log.Fatalf("default version not found in versions map")
						}
					} else {
						for k := range alts {
							defaultVersion = k
							break
						}
					}

					method.Responses[status] = Response{
						Description:    response.Description,
						Schema:         r,
						Versions:       alts,
						DefaultVersion: defaultVersion,
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
					//api.Resources = append(api.Resources, r)
				}

				api.Methods = append(api.Methods, *method)
			}
		}
	}
}

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

	example, err := json.MarshalIndent(&s.Example, "", "    ")
	if err != nil {
		log.Printf("error encoding example json: %s", err)
	}

	r := &Resource{
		ID:          id,
		Title:       s.Title,
		Description: string(github_flavored_markdown.Markdown([]byte(s.Description))),
		Example:     string(example),
		Type:        s.Type,
		Properties:  make(map[string]*Resource),
		FQNS:        myFQNS,
	}

	if len(s.Enum) > 0 {
		for _, e := range s.Enum {
			r.Enum = append(r.Enum, fmt.Sprintf("%s", e))
		}
	}

	required := make(map[string]bool)
	for _, r := range s.Required {
		required[r] = true
	}

	props := make(map[string]interface{})

	for name, property := range s.Properties {
		newFQNS := append([]string{}, myFQNS...)
		if chopped && len(id) > 0 {
			newFQNS = append(newFQNS, id)
		}
		newFQNS = append(newFQNS, name)
		r.Properties[name] = resourceFromSchema(&property, newFQNS)
		if _, ok := required[name]; ok {
			r.Properties[name].Required = true
		}

		// FIXME this is as nasty as it looks...
		if strings.ToLower(r.Properties[name].Type[0]) != "object" {
			props[name] = r.Properties[name].Schema
		} else {
			var f interface{}
			_ = json.Unmarshal([]byte(r.Properties[name].Schema), &f)
			props[name] = f
		}
	}

	// FIXME also as nasty as it looks
	if strings.ToLower(r.Type[0]) != "object" {
		r.Schema = r.Type[0]
	} else {
		schema, err := json.MarshalIndent(props, "", "    ")
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
