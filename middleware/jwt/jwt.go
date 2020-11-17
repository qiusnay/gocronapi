package jwt

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/qiusnay/gocronapi/models/db"
	"github.com/qiusnay/gocronapi/util"
)

type (
	errorHandler func(iris.Context, string)

	TokenExtractor func(iris.Context) (string, error)

	Jwt struct {
		Config JwtConfig
	}
)

var Myjwt *Jwt

var lock sync.Mutex

func Serve(ctx iris.Context) bool {
	CreateJwt()
	if err := Myjwt.CheckToken(ctx); err != nil {
		golog.Errorf("Check jwt error, %s", err)
		return false
	}
	return true
}

// token校验
func (m *Jwt) CheckToken(ctx iris.Context) error {
	//option请求不做校验
	if !m.Config.EnableAuthOnOptions {
		if ctx.Method() == iris.MethodOptions {
			fmt.Println(1234)
			return nil
		}
	}
	// 从头信息中获取token
	token, err := m.Config.Extractor(ctx)
	if err != nil {
		m.logf("Error extracting JWT: %v", err)
		m.Config.ErrorHandler(ctx, "token不存在或header设置不正确")
		return fmt.Errorf("Error extracting token: %v", err)
	}
	if token == "" {
		// token丢失
		m.Config.ErrorHandler(ctx, "解析错误,token为空")
		return fmt.Errorf("解析错误,token为空")
	}
	parsedToken, err := jwt.Parse(token, m.Config.ValidationKeyGetter)
	// Check if there was an error in parsing...
	if err != nil {
		m.logf("Error parsing token: %v", err)
		m.Config.ErrorHandler(ctx, "token解析失败")
		return fmt.Errorf("Error parsing token: %v", err)
	}
	if m.Config.SigningMethod != nil && m.Config.SigningMethod.Alg() != parsedToken.Header["alg"] {
		message := fmt.Sprintf("Expected %s signing method but token specified %s",
			m.Config.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		m.logf("Error validating token algorithm: %s", message)
		m.Config.ErrorHandler(ctx, "token解析错误") // 算法错误
		return fmt.Errorf("Error validating token algorithm: %s", message)
	}

	// Check if the parsed token is valid...
	if !parsedToken.Valid {
		m.logf("解析错误,token无效")
		m.Config.ErrorHandler(ctx, "解析错误,token无效")
		return fmt.Errorf("解析错误,token无效")
	}

	if m.Config.Expiration {
		// fmt.Println(fmt.Sprintf("token expires : %+v", m.Config.Expiration))
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			if expired := claims.VerifyExpiresAt(time.Now().Unix(), true); !expired {
				return fmt.Errorf("token己过期")
			}
		}
	}
	ctx.Values().Set(m.Config.ContextKey, parsedToken)
	return nil
}

func (m *Jwt) DecodeToken(ctx iris.Context, token string) (u db.User, e error) {
	parsedToken, err := jwt.Parse(token, m.Config.ValidationKeyGetter)
	if err != nil {
		m.logf("Error parsing token1: %v", err)
		m.Config.ErrorHandler(ctx, "token解析失败")
		return u, err
	}
	tokenUser, ok := parsedToken.Claims.(jwt.MapClaims)
	//解析token
	if ok {
		u.Id = int64(tokenUser["id"].(float64))
		u.Username = tokenUser["username"].(string)
		return u, nil
	}
	return u, errors.New("token解析失败")
}

func (m *Jwt) logf(format string, args ...interface{}) {
	// fmt.Printf("aaaaaaaaaaaaaaa : %+v", args...)
	if m.Config.Debug {
		fmt.Printf(format, args...)
	}
}

func (m *Jwt) FromAuthHeader(ctx iris.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}
	return authHeaderParts[1], nil
}

// jwt中间件
func CreateJwt() {
	if Myjwt != nil {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	if Myjwt != nil {
		return
	}
	Myjwt = &Jwt{Config: JwtConfig{
		ContextKey: DefaultContextKey,
		//这个方法将验证jwt的token
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//自己加密的秘钥或者说盐值
			return []byte(util.MyAppConfig.Secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(ctx iris.Context, errMsg string) {
			util.JsonFail(ctx, iris.StatusUnauthorized, errMsg, nil)
		},
		Extractor:           Myjwt.FromAuthHeader,
		Expiration:          true,
		Debug:               true,
		EnableAuthOnOptions: false,
	}}
}

type Claims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	//Password string `json:"password"`
	//User models.User `json:"user"`
	jwt.StandardClaims
}

// 在登录成功的时候生成token
func GenerateToken(user *db.User) (string, error) {
	expireTime := time.Now().Add(time.Duration(util.MyAppConfig.JWTTimeout) * time.Second)

	claims := Claims{
		user.Id,
		user.Username,
		//user.Password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "iris-casbins-jwt",
		},
	}
	// fmt.Println(fmt.Sprintf("sdfsdfsdfsfdfdf : %+v", claims))
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(util.MyAppConfig.Secret))
	return token, err
}
