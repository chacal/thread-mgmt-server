package mqtt

//go:generate mockgen -destination=../mocks/mock_mqtt_sender.go -package=mocks github.com/chacal/thread-mgmt-server/pkg/mqtt MqttSender

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
)

const MQTT_STATE_TAG = "d"

type MqttSender interface {
	Connect() chan bool
	PublishState(state device_registry.State)
}

type mqttSender struct {
	client mqtt.Client
}

func CreateSender(connectionUrl string, username string, password string) *mqttSender {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(connectionUrl)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID("thread-mgmt-srv-" + strconv.Itoa(randBetween(1000000, 9999999)))
	opts.SetOnConnectHandler(onMqttConnected)
	opts.SetConnectionLostHandler(onMqttConnectionLost)
	opts.SetConnectRetry(true)
	opts.SetReconnectingHandler(onMqttReconnecting)
	opts.SetTLSConfig(&tls.Config{})
	return &mqttSender{mqtt.NewClient(opts)}
}

func (s *mqttSender) Connect() chan bool {
	ret := make(chan bool)
	token := s.client.Connect()
	go func() {
		_ = token.Wait()
		if token.Error() != nil {
			log.Errorf("connection failed: %v", errors.WithStack(token.Error()))
		}
		ret <- true
	}()
	return ret
}

func (s *mqttSender) PublishState(state device_registry.State) {
	if !s.client.IsConnectionOpen() {
		log.Warnf("Can't publish state for device %v. MQTT not connected.", state.Instance)
	} else {
		buf, err := json.Marshal(state)
		if err != nil {
			log.Error(publishErrorMsg(state, err))
		}

		topic := fmt.Sprintf("/sensor/%s/%s/state", state.Instance, MQTT_STATE_TAG)
		t := s.client.Publish(topic, 1, true, buf)
		go func() {
			_ = t.Wait()
			if t.Error() != nil {
				log.Error(publishErrorMsg(state, err))
			}
		}()
	}
}

func randBetween(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func onMqttConnected(client mqtt.Client) {
	r := client.OptionsReader()
	log.Infof("MQTT connected to %v", r.Servers())
}

func onMqttConnectionLost(client mqtt.Client, err error) {
	log.Infof("MQTT connection lost. Error: %v", err)
}

func onMqttReconnecting(client mqtt.Client, opts *mqtt.ClientOptions) {
	log.Infof("MQTT reconnecting to %v", opts.Servers)
}

func publishErrorMsg(state device_registry.State, err error) string {
	return fmt.Sprintf("Failed to publish state for device %s. Error: %v", state.Instance, err)
}
