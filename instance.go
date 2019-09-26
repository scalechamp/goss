package goss

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (api *API) waitUntilReady(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		response, err := api.sling.Path("/v1/instances/").Get(id).Receive(&data, &failed)
		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintf("waitUntilReady failed, status: %v, message: %s", response.StatusCode, failed))
		}
		if data["status"] == 2 {
			data["id"] = id
			return data, nil
		}

		time.Sleep(10 * time.Second)
	}
}

func (api *API) CreateInstance(params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	response, err := api.sling.Post("/v1/instances").BodyJSON(params).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("CreateInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	data["id"] = strconv.FormatFloat(data["id"].(float64), 'f', 0, 64)
	return api.waitUntilReady(data["id"].(string))
}

func (api *API) ReadInstance(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	response, err := api.sling.Path("/v1/instances/").Get(id).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) ReadInstances() ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	response, err := api.sling.Get("/v1/instances").Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadInstances failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) DeleteInstance(id string) error {
	failed := make(map[string]interface{})
	response, err := api.sling.Path("/v1/instances/").Delete(id).Receive(nil, &failed)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("DeleteInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}
	return nil
}