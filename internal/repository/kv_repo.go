package repository

import (
	"errors"
	"fmt"
	"github.com/shuklarituparn/tarantool_crud_api/config"
	"github.com/shuklarituparn/tarantool_crud_api/pkg/logger"

	tarantool "github.com/tarantool/go-tarantool"
)

var (
	ErrKeyExists   = errors.New("key already exists")
	ErrKeyNotFound = errors.New("key not found")
)

func extractValueFromTarantoolResponse(res []interface{}) (interface{}, error) {
	if len(res) == 0 {
		return nil, ErrKeyNotFound
	}

	tuple, ok := res[0].([]interface{})
	if !ok || len(tuple) < 2 {
		return nil, errors.New("invalid response format from database")
	}

	return tuple[1], nil
}

func CreateKV(key string, value interface{}) (interface{}, error) {
	logger.Log.Info(fmt.Sprintf("Attempting to create key: %s with value: %v", key, value))

	resp, err := config.Conn.Insert("kv_store", []interface{}{key, value})
	if err != nil {
		if tarantoolErr, ok := err.(tarantool.Error); ok && tarantoolErr.Code == tarantool.ErrTupleFound {
			logger.Log.Error(fmt.Sprintf("Key %s already exists", key))
			return nil, ErrKeyExists
		}
		logger.Log.Error(fmt.Sprintf("Database error while inserting key %s: %v", key, err))
		return nil, err
	}

	if len(resp.Data) > 0 {
		val, err := extractValueFromTarantoolResponse(resp.Data)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Error extracting value for key %s: %v", key, err))
			return nil, err
		}
		logger.Log.Info(fmt.Sprintf("Successfully created key: %s", key))
		return val, nil
	}

	logger.Log.Info(fmt.Sprintf("Successfully created key: %s", key))
	return value, nil
}

func GetKV(key string) (interface{}, error) {
	logger.Log.Info(fmt.Sprintf("Fetching key: %s", key))

	var res []interface{}
	err := config.Conn.SelectTyped("kv_store", "primary", 0, 1, tarantool.IterEq, []interface{}{key}, &res)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Database error while fetching key %s: %v", key, err))
		return nil, err
	}

	val, err := extractValueFromTarantoolResponse(res)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Key %s not found or invalid format", key))
		return nil, err
	}

	logger.Log.Info(fmt.Sprintf("Successfully retrieved key: %s", key))
	return val, nil
}

func UpdateKV(key string, value interface{}) (interface{}, error) {
	logger.Log.Info(fmt.Sprintf("Updating key: %s with new value: %v", key, value))

	_, err := GetKV(key)
	if err != nil {
		if err == ErrKeyNotFound {
			logger.Log.Error(fmt.Sprintf("Cannot update non-existent key %s", key))
			return nil, ErrKeyNotFound
		}
		logger.Log.Error(fmt.Sprintf("Error checking existence of key %s: %v", key, err))
		return nil, err
	}

	resp, err := config.Conn.Replace("kv_store", []interface{}{key, value})
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Database error while updating key %s: %v", key, err))
		return nil, err
	}

	if len(resp.Data) > 0 {
		val, err := extractValueFromTarantoolResponse(resp.Data)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Error extracting updated value for key %s: %v", key, err))
			return nil, err
		}
		logger.Log.Info(fmt.Sprintf("Successfully updated key: %s", key))
		return val, nil
	}

	logger.Log.Info(fmt.Sprintf("Successfully updated key: %s", key))
	return value, nil
}

func DeleteKV(key string) error {
	logger.Log.Info(fmt.Sprintf("Deleting key: %s", key))

	_, err := GetKV(key)
	if err != nil {
		return err
	}

	_, err = config.Conn.Delete("kv_store", "primary", []interface{}{key})
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Database error while deleting key %s: %v", key, err))
		return err
	}

	logger.Log.Info(fmt.Sprintf("Successfully deleted key: %s", key))
	return nil
}
