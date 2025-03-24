package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	CreateKeyValuePairCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "create_key_value_pair_counter",
		Help: "Количество запросов, выполненных к конечной точке Key-Value Store для создания пары ключ-значение",
	})
	UpdateKeyValuePairCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "update_value_from_key_counter",
		Help: "Количество запросов, выполненных к конечной точке Key-Value Store для обновления пары ключ-значение",
	})
	ReadKeyValuePairCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "get_value_from_key_counter",
		Help: "Количество запросов, выполненных к конечной точке Key-Value Store для получения значения по ключу",
	})
	DeleteKeyValuePairCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "delete_value_from_key_counter",
		Help: "Количество запросов, выполненных к конечной точке Key-Value Store для удаления значения по ключу",
	})
)

func Init() {
	prometheus.MustRegister(CreateKeyValuePairCounter)
	prometheus.MustRegister(UpdateKeyValuePairCounter)
	prometheus.MustRegister(ReadKeyValuePairCounter)
	prometheus.MustRegister(DeleteKeyValuePairCounter)

}
