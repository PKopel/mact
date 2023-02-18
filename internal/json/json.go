package internal

import (
	"log"
	"strings"

	conf "github.com/PKopel/mact/internal/config"
	"github.com/PKopel/mact/internal/utils"
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

func ApplyChanges(body JSON, changes []conf.Change) JSON {

	for _, change := range changes {
		fullPath := strings.Split(change.Field, ".")
		parents, field := utils.Unsnoc(fullPath)
		object, path := FindField(body, parents)
		if path != nil {
			log.Printf("couldn't match path: %v to response JSON", fullPath)
			continue
		}
		switch change.Type {
		case conf.Add, conf.Modify:
			object[*field] = change.Value
		case conf.Remove:
			delete(object, *field)
		}
	}

	return body
}
