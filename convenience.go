package ddb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func DeserializeExclusiveStartKey(key string) (map[string]types.AttributeValue, error) {
	if key == "" {
		return nil, nil
	}

	jsonBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("DeserializeExclusiveStartKey: StdEncoding.DecodeString %w", err)
	}

	outputJSON := map[string]any{}
	if err := json.Unmarshal(jsonBytes, &outputJSON); err != nil {
		return nil, fmt.Errorf("deserializeExclusiveStartKey: json.Unmarshal %w", err)
	}

	return attributevalue.MarshalMap(outputJSON)
}

func SerializeExclusiveStartKey(key map[string]types.AttributeValue) (string, error) {
	if len(key) == 0 {
		return "", nil
	}

	convertedMap := map[string]any{}
	if err := attributevalue.UnmarshalMap(key, &convertedMap); err != nil {
		return "", fmt.Errorf("SerializeExclusiveStartKey: attributevalue.UnmarshalMap %w", err)
	}

	jsonBytes, err := json.Marshal(convertedMap)
	if err != nil {
		return "", fmt.Errorf("SerializeExclusiveStartKey: json.Marshal %w", err)
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

func marshalMapList(items []any) ([]map[string]types.AttributeValue, error) {
	if len(items) == 0 {
		return nil, nil
	}

	var (
		out = make([]map[string]types.AttributeValue, len(items))
		err error
	)

	for i := range items {
		out[i], err = attributevalue.MarshalMap(items[i])
		if err != nil {
			return nil, fmt.Errorf("MarshalMapList: MarshalMap: %w", err)
		}
	}

	return out, nil
}

func makePuts(table string, rows ...PutRow) ([]types.Put, error) {
	items := make([]types.Put, len(rows))
	for i := range rows {
		item, err := attributevalue.MarshalMap(rows[i].Row)
		if err != nil {
			return nil, fmt.Errorf("TransactionPuts: MarshalMap: %w", err)
		}

		items[i] = types.Put{
			Item:                item,
			ConditionExpression: rows[i].Condition,
			TableName:           &table,
		}
	}

	return items, nil
}

func makeDeletes(table string, rows ...DeleteRow) []types.Delete {
	items := make([]types.Delete, len(rows))
	for i := range rows {
		items[i] = types.Delete{
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: rows[i].PK},
				"SK": &types.AttributeValueMemberS{Value: rows[i].SK},
			},
			ConditionExpression: rows[i].Condition,
			TableName:           &table,
		}
	}
	return items
}

func makeTransactionWriteItems(
	puts []types.Put,
	deletes []types.Delete,
	updates []types.Update,
) []types.TransactWriteItem {
	out := make([]types.TransactWriteItem, len(puts)+len(deletes)+len(updates))
	i := 0
	for j := range puts {
		out[i] = types.TransactWriteItem{
			Put: &puts[j],
		}
		i++
	}
	for j := range deletes {
		out[i] = types.TransactWriteItem{
			Delete: &deletes[j],
		}
		i++
	}
	for j := range updates {
		out[i] = types.TransactWriteItem{
			Update: &updates[j],
		}
		i++
	}
	return out
}
