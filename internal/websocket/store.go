package websocket

import (
	"errors"
	"maps"
	"sync"
)

type WebsocketConnectionStore struct {
	connections map[string]*WebsocketConnection
	mutex       *sync.RWMutex
}

type WebsocketTopicStore struct {
	topics map[string]bool
	mutex  *sync.RWMutex
}

func NewConnectionStore() *WebsocketConnectionStore {
	return &WebsocketConnectionStore{
		connections: make(map[string]*WebsocketConnection),
		mutex:       &sync.RWMutex{},
	}
}

func NewTopicStore() *WebsocketTopicStore {
	return &WebsocketTopicStore{
		topics: make(map[string]bool),
		mutex:  &sync.RWMutex{},
	}
}

func (connectionStore *WebsocketConnectionStore) Get(ID string) (*WebsocketConnection, error) {
	connectionStore.mutex.RLock()
	defer connectionStore.mutex.RUnlock()

	connection, ok := connectionStore.connections[ID]
	if !ok {
		return nil, errors.New("connection not found")
	}

	return connection, nil
}

func (connectionStore *WebsocketConnectionStore) GetByTopic(topic string) []*WebsocketConnection {
	connectionStore.mutex.RLock()
	defer connectionStore.mutex.RUnlock()

	connections := make([]*WebsocketConnection, 0)

	for _, connection := range connectionStore.connections {
		if connection.Topics.Has(topic) {
			connections = append(connections, connection)
		}
	}

	return connections
}

func (connectionStore *WebsocketConnectionStore) GetAll() []*WebsocketConnection {
	connectionStore.mutex.RLock()
	defer connectionStore.mutex.RUnlock()

	connections := make([]*WebsocketConnection, 0, len(connectionStore.connections))

	for _, connection := range connectionStore.connections {
		connections = append(connections, connection)
	}

	return connections
}

func (connectionStore *WebsocketConnectionStore) Set(connection *WebsocketConnection) {
	connectionStore.mutex.Lock()
	defer connectionStore.mutex.Unlock()

	connectionStore.connections[connection.ID.String()] = connection
}

func (connectionStore *WebsocketConnectionStore) Remove(ID string) {
	connectionStore.mutex.Lock()
	defer connectionStore.mutex.Unlock()
	delete(connectionStore.connections, ID)
}

func (connectionStore *WebsocketConnectionStore) Clear() {
	connectionStore.mutex.Lock()
	defer connectionStore.mutex.Unlock()

	keys := maps.Keys(connectionStore.connections)

	for key := range keys {
		delete(connectionStore.connections, key)
	}
}

func (topicStore *WebsocketTopicStore) Get() []string {
	topicStore.mutex.RLock()
	defer topicStore.mutex.RUnlock()

	keys := maps.Keys(topicStore.topics)
	topics := make([]string, 0, len(topicStore.topics))

	for key := range keys {
		topics = append(topics, key)
	}

	return topics
}

func (topicStore *WebsocketTopicStore) Has(topic string) bool {
	topicStore.mutex.RLock()
	defer topicStore.mutex.RUnlock()

	return topicStore.topics[topic]
}

func (topicStore *WebsocketTopicStore) Set(topic string) {
	topicStore.mutex.Lock()
	defer topicStore.mutex.Unlock()

	topicStore.topics[topic] = true
}

func (topicStore *WebsocketTopicStore) Remove(topic string) {
	topicStore.mutex.Lock()
	defer topicStore.mutex.Unlock()
	delete(topicStore.topics, topic)
}

func (topicStore *WebsocketTopicStore) Clear() {
	topicStore.mutex.Lock()
	defer topicStore.mutex.Unlock()

	keys := maps.Keys(topicStore.topics)

	for key := range keys {
		delete(topicStore.topics, key)
	}
}
