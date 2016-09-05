package spec

import (
	"github.com/mitchellh/hashstructure"
	"github.com/zxchris/swaggerly/logger"
)

func (i *APISuiteMap) Merge(s *APISpecification) error {
	// 1.     Calculate hash of new specification
	// 1.1.   Add hash to specification
	// 1.2.   If specification already exists in suite and version is greater
	//   2.2.   If hashes do not agree, replace existing with new
	// 1.3    else
	//   3.1.   Store specification in suite

	logger.Tracef(nil, "Hashing specification %s", s.ID)
	hash, err := hashstructure.Hash(s, nil)
	if err != nil {
		logger.Tracef(nil, "Hash error %s", err)
		return err
	}

	s.Hash = hash
	logger.Tracef(nil, "Hashed spec %s as %d", s.ID, s.Hash)

	if existing, ok := (*i)[s.ID]; ok {
		// TODO check version numbers
		if existing.Hash != s.Hash {
			(*i)[s.ID] = s
			logger.Tracef(nil, "Replace existing specification")
		}
	} else {
		(*i)[s.ID] = s
		logger.Tracef(nil, "Merged %s into APISuite", s.ID)
	}

	return nil
}

// For the given specification:
//   For each API (endpoint)
//     For each method
//       See if the API/Method exists and it's version
//       If the versions are differ

func (i *APISuiteMap) FoldVersions(s *APISpecification) error {
	for _, endpoint := range s.APIs {
		for _, m := range endpoint.Methods {
			logger.Tracef(nil, "API name: %s Method Name: %s  Method: %s  Path: %s", endpoint.Name, m.Name, m.Method, m.Path)
		}
	}
	return nil
}


// end
