package base

import (
	"fmt"
	"strings"
	"time"
)

const orderIDMaxLen = 30

func GenOrderID(uid string) string {
	now := time.Now().Unix()
	orderID := fmt.Sprintf("%d%s", now, uid)
	orderID = strings.TrimSpace(orderID)
	n := len(orderID)
	if n > orderIDMaxLen {
		orderID = orderID[n-orderIDMaxLen:]
	}
	return orderID
}
