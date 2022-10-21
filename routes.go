package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/imraan-go/aws-step-order-service/entity"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

var validState = []string{"AK", "AL", "AR", "AZ", "CA", "CO", "CT", "DC", "DE", "FL", "GA",
	"HI", "IA", "ID", "IL", "IN", "KS", "KY", "LA", "MA", "MD", "ME",
	"MI", "MN", "MO", "MS", "MT", "NC", "ND", "NE", "NH", "NJ", "NM",
	"NV", "NY", "OH", "OK", "OR", "PA", "RI", "SC", "SD", "TN", "TX",
	"UT", "VA", "VT", "WA", "WI", "WV", "WY"}

func setupRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, "Order service running successfully!")
	})
	e.GET("/tables", tablesHandler)
	e.GET("/getItem/:itemId", getItemHandler)
	e.POST("/order", orderHandler)
}

func tablesHandler(c echo.Context) error {
	resp, err := dnc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}

	fmt.Println("Tables:")
	for _, tableName := range resp.TableNames {
		fmt.Println(tableName)
	}

	return c.JSON(200, resp)

}

func getItemHandler(c echo.Context) error {
	itemId := c.Param("itemId")
	fmt.Println(itemId)

	tableName := "Inventory"

	resp, err := dnc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"ItemId": &types.AttributeValueMemberS{Value: itemId},
		},
	})

	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}

	return c.JSON(200, resp.Item)

}

func orderHandler(c echo.Context) error {
	data := &entity.CreateOrderRequest{}
	err := c.Bind(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "invalid.payload",
			"message": err.Error(),
		})
	}

	itemId := data.Order.ItemID
	fmt.Println(itemId)

	tableName := "Inventory"

	resp, err := dnc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"ItemId": &types.AttributeValueMemberS{Value: itemId},
		},
	})

	if err != nil {
		return err
	}

	if resp.Item == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"statusCode": 100,
			"body":       "Item is not in our inventory!",
		})
	}

	var item entity.Item
	err = attributevalue.UnmarshalMap(resp.Item, &item)
	if err != nil {
		return err
	}

	if item.Count < data.Order.Quantity {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"statusCode": 100,
			"body":       "Not enough inventory!",
		})
	}
	if !slices.Contains(validState, data.DeliveryDetails.ShippingAddress.StateOrRegion) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"statusCode": 100,
			"body":       "can not delivery to this state!",
		})
	}

	// Everything ok
	// Insert to order table

	_, err = dnc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Order"),
		Item: map[string]types.AttributeValue{
			"OrderId":      &types.AttributeValueMemberS{Value: data.Order.OrderID},
			"ItemId":       &types.AttributeValueMemberS{Value: data.Order.ItemID},
			"ItemName":     &types.AttributeValueMemberS{Value: data.Order.ItemName},
			"PurchaseDate": &types.AttributeValueMemberS{Value: data.Order.PurchaseDate.String()},
			"Amount":       &types.AttributeValueMemberS{Value: data.Order.OrderTotal.Amount},
			"DeliveryId":   &types.AttributeValueMemberS{Value: data.DeliveryDetails.DeliveryID},
			"PaymentId":    &types.AttributeValueMemberS{Value: data.Payment.PaymentID},
		},
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"statusCode": 100,
			"body":       "error: fail to save into database",
		})
	}

	return c.JSON(200, echo.Map{
		"success": true,
		"message": "Order palced successfully",
	})

}
