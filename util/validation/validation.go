package validation

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var trans ut.Translator

// LocalTrans 修改 gin 框架中的 validator 引擎属性，实现语言定制
func LocalTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取 json 的 tag 的自定义方法，用于修改返回 error 中的错误 key 字段格式
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" { // json tag 中的一种约束
				return ""
			}
			return name
		})
		zhT := zh.New()              // 中文翻译
		enT := en.New()              // 英文翻译
		uni := ut.New(enT, zhT, enT) // 第一个参数是备用的语言环境，后面参数是应该支持的语言环境
		if trans, ok = uni.GetTranslator(locale); !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		} else {
			switch locale {
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
			case "zh":
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
			default:
				_ = en_translations.RegisterDefaultTranslations(v, trans)
			}
			return
		}
	}
	return
}

// Error 解析 validator 获得错误信息返回 string
func Error(err error) (message string) {
	if validationErrors, ok := err.(validator.ValidationErrors); !ok {
		return err.Error()
	} else {
		n := len(validationErrors)
		for i := 0; i < n; i++ {
			message += validationErrors[i].Translate(trans)
			if i < n-1 {
				message += ";"
			}
		}
	}
	return message
}
