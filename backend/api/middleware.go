package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"ticket/backend/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader(authorizationHeaderKey)
		if len(header) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(header)
		if len(fields) != 2 {
			err := errors.New("authorization header invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// It's probably better to put it in config
const (
	refillIntervalInSecond int64 = 1
	bucketSize             int64 = 5
	maxTxRetries           int   = 5
)

var exceedRateLimitErr error = fmt.Errorf("User has exceed request limit, please try again later.")

func RateLimitMiddleware(cache *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Token Bucket Strategy:
		// Bucket 預設有 N 個 Token，每個 request 消耗一個
		// 每秒補充 M 個 Token
		// 每次 request 要拿 Token 前檢查從上次補充 (last_refill_time) 到現在過了幾秒
		// 來決定要補充幾個，補充完再拿 Token

		payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

		lastRefillTimeKey := fmt.Sprintf("%d_last_refill_time", payload.UserID)
		userBucketKey := fmt.Sprintf("%d_bucket", payload.UserID)

		txf := func(tx *redis.Tx) error {
			userRequestLeft, err := tx.Get(ctx, userBucketKey).Int64()
			if err == redis.Nil {
				userRequestLeft = bucketSize
			} else if err != nil {
				return err
			}

			lastRefillTime, err := tx.Get(ctx, lastRefillTimeKey).Int64()

			if err == redis.Nil {
				// Never requested before: Initialize
				lastRefillTime = time.Now().Unix()
			} else if err != nil {
				return err
			} else {
				// Refill Token
				userRequestLeft = min(userRequestLeft+(time.Now().Unix()-lastRefillTime)/refillIntervalInSecond, bucketSize)
				lastRefillTime = time.Now().Unix()
			}

			if userRequestLeft > 0 {
				userRequestLeft--
			} else {
				return exceedRateLimitErr
			}

			_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
				p.Set(ctx, lastRefillTimeKey, lastRefillTime, 0)
				p.Set(ctx, userBucketKey, userRequestLeft, 0)
				return nil
			})

			return err
		}

		for i := 0; i < maxTxRetries; i++ {
			err := cache.Watch(ctx, txf, lastRefillTimeKey, userBucketKey)

			if err == nil {
				// Success
				break
			}

			if err == redis.TxFailedErr {
				// Optimistic Lock Trigger
				continue
			}

			if err == exceedRateLimitErr {
				// Rate limit Exceeded
				ctx.AbortWithStatusJSON(http.StatusTooManyRequests, errorResponse(err))
				return
			}

			// Something unexpected happened
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.Next()
	}
}
