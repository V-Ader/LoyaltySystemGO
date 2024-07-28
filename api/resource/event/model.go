package event

import (
	"fmt"

	"time"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
)

type Event struct {
	Id        int       `json:"id"`
	Card_id   int       `json:"card_id"`
	Timestamp time.Time `json:"timestamp"`
	Quantity  int       `json:"quantity"`
}

func (event *Event) GetHash() string {
	return common.GenerateETag([]byte(fmt.Sprintf("%v", event)))
}

type EventDataRequest struct {
	Card_id  int `json:"card_id"`
	Quantity int `json:"quantity"`
}
