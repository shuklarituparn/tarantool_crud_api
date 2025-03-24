package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shuklarituparn/tarantool_crud_api/internal/prometheus"
	"github.com/shuklarituparn/tarantool_crud_api/internal/repository"
	"github.com/shuklarituparn/tarantool_crud_api/internal/utils"
	"github.com/shuklarituparn/tarantool_crud_api/pkg/logger"

	"github.com/gorilla/mux"
)

// swagger:model
type KVRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value" swaggertype:"object"`
}

// swagger:model
type UpdateKVRequest struct {
	Value interface{} `json:"value" swaggertype:"object"`
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode JSON response: ", err)
	}
}

// @Summary Создает пару ключ-значение
// @Description Сохраняет пару ключ-значение в базе данных
// @Tags Key-Value Store
// @Accept json
// @Produce json
// @Param request body KVRequest true "Key-Value Data"
// @Success 201
// @Failure 400
// @Failure 409
// @Router /api/v1/kv [post]
func CreateKV(w http.ResponseWriter, r *http.Request) {
	prometheus.CreateKeyValuePairCounter.Inc()
	var req KVRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body: ", err)
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
		return
	}

	if req.Key == "" {
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "Key cannot be empty"})
		return
	}
	result, err := repository.CreateKV(req.Key, req.Value)
	if err != nil {
		logger.Log.Error("Failed to create KV pair: ", err)
		if err == repository.ErrKeyExists {
			sendJSONResponse(w, http.StatusConflict, map[string]interface{}{"error": err.Error()})
		} else {
			sendJSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": "Internal server error"})
		}
		return
	}

	logger.Log.Info("KV pair created successfully: ", req.Key)
	sendJSONResponse(w, http.StatusCreated, map[string]interface{}{"value": result, "key": req.Key})
}

// @Summary Получение пару ключ-значение
// @Description Получает значение, связанное с ключом
// @Tags Key-Value Store
// @Produce json
// @Param key path string true "Key"
// @Success 200
// @Failure 404
// @Router /api/v1/kv/{key} [get]
func GetKV(w http.ResponseWriter, r *http.Request) {
	prometheus.ReadKeyValuePairCounter.Inc()

	key := mux.Vars(r)["key"]
	if key == "" {
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "Key cannot be empty"})
		return
	}
	result, err := repository.GetKV(key)
	if err != nil {
		logger.Log.Error("Key not found: ", key)
		if err == repository.ErrKeyNotFound {
			sendJSONResponse(w, http.StatusNotFound, map[string]interface{}{"error": err.Error()})
		} else {
			sendJSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": "Internal server error"})
		}
		return
	}
	safeResult := utils.ConvertToJSONSafe(result)
	response := map[string]interface{}{
		"value": safeResult,
	}

	logger.Log.Info("KV pair retrieved successfully: ", key)
	sendJSONResponse(w, http.StatusOK, response)
}

// @Summary Обновление пару ключ-значение
// @Description Обновляет значение, связанное с ключом
// @Tags Key-Value Store
// @Accept json
// @Produce json
// @Param key path string true "Key"
// @Param request body UpdateKVRequest true "Updated value data"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /api/v1/kv/{key} [put]
func UpdateKV(w http.ResponseWriter, r *http.Request) {
	prometheus.UpdateKeyValuePairCounter.Inc()

	var req UpdateKVRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Warn("Invalid request body: ", err)
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
		return
	}

	key := mux.Vars(r)["key"]
	if key == "" {
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "Key cannot be empty"})
		return
	}
	result, err := repository.UpdateKV(key, req.Value)
	if err != nil {
		logger.Log.Warn("Error updating key: ", key, err)
		if err == repository.ErrKeyNotFound {
			sendJSONResponse(w, http.StatusNotFound, map[string]interface{}{"error": err.Error()})
		} else {
			sendJSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": "Internal server error"})
		}
		return
	}
	logger.Log.Info("Key updated successfully: ", key)
	sendJSONResponse(w, http.StatusOK, map[string]interface{}{"value": result})
}

// @Summary Удаление пары ключ-значение
// @Description Удаляет пару ключ-значение из базы данных
// @Tags Key-Value Store
// @Produce json
// @Param key path string true "Key"
// @Success 200
// @Failure 404
// @Router /api/v1/kv/{key} [delete]
func DeleteKV(w http.ResponseWriter, r *http.Request) {
	prometheus.DeleteKeyValuePairCounter.Inc()
	key := mux.Vars(r)["key"]
	if key == "" {
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "Key cannot be empty"})
		return
	}
	err := repository.DeleteKV(key)
	if err != nil {
		logger.Log.Warn("Error deleting key: ", key, err)
		if err == repository.ErrKeyNotFound {
			sendJSONResponse(w, http.StatusNotFound, map[string]interface{}{"error": err.Error()})
		} else {
			sendJSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": "Internal server error"})
		}
		return
	}

	logger.Log.Info("Key deleted successfully: ", key)
	sendJSONResponse(w, http.StatusOK, map[string]interface{}{"message": "Key deleted successfully"})
}
