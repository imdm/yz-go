package yz_go

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
)

type ICache interface {
	ExpiresIn() int64
}

type Cache interface {
	Set(data ICache) error
	Get(data ICache) error
}

type FileCache struct {
	Path string
}

func newFileCache(path string) *FileCache {
	return &FileCache{
		Path: path,
	}
}

func (f *FileCache) Set(data ICache) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.Path, bytes, 0644)
}

func (f *FileCache) Get(data ICache) error {
	bytes, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}
	expires := data.ExpiresIn()
	if time.Now().Unix() > expires-60 {
		return errors.New("data is expired")
	}
	return nil
}
