package api

import (
	"Completed_Api/user-web/forms"
	"Completed_Api/user-web/global"
	"Completed_Api/user-web/global/CodesType"
	"Completed_Api/user-web/global/request"
	"Completed_Api/user-web/global/respons"
	"Completed_Api/user-web/middlewares"
	"Completed_Api/user-web/models"
	"Completed_Api/user-web/proto"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

//参数验证
func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

//TODO GetUserList服务（获取用户列表）
func GetUserList(ctx *gin.Context) {
	var res respons.GetUserInfoListRespon
	var req request.GetUserListRequst
	var userinfo respons.UserInfo
	if err := ctx.ShouldBind(&req); err != nil {
		res.Code = CodesType.Code_Err_ArgumentError
		ctx.JSON(http.StatusOK, res)
		_ = ctx.Error(err)
		return
	}
	if req.Pn == 0 || req.PSize == 0 {
		res.Code = CodesType.Code_Err_ArgumentError
		res.Desc = "参数错误"
		ctx.JSON(http.StatusOK, res)
		return
	}

	list, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    req.Pn,
		PSize: req.PSize,
	})
	if err != nil {
		zap.S().Infof("[GetUserList]获取用户列表失败,err:%s", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	var userinfos respons.UserLists

	for _, value := range list.Data {
		userinfo.ID = value.Id
		userinfo.Mobile = value.Mobile
		userinfo.Gender = value.Gender
		userinfo.Birthday = time.Time(time.Unix(int64(value.Birthday), 0))
		userinfo.NickName = value.NickName
		userinfos = append(userinfos, userinfo)
	}
	res.UserList = userinfos
	res.Code = 200
	res.Desc = "查找用户列表成功"
	ctx.JSON(http.StatusOK, res)
}

//TODO PasswordLogin（用户登录服务）
func PasswordLogin(ctx *gin.Context) {
	PassWordReq := forms.PassWordLoginForm{}
	if err := ctx.ShouldBind(&PassWordReq); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	if !store.Verify(PassWordReq.CaptchaId, PassWordReq.Captcha, true) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"captcha": "验证码错误",
		})
		return
	}
	//首先查找是否存在这个用户
	mobile, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: PassWordReq.Mobile})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "没有找到相关用户",
				})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    100,
					"message": err.Error(),
				})
				return
			}
		}
	} else {
		password, err := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:         PassWordReq.PassWord,
			EncrytedPassword: mobile.Password,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    100,
				"message": "登录失败",
			})
			return
		}
		zap.S().Infof("测试info：%s，%s,%s", PassWordReq.PassWord, mobile.Password, password.Success)
		if !password.Success {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    100,
				"message": "密码错误",
			})
			return
		} else {
			//生成token
			jwtObject := middlewares.NewJWT()
			claims := models.CustomClaims{
				ID:          uint(mobile.Id),
				NickName:    mobile.NickName,
				AuthorityId: uint(mobile.Role),
				StandardClaims: jwt.StandardClaims{
					NotBefore: time.Now().Unix(),               //设置签名生效时间
					ExpiresAt: time.Now().Unix() + 60*60*24*30, //设置过期时间
					Issuer:    "lzz",                           //签发人
				},
			}
			token, err := jwtObject.CreateToken(claims)
			if err != nil {
				zap.S().Infof("生成token错误：%s", err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"Code":  "200",
				"Desc":  "登录成功",
				"token": token,
				"Data":  mobile,
			})
		}
	}
}

//TODO Register（用户注册服务）
func Register(ctx *gin.Context) {
	RegisterForm := forms.RegisterForm{}
	//表单验证
	if err := ctx.ShouldBind(&RegisterForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserinfo{
		NickName: RegisterForm.Mobile,
		PassWord: RegisterForm.PassWord,
		Mobile:   RegisterForm.Mobile,
	})
	if err != nil {
		zap.S().Infof("[Register]Register server is fail：%s", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	jwtObject := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //设置签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //设置过期时间
			Issuer:    "lzz",                           //签发人
		},
	}
	token, err := jwtObject.CreateToken(claims)
	if err != nil {
		zap.S().Infof("生成token错误：%s", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Code":   "200",
		"Desc":   "注册成功",
		"userid": user.Id,
		"token":  token,
		"Data":   user.Mobile,
	})
}
