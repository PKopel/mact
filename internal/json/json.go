package internal

import (
	"log"
	"strings"

	"github.com/PKopel/mact/internal/utils"
	"github.com/PKopel/mact/types"
)

type JSON map[string]interface{}

// FindField returns the map object for the given field path or the last map found on that path and the unmatched path
func FindField(object JSON, fieldPath []string) (JSON, []string) {
	for i, field := range fieldPath {
		val, ok := object[field]
		if !ok {
			return object, fieldPath[i:]
		}
		newObj, ok := val.(*JSON)
		if !ok {
			return object, fieldPath[i:]
		}
		object = *newObj
	}
	return object, nil
}

func ApplyChanges(body JSON, changes []types.Change) JSON {

	for _, change := range changes {
		fullPath := strings.Split(change.Field, ".")
		parents, field := utils.Unsnoc(fullPath)
		object, path := FindField(body, parents)
		if path != nil {
			log.Printf("couldn't match path: %v to response JSON", fullPath)
			continue
		}
		switch change.Type {
		case types.Add, types.Modify:
			object[*field] = change.Value
		case types.Remove:
			delete(object, *field)
		}
	}

	return body
}
