package main

import (
	"io/ioutil"
	"firebase.google.com/go/v4/storage"
)

func getWorldList(client *storage.Client) ([]string, error) {

	bucket, err := client.DefaultBucket()

	if err != nil {
		return nil, err
	}

	iter := bucket.Objects(ctx, nil)

	var names []string

	for {
		attrs, err := iter.Next()
		if err != nil {
			break
		}

		if attrs == nil {
			break
		}

		names = append(names, attrs.Name)
	}

	return names, nil

}

func downloadWorld(client *storage.Client, name string) ([]byte, error) {

	bucket, err := client.DefaultBucket()
	if err != nil {
		return nil, err
	}

	obj := bucket.Object(name)

	file, err := obj.NewReader(ctx)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return data, nil

}
