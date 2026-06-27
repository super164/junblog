package jwt

import (
	"errors"
	"time"

	"blog_backend/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明
type Claims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

// GenerateToken 生成 access token
func GenerateToken(userID uint, username, role string) (string, error) {
	return generateTokenByType(userID, username, role, "access",
		time.Duration(config.Get().JWT.AccessExpire)*time.Minute)
}

// GenerateRefreshToken 生成 refresh token
func GenerateRefreshToken(userID uint, username, role string) (string, error) {
	return generateTokenByType(userID, username, role, "refresh",
		time.Duration(config.Get().JWT.RefreshExpire)*time.Hour)
}

// generateTokenByType 通用 token 生成
func generateTokenByType(userID uint, username, role, tokenType string, expireDuration time.Duration) (string, error) {
	cfg := config.Get()
	if cfg == nil {
		return "", errors.New("配置未初始化")
	}
	if cfg.JWT.Secret == "" {
		return "", errors.New("jwt secret 未配置")
	}

	now := time.Now()
	claims := &Claims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
			Issuer:    cfg.JWT.Issuer,
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析并校验 token（不限制 token 类型）
func ParseToken(tokenString string) (*Claims, error) {
	return parseTokenWithValidation(tokenString, false, "")
}

// ParseAccessToken 解析并校验 access token
func ParseAccessToken(tokenString string) (*Claims, error) {
	return parseTokenWithValidation(tokenString, true, "access")
}

// ParseRefreshToken 解析并校验 refresh token
func ParseRefreshToken(tokenString string) (*Claims, error) {
	return parseTokenWithValidation(tokenString, true, "refresh")
}

// parseTokenWithValidation 通用 token 解析与校验
func parseTokenWithValidation(tokenString string, checkType bool, expectedType string) (*Claims, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, errors.New("配置未初始化")
	}
	if cfg.JWT.Secret == "" {
		return nil, errors.New("jwt secret 未配置")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token 签名算法不正确")
		}
		return []byte(cfg.JWT.Secret), nil
	}, jwt.WithIssuer(cfg.JWT.Issuer))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token 无效")
	}

	// 校验 token 类型
	if checkType && claims.TokenType != expectedType {
		return nil, errors.New("token 类型不匹配")
	}

	return claims, nil
}
