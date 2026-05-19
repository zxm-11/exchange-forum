package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword函数用于加密密码
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12) //加密密码,生成12位盐
	return string(hash), err
}

// GenerateJWT函数用于生成JWT token
func GenerateJWT(username string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //创建jwt token对象,参1:签名方法,参2:有效载荷
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), //设置过期时间为72小时
	})

	Signedtoken, err := token.SignedString([]byte("secret")) //token原本是[]byte类型的,signdedstring添加了密钥secret又把token转化成了字符串
	return "Bearer " + Signedtoken, err
}

// CheckPasswordHash函数用于检查密码是否匹配
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) //比较密码和哈希值是否匹配
	return err == nil
}

// ParseJWT函数用于解析JWT token
func ParseJWT(tokenstring string) (string, error) {
	if len(tokenstring) > 7 && tokenstring[:7] == "Bearer " {
		tokenstring = tokenstring[7:] //去掉Bearer前缀
	}

	//解析器
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		} //类型断言,返回值1:转换成功后具体类型的值  返回值2:类型断言是否成功
		return []byte("secret"), nil //返回密钥
	})

	if err != nil {
		return "", err
	}

	//类型断言,判断token.Claims是否是jwt.MapClaims类型,并且token.Valid是否为true
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	}
	return "", errors.New("invalid token")
}
