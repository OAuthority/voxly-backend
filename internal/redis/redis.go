package redis 

import( 
	"github.com/redis/go-redis/v9"
	"sync"
	"context"
	"fmt"
	"time"
	"encoding/json"
)

// struct for the config object
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// struct for SessionManagement through redis
type SessionManager struct {
	client *redis.Client
}

// Struct to represent a session stored in Redis
type Session struct {
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

var (
	sessionManager *SessionManager
	once          sync.Once
)

// Initialize a new Redis client with the given configuration
// and return any errors we may encounter along the way
func Initialize(config Config) error {
	var initErr error
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
			Password: config.Password,
			DB:       config.DB,
		})

		// Test the connection by pinging redis, if we get an error then something erronous has happened
		// along the way. We call this function in the main setup so we need to ensure no error is returned.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := client.Ping(ctx).Result()
		if err != nil {
			initErr = fmt.Errorf("failed to connect to Redis: %v", err)
			return
		}

		sessionManager = &SessionManager{
			client: client,
		}
	})

	return initErr
}

// Get a connection to Redis' SessionManager
func GetConnection() (*SessionManager, error) {
	if sessionManager == nil {
		return nil, fmt.Errorf("redis connection not initialized")
	}
	return sessionManager, nil
}

// Create a new session in redis!
func (sm *SessionManager) CreateSession(sessionId string, userId string, duration time.Duration) (*Session, error) {
	session := &Session{
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	data, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session: %v", err)
	}

	ctx := context.Background()
	err = sm.client.Set(ctx, fmt.Sprintf("session:%s", sessionId), data, duration).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store session: %v", err)
	}

	return session, nil
}

// Retrieve a session from Redis 
func (sm *SessionManager) GetSession(sessionId string) (*Session, error) {
	ctx := context.Background()
	data, err := sm.client.Get(ctx, fmt.Sprintf("session:%s", sessionId)).Result()
	if err != nil {
		if err == redis.Nil {
			// we didn't find a session in redis for this user
			return nil, nil 
		}
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	var session Session
	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %v", err)
	}

	return &session, nil
}

// Extend a session in Redis
func (sm *SessionManager) ExtendSession(sessionId string, duration time.Duration) error {
	session, err := sm.GetSession(sessionId)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found")
	}

	session.ExpiresAt = time.Now().Add(duration)
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %v", err)
	}

	ctx := context.Background()
	return sm.client.Set(ctx, fmt.Sprintf("session:%s", sessionId), data, duration).Err()
}

// Remove a session from Redis
// useful if we need to somehow log everyone out!
func (sm *SessionManager) DeleteSession(sessionId string) error {
	ctx := context.Background()
	return sm.client.Del(ctx, fmt.Sprintf("session:%s", sessionId)).Err()
}