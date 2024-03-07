package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jicodes/go-micro-service/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(orderID uint64) string {
	return fmt.Sprintf("order:%d", orderID)
}

func (r *RedisRepo) Insert(ctx context.Context, order *model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	txn := r.Client.TxPipeline()

	err = txn.SetNX(ctx, key, string(data), 0).Err()
	if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add order to orders set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}
	
	return nil
}

var ErrNotExist = errors.New("order does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, orderID uint64) (model.Order, error) {
	key := orderIDKey(orderID)

	data, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(data), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return order, nil
}

type FindAllPage struct {
	Size uint64
	Cursor uint64
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "orders", page.Cursor, "*", int64(page.Size))

	keys, cursor, err := res.Result()

	if err != nil {
		return FindResult{}, fmt.Errorf("failed to scan orders: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}
	orders := make([]model.Order, len(xs))

	for i, x := range xs {
		x := x.(string)
		var order model.Order
		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to unmarshal order: %w", err)
		}
		orders[i] = order
	}
	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}

func (r *RedisRepo) UpdateByID(ctx context.Context, orderID uint64, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	key := orderIDKey(orderID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, orderID uint64) error {
	key := orderIDKey(orderID)

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove order from orders set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	return nil
}