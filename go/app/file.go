package main

import (
	"encoding/json"
	"io/ioutil"
)

func readItemListFromFile() (*Items, error) {
	filename := "./app/items.json"
	//read file
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	//update item
	var items Items
	if len(buf) != 0 {
		err = json.Unmarshal(buf, &items)
		if err != nil {
			return nil, err
		}
	}
	return &items, nil
}

func updateFile(name string, category string, imgHash string) error {
	filename := "./app/items.json"
	items, err := readItemListFromFile()
	if err != nil {
		return err
	}

	//items.Items = append(items.Items, Item{Name: name, Category: category, ImageFilename: imgHash})

	items.Items = append(items.Items, Item{Name: name, Category: category, ImageFilename: imgHash})

	buf, err := json.Marshal(items)
	if err != nil {
		return err
	}

	//write file
	err = ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		return err
	}
	return nil
}
