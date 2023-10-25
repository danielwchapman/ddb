package ddb

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type options struct {
	updates         expression.UpdateBuilder
	updatesCount    int
	conditions      expression.ConditionBuilder
	conditionsCount int
	returnValues    types.ReturnValue
	returnValuesOut any

	// for use with query
	pageSize      *int32
	startKey      map[string]types.AttributeValue
	pageOut       *string
	scanBackwards bool
}

type Option func(options *options) error

// WithFieldUpdates adds field updates to the options. For use with Update.
func WithFieldUpdates(updates map[string]any) Option {
	return func(options *options) error {
		item, err := attributevalue.MarshalMap(updates)
		if err != nil {
			return fmt.Errorf("WithFieldUpdates: MarshalMap: %w", err)
		}

		for k, v := range item {
			options.updates = options.updates.Set(expression.Name(k), expression.Value(v))
		}

		options.updatesCount += len(item)

		return nil
	}
}

// WithItemExists adds a condition that the item exists. For use with Update and Put.
func WithItemExists() Option {
	return func(options *options) error {
		cond := expression.AttributeExists(expression.Name("PK"))

		if options.conditionsCount == 0 {
			options.conditions = cond
		} else {
			options.conditions.And(cond)
		}
		options.conditionsCount++

		return nil
	}
}

// WithItemNotExist adds a condition that the item does not exist. For use with Update and Put.
func WithItemNotExist() Option {
	return func(options *options) error {
		cond := expression.AttributeNotExists(expression.Name("PK"))

		if options.conditionsCount == 0 {
			options.conditions = cond
		} else {
			options.conditions.And(cond)
		}
		options.conditionsCount++

		return nil
	}
}

func WithPage(serializedPage string, out *string) Option {
	return func(options *options) error {
		startKey, err := DeserializeExclusiveStartKey(serializedPage)
		if err != nil {
			return fmt.Errorf("WithPage: %w", err)
		}
		options.startKey = startKey
		options.pageOut = out
		return nil
	}
}

// WithPageSize adds a page size for Querying.
func WithPageSize(pageSize int) Option {
	return func(options *options) error {
		if pageSize == 0 {
			return nil
		}
		if pageSize < 0 {
			return &InvalidArgumentError{err: errors.New("WithPageSize: pageSize cannot be negative")}
		}
		size := int32(pageSize)
		options.pageSize = &size
		return nil
	}
}

func WithCondition(condition expression.ConditionBuilder) Option {
	return func(options *options) error {
		if options.conditionsCount == 0 {
			options.conditions = condition
		} else {
			options.conditions.And(condition)
		}
		options.conditionsCount++

		return nil
	}
}

// WithReturnValues adds a condition to the options. For use with Update and Put.
func WithReturnValues(returnValues types.ReturnValue, out any) Option {
	return func(options *options) error {
		options.returnValues = returnValues
		options.returnValuesOut = out
		return nil
	}
}

func WithScanBackwards() Option {
	return func(options *options) error {
		options.scanBackwards = true
		return nil
	}
}
