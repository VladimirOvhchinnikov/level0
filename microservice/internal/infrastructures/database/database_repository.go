package database

import (
	"context"
	"fmt"
	"microservice/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type DataBase struct {
	logger *zap.Logger
	dbPool *pgxpool.Pool
}

func NewDataBase(logger *zap.Logger, dbPool *pgxpool.Pool) *DataBase {
	return &DataBase{
		logger: logger,
		dbPool: dbPool,
	}
}

func (db *DataBase) GetDataByID(ctx context.Context, id int) (models.Order, error) {

	var order models.Order

	orderQuery, args, err := squirrel.Select("*").From("Orders").Where(squirrel.Eq{"order_uid": fmt.Sprint(id)}).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		db.logger.Info("Ошибка создания SQL запроса: " + err.Error())
		return models.Order{}, err
	}
	err = db.dbPool.QueryRow(ctx, orderQuery, args...).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {

			db.logger.Info("Запись не найдена")
			db.logger.Error(err.Error())
			return models.Order{}, nil
		} else {

			db.logger.Info("Ошибка при выполнении запроса: " + err.Error())
			return models.Order{}, err
		}
	}

	deliveryQuery, args, err := squirrel.Select("nameD", "phone", "zip", "city", "addressD", "region", "email").From("Delivery").Where(squirrel.Eq{"order_uid": id}).ToSql()
	if err != nil {
		db.logger.Info("Ошибка создания SQL запроса доставки: " + err.Error())
		return models.Order{}, err
	}

	var delivery models.Delivery
	err = db.dbPool.QueryRow(ctx, deliveryQuery, args...).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			db.logger.Info("Ошибка при выполнении запроса доставки: " + err.Error())
			return models.Order{}, err
		} else {

			db.logger.Info("Ошибка при выполнении запроса: " + err.Error())
			return models.Order{}, err
		}
	}

	order.Delivery = delivery

	paymentQuery, args, err := squirrel.Select(
		"transactionD", "request_id", "currency", "providerD",
		"amount", "payment_dt", "bank", "delivery_cost",
		"goods_total", "custom_fee",
	).From("Payment").Where(squirrel.Eq{"order_uid": id}).ToSql()

	if err != nil {
		db.logger.Info("Ошибка создания SQL запроса оплаты: " + err.Error())
		return models.Order{}, err
	}

	var payment models.Payment
	err = db.dbPool.QueryRow(ctx, paymentQuery, args...).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			db.logger.Info("Ошибка при выполнении запроса оплаты: " + err.Error())
			return models.Order{}, err
		} else {

			db.logger.Info("Ошибка при выполнении запроса: " + err.Error())
			return models.Order{}, err
		}
	}

	order.Payment = payment

	// Получение данных товаров
	itemsQuery, args, err := squirrel.Select(
		"chrt_id", "track_number", "price", "rid", "nameD",
		"sale", "sizeD", "total_price", "nm_id", "brand", "statusD",
	).From("Items").Where(squirrel.Eq{"order_uid": id}).ToSql()

	if err != nil {
		db.logger.Info("Ошибка создания SQL запроса товаров: " + err.Error())
		return models.Order{}, err
	}

	rows, err := db.dbPool.Query(ctx, itemsQuery, args...)
	if err != nil {
		db.logger.Info("Ошибка при выполнении запроса товаров: " + err.Error())
		return models.Order{}, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err = rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			db.logger.Info("Ошибка при сканировании строки товара: " + err.Error())
			return models.Order{}, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		db.logger.Info("Ошибка при обработке результатов запроса товаров: " + err.Error())
		return models.Order{}, err
	}

	order.Items = items

	return order, nil

}
