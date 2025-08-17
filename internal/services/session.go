package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID    int32     `json:"user_id"`
	Email     string    `json:"email"`
	LoginTime time.Time `json:"login_time"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SessionStore struct {
	Client *redis.Client
}

const SESSION_ID_KEY_PREFIX = "__session_id__"

func (r *SessionStore) formatSessionID(sessionID string) string {
	return fmt.Sprintf("%s%s", SESSION_ID_KEY_PREFIX, sessionID)
}

func (r *SessionStore) Set(sessionID string, data SessionData) error {
	json, _ := json.Marshal(data)
	sessionID = r.formatSessionID(sessionID)
	return r.Client.Set(context.Background(), sessionID, json, 24*time.Hour).Err()
}

func (r *SessionStore) Get(sessionID string) (*SessionData, error) {
	sessionID = r.formatSessionID(sessionID)
	val, err := r.Client.Get(context.Background(), sessionID).Result()
	if err != nil {
		return nil, err
	}

	var data SessionData
	err = json.Unmarshal([]byte(val), &data)
	return &data, err
}

func (r *SessionStore) Expire(sessionID string) (error) {
	sessionID = r.formatSessionID(sessionID)
	_, err := r.Client.Del(context.Background(), sessionID).Result()
	return err
}

func (r *SessionStore) CreateSession(data SessionData) (string, error) {
	sessionID := uuid.New().String()
	err := r.Set(sessionID, data)
	return sessionID, err
}
