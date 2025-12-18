package cache_redis_config

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/go-redis/redis/v8"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/BurntSushi/toml"
)

// GetRedisConnection returns a new Redis client
func GetRedisConnection(redisConfig RedisConfig) (*redis.Client, error) {
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})

	// Establish a connection to Redis
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetTracer returns a new OpenTelemetry tracer
func GetTracer() (trace.Tracer, func()) {
	// Create a new Jaeger exporter
	exporter, err := jaeger.NewExporter(jaeger.WithCollectorEndpoint(
		japerger.WithEndpoint("http://jaeger-agent:6831/api/traces"),
	))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new tracer builder
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			"com.example.cache-redis-config",
			attribute.String("version", "1.0"),
		)),
	)

	// Get a new tracer
	tracer := otel.Tracer("cache-redis-config")

	// Return the tracer and a function to close it
	return tracer, func() {
		tracerProvider.Shutdown(context.Background())
	}
}

// GetGRPCServer returns a new GRPC server
func GetGRPCServer(redisConfig RedisConfig) (*grpc.Server, *redis.Client) {
	// Create a new GRPC server
	server := grpc.NewServer()

	// Create a new Redis client
	client, err := GetRedisConnection(redisConfig)
	if err != nil {
		log.Fatal(err)
	}

	return server, client
}

// LoadConfig loads the Redis configuration from a TOML file
func LoadConfig(configPath string) (RedisConfig, error) {
	// Load the Redis configuration from the TOML file
	var config RedisConfig
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return RedisConfig{}, err
	}

	return config, nil
}

// RedisConfig is the Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Db       int
}