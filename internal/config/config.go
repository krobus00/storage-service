package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/spf13/viper"
)

var (
	serviceName    = ""
	serviceVersion = ""
)

func ServiceName() string {
	return serviceName
}

func ServiceVersion() string {
	return serviceVersion
}

func DurableID() string {
	return fmt.Sprintf("%s-durable", serviceName)
}

func QueueGroup() string {
	return fmt.Sprintf("%s-queue-group", serviceName)
}

func Env() string {
	return viper.GetString("env")
}

func LogLevel() string {
	return viper.GetString("log_level")
}

func PortHTTP() string {
	return viper.GetString("ports.http")
}

func PortGRPC() string {
	return viper.GetString("ports.grpc")
}

func PortMetrics() string {
	return viper.GetString("ports.metrics")
}

func GracefulShutdownTimeOut() time.Duration {
	cfg := viper.GetString("graceful_shutdown_timeout")
	return parseDuration(cfg, DefaultGracefulShutdownTimeOut)
}

func DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		DatabaseUsername(),
		DatabasePassword(),
		DatabaseHost(),
		DatabaseName(),
		DatabaseSSLMode())
}

func DatabaseHost() string {
	return viper.GetString("database.host")
}

func DatabaseName() string {
	return viper.GetString("database.database")
}

func DatabaseUsername() string {
	return viper.GetString("database.username")
}

func DatabasePassword() string {
	return viper.GetString("database.password")
}

func DatabaseSSLMode() string {
	if viper.IsSet("database.sslmode") {
		return viper.GetString("database.sslmode")
	}
	return "disable"
}

func DatabasePingInterval() time.Duration {
	if viper.GetInt("database.ping_interval") <= 0 {
		return DefaultDatabasePingInterval
	}
	return time.Duration(viper.GetInt("database.ping_interval")) * time.Millisecond
}

func DatabaseRetryAttempts() float64 {
	if viper.GetInt("database.retry_attempts") > 0 {
		return float64(viper.GetInt("database.retry_attempts"))
	}
	return DefaultDatabaseRetryAttempts
}

func DatabaseMaxIdleConns() int {
	if viper.GetInt("database.max_idle_conns") <= 0 {
		return DefaultDatabaseMaxIdleConns
	}
	return viper.GetInt("database.max_idle_conns")
}

func DatabaseMaxOpenConns() int {
	if viper.GetInt("database.max_open_conns") <= 0 {
		return DefaultDatabaseMaxOpenConns
	}
	return viper.GetInt("database.max_open_conns")
}

func DatabaseConnMaxLifetime() time.Duration {
	if !viper.IsSet("database.conn_max_lifetime") {
		return DefaultDatabaseConnMaxLifetime
	}
	return time.Duration(viper.GetInt("database.conn_max_lifetime")) * time.Millisecond
}

func DatabaseConnReconnectFactor() int {
	if viper.GetInt("database.conn_reconnect_factor") <= 0 {
		return DefaultDatabaseMaxOpenConns
	}
	return viper.GetInt("database.conn_reconnect_factor")
}

func DatabaseConnReconnectMinJitter() time.Duration {
	cfg := viper.GetString("database.conn_reconnect_min_jitter")
	return parseDuration(cfg, DefaultDatabaseReconnectMinJitter)
}

func DatabaseConnReconnectMaxJitter() time.Duration {
	cfg := viper.GetString("database.conn_reconnect_Max_jitter")
	return parseDuration(cfg, DefaultDatabaseReconnectMaxJitter)
}

func DisableCaching() bool {
	return viper.GetBool("redis.disable_caching")
}

func RedisCacheHost() string {
	return viper.GetString("redis.cache_host")
}

func RedisDialTimeout() time.Duration {
	cfg := viper.GetString("redis.dial_timeout")
	return parseDuration(cfg, DefaultRedisDialTimeout)
}

func RedisWriteTimeout() time.Duration {
	cfg := viper.GetString("redis.write_timeout")
	return parseDuration(cfg, DefaultRedisWriteTimeout)
}

func RedisReadTimeout() time.Duration {
	cfg := viper.GetString("redis.read_timeout")
	return parseDuration(cfg, DefaultRedisReadTimeout)
}

func RedisCacheTTL() time.Duration {
	cfg := viper.GetString("cache_ttl")
	return parseDuration(cfg, DefaultRedisCacheTTL)
}

func GetS3Region() string {
	return viper.GetString("s3.region")
}

func GetS3Endpoint() string {
	return viper.GetString("s3.endpoint")
}

func GetS3BucketName() string {
	return viper.GetString("s3.bucket")
}

func GetS3AccessKey() string {
	return viper.GetString("s3.access_key")
}

func GetS3SecretKey() string {
	return viper.GetString("s3.secret_key")
}

func GetS3SignDuration() time.Duration {
	cfg := viper.GetString("s3.sign_duration")
	return parseDuration(cfg, DefaultS3SignDuration)
}

func GetS3Credential() *credentials.StaticCredentialsProvider {
	accessKeyIDValue := GetS3AccessKey()
	secretAccessKeyValue := GetS3SecretKey()
	if accessKeyIDValue != "" && secretAccessKeyValue != "" {
		credentials := credentials.NewStaticCredentialsProvider(accessKeyIDValue, secretAccessKeyValue, "")
		return &credentials
	}
	return nil
}

func JetstreamHost() string {
	return viper.GetString("js.host")
}

func JetstreamMaxPending() int {
	if viper.GetInt("js.max_pending") <= 0 {
		return DefaultJetstreamMaxPending
	}
	return viper.GetInt("js.max_pending")
}

func JetstreamMaxAge() time.Duration {
	cfg := viper.GetString("js.max_age")
	return parseDuration(cfg, DefaultJetstreamMaxAge)
}

func AuthGRPCHost() string {
	return viper.GetString("services.auth.grpc")
}

func JaegerProtocol() string {
	return viper.GetString("jaeger.protocol")
}

func JaegerHost() string {
	return viper.GetString("jaeger.host")
}

func JaegerPort() string {
	return viper.GetString("jaeger.port")
}

func JaegerSampleRate() float64 {
	return viper.GetFloat64("jaeger.sample_rate")
}

func parseDuration(in string, defaultDuration time.Duration) time.Duration {
	dur, err := time.ParseDuration(in)
	if err != nil {
		return defaultDuration
	}
	return dur
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("config not found")
		}
		return err
	}
	return nil
}
