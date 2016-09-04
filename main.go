package main

import (
	"os/exec"
	"encoding/json"
	"fmt"
	"errors"
)

type Image interface {
	Tags() map[string]interface{}
	AddTag(name string, value string) error
	AddTagValue(name string, value string) error
}

type ImageCache struct {
	filepath string
	tags map[string]interface{}
}

func NewImage(filepath string) (Image, error) {

	cmdOut, err := callTool("-j", filepath)

	if err != nil {
		return nil, err
	}

	var tags []map[string]interface{}

	if err := json.Unmarshal([]byte(cmdOut), &tags); err != nil {
		return nil, err
	}

	return &ImageCache{filepath:filepath,tags:tags[0]}, nil
}

func (img ImageCache) Tags() map[string]interface{} {
	return img.tags
}

func (img ImageCache) AddTag(name string, value string) error {

	if name == "" {
		return errors.New("name required")
	}
	if value == "" {
		return errors.New("value required")
	}

	if img.tags[name] != nil {
		return errors.New(fmt.Sprintf("Tag %v already exists", name))
	}

	out, err := callTool(fmt.Sprintf("-%v=%v", name, value), img.filepath)

	if err != nil {
		return errors.New(fmt.Sprintf("%v: %v", err, out))
	}

	img.tags[name] = value

	return nil
}

func (img ImageCache) RemoveTag(name string) error {

	if name == "" {
		return errors.New("name required")
	}

	if img.tags[name] == nil {
		return errors.New(fmt.Sprintf("Tag %v does not exist", name))
	}

	out, err := callTool(fmt.Sprintf("-%v=", name), img.filepath)

	if err != nil {
		return errors.New(fmt.Sprintf("%v: %v", err, out))
	}

	delete(img.tags, name)

	return nil
}

func (img ImageCache) AddTagValue(name string, value string) error {

	if name == "" {
		return errors.New("name required")
	}
	if value == "" {
		return errors.New("value required")
	}

	current := img.tags[name]

	var vals []string

	if(current == nil) {
		vals = make([]string, 0)
	} else {

		switch v := current.(type) {
		default:
			return errors.New(fmt.Sprintf("unexpected tag type %T", v))
		case string:
			vals = []string{current.(string)}
		case []string:
			vals = current.([]string)
		}
	}

	out, err := callTool(fmt.Sprintf("-%v+=%v", name, value), img.filepath)

	if err != nil {
		return errors.New(fmt.Sprintf("%v: %v", err, out))
	}

	vals = append(vals, value)
	img.tags[name] = vals

	return nil
}

func (img ImageCache) RemoveTagValue(name string, value string) error {

	if name == "" {
		return errors.New("name required")
	}
	if value == "" {
		return errors.New("value required")
	}

	current := img.tags[name]

	var vals []string

	if(current == nil) {
		return errors.New(fmt.Sprintf("Tag not found: %v", name))
	} else {

		switch v := current.(type) {
		default:
			return errors.New(fmt.Sprintf("unexpected tag type %T", v))
		case string:
			vals = []string{current.(string)}
		case []string:
			vals = current.([]string)
		}
	}

	out, err := callTool(fmt.Sprintf("-%v+=%v", name, value), img.filepath)

	if err != nil {
		return errors.New(fmt.Sprintf("%v: %v", err, out))
	}

	vals = append(vals, value)
	img.tags[name] = vals

	return nil
}


func callTool(args ...string) (string, error) {

	cmdName := "exiftool"
	cmdOut, err := exec.Command(cmdName, args...).CombinedOutput()

	return string(cmdOut), err
}