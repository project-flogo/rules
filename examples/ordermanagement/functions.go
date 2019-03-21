package main

import (
	"context"
	"fmt"
	"time"

	"github.com/project-flogo/rules/examples/ordermanagement/audittrail"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/project-flogo/core/support/log"
)

func init() {
	//rule orderevent, handle a order event, create a order in the orderAction
	config.RegisterConditionEvaluator("truecondition", truecondition)
	config.RegisterActionFunction("ordereventAction", ordereventAction)

	//rule itemevent, handle a item event, update order in the itemeventAction
	config.RegisterConditionEvaluator("itemeventCondition", itemeventCondition)
	config.RegisterActionFunction("itemeventAction", itemeventAction)

	//rule order
	config.RegisterConditionEvaluator("orderCondition", orderCondition)
	config.RegisterActionFunction("orderAction", orderAction)

	//rule ordertimeoutevent, order timeout if within a specified interval, all order items dont arrive
	config.RegisterConditionEvaluator("ordertimeouteventCondition", ordertimeouteventCondition)
	config.RegisterActionFunction("ordertimeouteventAction", ordertimeouteventAction)

	//rule ordershippedevent
	config.RegisterConditionEvaluator("ordershippedCondition", ordershippedCondition)
	config.RegisterActionFunction("ordershippedAction", ordershippedAction)

	//rule ordercancelled, notification on order cancelled
	config.RegisterConditionEvaluator("ordercancelledCondition", ordercancelledCondition)
	config.RegisterActionFunction("ordercancelledAction", ordercancelledAction)

	//rule shippingtimeoutevent, shipping timeout if within a specified interval, order is not shipped
	config.RegisterConditionEvaluator("shippingtimeouteventCondition", shippingtimeouteventCondition)
	config.RegisterActionFunction("shippingtimeouteventAction", shippingtimeouteventAction)

	//rule orderinvoice, apply discounts based on the customer level on hte final invoice
	config.RegisterConditionEvaluator("orderinvoiceCondition", orderinvoiceCondition)
	config.RegisterActionFunction("orderinvoiceAction", orderinvoiceAction)
}

func truecondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return true
}

func ordereventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	orderEvent := tuples["orderevent"]
	orderID, _ := orderEvent.GetString("orderId")
	totalItems, _ := orderEvent.GetInt("totalItems")
	level, _ := orderEvent.GetString("level")
	itemTimeout, _ := orderEvent.GetInt("itemTimeout")
	shippingTimeout, _ := orderEvent.GetInt("shippingTimeout")

	log.RootLogger().Infof("Received a new order, processing order with id[%s]", orderID)

	//assert a order
	order, _ := model.NewTupleWithKeyValues(model.TupleType("order"), orderID)
	order.SetInt(ctx, "expectedItems", totalItems)
	order.SetString(ctx, "level", level)
	order.SetInt(ctx, "receivedItems", 0)
	order.SetString(ctx, "status", "created")
	order.SetInt(ctx, "shippingTimeout", shippingTimeout)
	order.SetDouble(ctx, "invoice", 0)

	if itemTimeout > 0 {
		orderTimeout, _ := model.NewTupleWithKeyValues(model.TupleType("ordertimeoutevent"), orderID)
		rs.ScheduleAssert(ctx, uint64(1000*itemTimeout), orderID, orderTimeout)
	}

	rs.Assert(ctx, order)

	audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "created", ruleName, fmt.Sprintf("Order[%s] created", orderID))

	if itemTimeout > 0 {
		audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "created", ruleName, fmt.Sprintf("Item timeout set to [%d]ms from now[%s]", itemTimeout, time.Now().Format("2006-01-02 15:04:05")))
	}
}

func itemeventCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	itemEvent := tuples["itemevent"]
	itemOrderID, _ := itemEvent.GetString("orderId")
	itemType, _ := itemEvent.GetString("type")
	itemQty, _ := itemEvent.GetInt("quantity")

	order := tuples["order"]
	orderID, _ := order.GetString("orderId")
	orderStatus, _ := order.GetString("status")

	currentQty := itemInventory[itemType]

	return itemOrderID == orderID && itemQty <= currentQty.quantity && orderStatus == "created"
}

func itemeventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	itemEvent := tuples["itemevent"]
	itemID, _ := itemEvent.GetString("itemId")
	itemType, _ := itemEvent.GetString("type")
	itemQty, _ := itemEvent.GetInt("quantity")

	order := tuples["order"].(model.MutableTuple)
	orderID, _ := order.GetString("orderId")

	existingReceivedItems, _ := order.GetInt("receivedItems")
	order.SetInt(ctx, "receivedItems", existingReceivedItems+1)

	currentItem := itemInventory[itemType]
	currentItem.quantity = currentItem.quantity - itemQty

	invoice, _ := order.GetDouble("invoice")
	order.SetDouble(ctx, "invoice", invoice+(float64(itemQty)*currentItem.pricePerItem))

	log.RootLogger().Infof("Received a new item id[%s] of type[%s] for order id[%s]", itemID, itemType, orderID)

	audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "created", ruleName, fmt.Sprintf("Item[%s] received", itemID))
}

func orderCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	order := tuples["order"]
	expectedItems, _ := order.GetInt("expectedItems")
	receivedItems, _ := order.GetInt("receivedItems")
	orderStatus, _ := order.GetString("status")

	return expectedItems == receivedItems && orderStatus == "created"
}

func orderAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	order := tuples["order"].(model.MutableTuple)
	receivedItems, _ := order.GetInt("receivedItems")
	orderID, _ := order.GetString("orderId")
	log.RootLogger().Infof("Received a all items[%d] for order id[%s]. Handing over to shipping", receivedItems, orderID)
	rs.CancelScheduledAssert(ctx, orderID)

	order.SetString(ctx, "status", "inshipping")

	shippingTimeout, _ := order.GetInt("shippingTimeout")
	if shippingTimeout > 0 {
		shipping, _ := model.NewTupleWithKeyValues(model.TupleType("shippingtimeoutevent"), orderID)
		rs.ScheduleAssert(ctx, uint64(1000*shippingTimeout), "shipping_"+orderID, shipping)
	}

	audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "inshipping", ruleName, fmt.Sprintf("Received all items[%d] received", receivedItems))

	if shippingTimeout > 0 {
		audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "inshipping", ruleName, fmt.Sprintf("Shipping timeout set to [%d]ms from now[%s]", shippingTimeout, time.Now().Format("2006-01-02 15:04:05")))
	}
}

func ordertimeouteventCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	ordertimeout := tuples["ordertimeoutevent"]
	timeoutOrderID, _ := ordertimeout.GetString("orderId")

	order := tuples["order"]
	orderID, _ := order.GetString("orderId")
	orderStatus, _ := order.GetString("status")

	return timeoutOrderID == orderID && orderStatus == "created"
}

func ordertimeouteventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	order := tuples["order"].(model.MutableTuple)
	orderID, _ := order.GetString("orderId")
	expectedItems, _ := order.GetInt("expectedItems")
	receivedItems, _ := order.GetInt("receivedItems")

	log.RootLogger().Infof("Order[%s] timed out, received only [%d] order items, expected [%d] items. Cancelling Order !!", orderID, receivedItems, expectedItems)
	order.SetString(ctx, "status", "cancelled")
	order.SetString(ctx, "statusMessage", fmt.Sprintf("Received only [%d] order items, expected [%d] items", receivedItems, expectedItems))

	audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "cancelled", ruleName, fmt.Sprintf("Order[%s] timed out", orderID))
}

func ordershippedCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	ordershippedevent := tuples["ordershippedevent"]
	ordershippedeventorderID, _ := ordershippedevent.GetString("orderId")

	order := tuples["order"]
	orderID, _ := order.GetString("orderId")
	orderStatus, _ := order.GetString("status")

	return orderStatus == "inshipping" && orderID == ordershippedeventorderID
}

func ordershippedAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	order := tuples["order"].(model.MutableTuple)
	orderID, _ := order.GetString("orderId")
	order.SetString(ctx, "status", "shipped")

	rs.CancelScheduledAssert(ctx, "shipping_"+orderID)

	log.RootLogger().Infof("Order[%s] has been shipped", orderID)

	audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "shipped", ruleName, fmt.Sprintf("Order[%s] shipped", orderID))
}

func ordercancelledCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	order := tuples["order"].(model.MutableTuple)
	status, _ := order.GetString("status")

	return status == "cancelled"
}

func ordercancelledAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	order := tuples["order"].(model.MutableTuple)
	orderID, _ := order.GetString("orderId")
	orderMessage, _ := order.GetString("statusMessage")

	log.RootLogger().Infof("Order id[%s] has been cancelled. Reason - %s", orderID, orderMessage)
	rs.Retract(ctx, order)

	audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "cancelled", ruleName, fmt.Sprintf("Order[%s] has been cancelled. Reason - %s", orderID, orderMessage))
}

func shippingtimeouteventCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	ordertimeout := tuples["shippingtimeoutevent"]
	timeoutOrderID, _ := ordertimeout.GetString("orderId")

	order := tuples["order"]
	orderID, _ := order.GetString("orderId")
	status, _ := order.GetString("status")

	return timeoutOrderID == orderID && status == "inshipping"
}

func shippingtimeouteventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	order := tuples["order"].(model.MutableTuple)
	orderID, _ := order.GetString("orderId")

	log.RootLogger().Infof("Shipping of Order[%s] timed out. Cancelling Order !!", orderID)
	order.SetString(ctx, "status", "cancelled")
	order.SetString(ctx, "statusMessage", fmt.Sprint("Shipping timed out"))

	audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "cancelled", ruleName, fmt.Sprintf("Order[%s] shipping timed out", orderID))
}

func orderinvoiceCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	ordertimeout := tuples["order"]
	status, _ := ordertimeout.GetString("status")

	return status == "shipped"
}

func orderinvoiceAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	order := tuples["order"].(model.MutableTuple)
	orderID, _ := order.GetString("orderId")
	invoice, _ := order.GetDouble("invoice")
	level, _ := order.GetString("level")

	discountApplied, ok := levelDiscount[level]
	if !ok {
		discountApplied = 100
	}
	discountedInvoice := invoice - ((discountApplied / 100.00) * invoice)

	log.RootLogger().Infof("Order [%s] completed. Invoice submitted,\n\tOrginal invoice - $%.2f\n\t'%s' level discount - %.2f%%\n\tFinal invoice - $%.2f", orderID, invoice, level, discountApplied, discountedInvoice)
	order.SetString(ctx, "status", "complete")
	rs.Retract(ctx, order)

	audittrail.PublishAuditTrailItem(*awsStreamName, orderID, "complete", ruleName, fmt.Sprintf("Order[%s] processing completed", orderID))
}
