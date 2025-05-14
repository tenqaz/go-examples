
type RedisMutex struct {
	redisClient *redis.ClusterClient
	name        string
	expire      time.Duration
	autoRenewal bool
	// 通知续约退出
	donec chan struct{}
}

func NewRedisMutex(client *redis.ClusterClient, name string, options ...Option) *RedisMutex {
	donec := make(chan struct{})
	redisMutex := &RedisMutex{
		redisClient: client,
		name:        name,
		expire:      60 * time.Second,
		autoRenewal: false,
		donec:       donec,
	}

	for _, option := range options {
		option.Apply(redisMutex)
	}

	return redisMutex
}

// Lock 获取锁
func (m *RedisMutex) Lock() (bool, error) {
	value := uuid.New().String()
	var isOk bool
	var err error
	err = RedisWithRetry(func() error {
		isOk, err = m.redisClient.SetNX(context.Background(), m.name, value, m.expire).Result()
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	if !isOk {
		return false, nil
	}
	logger.Infof("redis setnx success, name: %s, value: %s", m.name, value)

	if m.autoRenewal {
		go func() {
			err := m.autoRenewalLock()
			if err != nil {
				logger.Errorf("autoRenewalLock failed, err: %+v", err)
			}
		}()
	}
	return true, nil
}

// Unlock 释放锁
func (m *RedisMutex) Unlock() error {
	logger.Infof("redis del, name: %s", m.name)
	if m.autoRenewal {
		close(m.donec)
	}
	err := RedisWithRetry(func() error {
		err := m.redisClient.Del(context.Background(), m.name).Err()
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// Extend 续约
func (m *RedisMutex) Extend() error {
	logger.Infof("redis expire, name: %s", m.name)
	err := RedisWithRetry(func() error {
		err := m.redisClient.Expire(context.Background(), m.name, m.expire).Err()
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// autoRenewalLock 自动续约
func (m *RedisMutex) autoRenewalLock() error {
	for {
		select {
		case <-m.donec:
			return nil
		case <-time.After(m.expire * 2 / 3):
			err := m.Extend()
			if err != nil {
				return err
			}
		}
	}
}

type Option interface {
	Apply(*RedisMutex)
}

type OptionFunc func(*RedisMutex)

func (f OptionFunc) Apply(mutex *RedisMutex) {
	f(mutex)
}

// WithExpire 设置过期时间
func WithExpire(expire time.Duration) Option {
	return OptionFunc(func(m *RedisMutex) {
		m.expire = expire
	})
}

// WithAutoRenewal 是否自动续约
func WithAutoRenewal(a bool) Option {
	return OptionFunc(func(m *RedisMutex) {
		m.autoRenewal = a
	})
}
