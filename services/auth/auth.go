package auth

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type claimKey struct{}

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	window   time.Duration
	maxReqs  int
}

func NewRateLimiter(window time.Duration, max int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		maxReqs:  max,
	}
}

func (rl *RateLimiter) filterOldRequests(reqs []time.Time) []time.Time {
	now := time.Now()
	valid := []time.Time{}

	for _, t := range reqs {
		if now.Sub(t) < rl.window {
			valid = append(valid, t)
		}
	}

	return valid
}

func FromCtx(ctx context.Context) Claims {
	claims, ok := ctx.Value(claimKey{}).(Claims)
	if !ok {
		return Claims{}
	}

	return claims
}

func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.requests[ip] = rl.filterOldRequests(rl.requests[ip])
	if len(rl.requests[ip]) >= rl.maxReqs {
		return false
	}

	rl.requests[ip] = append(rl.requests[ip], time.Now())
	return true
}

func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !rl.isAllowed(r.RemoteAddr) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func GenerateJWT(email string) (string, error) {
	claims := Claims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "LON",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}

func validateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func GetEmailFromToken(tokenStr string) (string, error) {
	claims, err := validateJWT(tokenStr)
	if err != nil {
		return "", nil
	}
	return claims.Email, nil
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No auth token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next(w, r.WithContext(ctx))
	}
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authorization := md["authorization"]
	if len(authorization) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	claims, err := validateJWT(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token")
	}

	newCtx := context.WithValue(ctx, claimKey{}, *claims)
	return handler(newCtx, req)
}

func WithAuth(ctx context.Context, r *http.Request) metadata.MD {
	md := make(map[string]string)
	if auth := r.Header.Get("Authorization"); auth != "" {
		md["authorization"] = auth
	}
	return metadata.New(md)
}

func sanitizeEmail(email string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(email, ".", "_"),
		"@", "_at_",
	)
}
