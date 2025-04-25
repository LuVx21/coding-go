package mongodb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoToDbBackuper 备份器
type MongoToDbBackuper struct {
	mongoClient   *mongo.Client
	db            *sql.DB
	batchSize     int
	typeConverter TypeConverter
}

// NewBackuper 创建备份器实例
func NewBackuper(mongoURI string, db *sql.DB, batchSize int) (*MongoToDbBackuper, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to MongoDB")
	}

	return &MongoToDbBackuper{
		mongoClient:   mongoClient,
		db:            db,
		batchSize:     batchSize,
		typeConverter: DefaultTypeConverter{},
	}, nil
}

// Backup 备份单个集合
func (b *MongoToDbBackuper) Backup(ctx context.Context, mongoDB, mongoCollection, dbTable string) error {
	collection := b.mongoClient.Database(mongoDB).Collection(mongoCollection)

	if err := b.createDbTable(ctx, collection, dbTable); err != nil {
		return errors.Wrap(err, "failed to create DB table")
	}

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return errors.Wrap(err, "failed to query MongoDB")
	}
	defer cursor.Close(ctx)

	var batch []any
	inserted := 0

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Failed to decode document: %v", err)
			continue
		}

		// 转换数据类型
		dbData, err := b.convertToDbData(doc)
		if err != nil {
			log.Printf("Failed to convert data: %v", err)
			continue
		}

		batch = append(batch, dbData)
		inserted++

		// 批量插入
		if len(batch) >= b.batchSize {
			if err := b.insertBatch(ctx, dbTable, batch); err != nil {
				return errors.Wrap(err, "failed to insert batch")
			}
			batch = batch[:0] // 清空batch
			log.Printf("Inserted %d documents into %s", inserted, dbTable)
		}
	}

	// 插入剩余数据
	if len(batch) > 0 {
		if err := b.insertBatch(ctx, dbTable, batch); err != nil {
			return errors.Wrap(err, "failed to insert final batch")
		}
		log.Printf("Inserted total %d documents into %s", inserted, dbTable)
	}

	return nil
}

// createDbTable 根据MongoDB集合结构创建DB表
func (b *MongoToDbBackuper) createDbTable(ctx context.Context, collection *mongo.Collection, tableName string) error {
	// 获取一个文档样本推断结构
	var sampleDoc bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})
	err := collection.FindOne(ctx, bson.D{}, opts).Decode(&sampleDoc)
	if err != nil {
		return errors.Wrap(err, "failed to get sample document")
	}

	// 构建CREATE TABLE语句
	columns := []string{"id VARCHAR(255) PRIMARY KEY"} // MongoDB _id 作为主键
	for key, value := range sampleDoc {
		if key == "_id" {
			continue
		}

		sqlType := b.typeConverter.GoTypeToSQL(reflect.TypeOf(value))
		columns = append(columns, fmt.Sprintf("%s %s", key, sqlType))
	}

	createSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ", "))
	_, err = b.db.ExecContext(ctx, createSQL)
	return err
}

// insertBatch 批量插入数据到DB
func (b *MongoToDbBackuper) insertBatch(ctx context.Context, tableName string, batch []any) error {
	if len(batch) == 0 {
		return nil
	}

	// 构建INSERT语句
	var placeholders []string
	var columns []string
	var values []any

	// 从第一个文档获取列名
	firstDoc := batch[0].(bson.M)
	for col := range firstDoc {
		columns = append(columns, col)
		placeholders = append(placeholders, "?")
	}

	// 构建所有值
	for _, doc := range batch {
		d := doc.(bson.M)
		for _, col := range columns {
			values = append(values, d[col])
		}
	}

	// 构建完整SQL
	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		b.buildUpdateClause(columns),
	)

	_, err := b.db.ExecContext(ctx, insertSQL, values...)
	return err
}

// buildUpdateClause 构建UPDATE子句
func (b *MongoToDbBackuper) buildUpdateClause(columns []string) string {
	var updates []string
	for _, col := range columns {
		if col == "_id" {
			continue
		}
		updates = append(updates, fmt.Sprintf("%s=VALUES(%s)", col, col))
	}
	return strings.Join(updates, ", ")
}

// convertToDbData 转换数据类型
func (b *MongoToDbBackuper) convertToDbData(doc bson.M) (bson.M, error) {
	result := make(bson.M)
	for key, value := range doc {
		converted, err := b.typeConverter.Convert(value)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert field %s", key)
		}
		result[key] = converted
	}
	return result, nil
}

// TypeConverter 类型转换接口
type TypeConverter interface {
	Convert(value any) (any, error)
	GoTypeToSQL(t reflect.Type) string
}

// DefaultTypeConverter 默认类型转换实现
type DefaultTypeConverter struct{}

func (c DefaultTypeConverter) Convert(value any) (any, error) {
	switch v := value.(type) {
	// case bson.TypeObjectID:
	// 	return v.Hex(), nil
	case primitive.DateTime:
		return v.Time(), nil
	case primitive.Decimal128:
		return v.String(), nil
	case primitive.Binary:
		return v.Data, nil
	case primitive.Timestamp:
		return time.Unix(int64(v.T), 0), nil
	case primitive.Regex:
		return v.String(), nil
	case primitive.JavaScript:
		return string(v), nil
	default:
		return value, nil
	}
}

func (c DefaultTypeConverter) GoTypeToSQL(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "TEXT"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32, reflect.Float64:
		return "DOUBLE"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return "DATETIME"
		}
		return "TEXT" // 复杂结构存储为JSON
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return "BLOB"
		}
		return "TEXT" // 数组存储为JSON
	default:
		return "TEXT"
	}
}
