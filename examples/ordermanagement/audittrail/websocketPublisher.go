package audittrail

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/project-flogo/core/support/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// AuditTrailChannel to receive audit trail item
var AuditTrailChannel chan Entry
var streamName string

// WebSocketPublisher to publish data to web client
type WebSocketPublisher interface {
	Start()
}

type webSocketPublisherImpl struct {
	port int64
}

// Create a new WebSocketPublisher Instance
func Create(port int64, stream string) WebSocketPublisher {
	streamName = stream
	return &webSocketPublisherImpl{port: port}
}

func (wsp *webSocketPublisherImpl) Start() {
	AuditTrailChannel = make(chan Entry)

	http.HandleFunc("/auditTrail", processAuditTrail)

	log.RootLogger().Infof("Started Websocket server on port [%s]", strconv.FormatInt(wsp.port, 10))
	http.ListenAndServe(":"+strconv.FormatInt(wsp.port, 10), nil)
}

// Fetch records from Kinesis and push them onto the web client
func processAuditTrail(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	// start loading stream records
	go KPInstance.GetRecords(streamName)

	for true {
		incoming := <-AuditTrailChannel

		formattedData := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", incoming.TimeStamp.Format("2006-01-02 15:04:05"), incoming.OrderID, incoming.RuleName, incoming.Status, incoming.Description)

		// Write message back to browser
		if err = conn.WriteMessage(websocket.TextMessage, []byte(formattedData)); err != nil {
			panic(err)
		}
	}
}
