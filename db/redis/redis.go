package redis

import (
	"encoding/json"
	"time"

	"github.com/carlware/go-common/errors"
	"github.com/go-redis/redis"
)

var purge = `
local keys=redis.call("KEYS", ARGV[1])
local total = 0
for k,_ in ipairs(keys) do
   redis.call("del",keys[k])
   total = total + 1
end
return total`

type Cache struct {
	Client    *redis.Client
	namespace string
}

// NewRedisClient instances a new redis client
func NewRedisClient(host, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}

func New(client *redis.Client, namespace string) *Cache {
	return &Cache{
		Client:    client,
		namespace: namespace,
	}
}

func (c *Cache) Get(key string, model interface{}) error {
	op := "redis.Get"

	value, err := c.Client.Get(c.namespace + ":" + key).Bytes()
	if err == redis.Nil {
		return errors.New(errors.NotFound, op, "not found", err)
	}
	if err != nil {
		return errors.New(errors.Internal, op, "", err)
	}
	err = json.Unmarshal(value, &model)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Set(key string, model interface{}, ttl int) error {
	op := "redis.Set"

	encoded, err := json.Marshal(model)
	if err != nil {
		return err
	}
	err = c.Client.Set(c.namespace+":"+key, encoded, time.Duration(ttl)*time.Millisecond).Err()
	if err != nil {
		return errors.New(errors.Internal, op, "", err)
	}
	return nil
}

func (c *Cache) Delete(key string) error {
	op := "redis.Delete"

	err := c.Client.Del(c.namespace + ":" + key).Err()
	if err != nil {
		return errors.New(errors.Internal, op, "", err)
	}
	return nil
}

func (c *Cache) Purge() (int64, error) {
	op := "redis.Purge"
	purge := redis.NewScript(purge)

	t, err := purge.Run(c.Client, []string{}, c.namespace+"*").Result()
	if err != nil {
		return 0, errors.New(errors.Internal, op, "", err)
	}
	switch t := t.(type) {
	case int64:
		return t, nil
	default:
		return 0, nil
	}
}
