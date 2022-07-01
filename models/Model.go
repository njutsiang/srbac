package models

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"reflect"
	"srbac/app"
	"srbac/utils"
	"srbac/code"
	"srbac/exception"
	"strings"
)

// 数据模型基础
// 表单模型和数据库数据模型，都可以使用该基类
type Model struct {
	translator ut.Translator
	validate *validator.Validate
	error error
	errorMessages map[string]string
	refValue reflect.Value
	db *gorm.DB
}

// 实例化验证器
func (this *Model) NewValidate() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("label")
	})
	return validate
}

// 实例化翻译器
func (this *Model) NewTranslator() ut.Translator {
	translator, found := ut.New(en.New(), zh.New()).GetTranslator("zh")
	if !found {
		app.Error("指定的语言不存在")
	}
	return translator
}

// 注册自定义验证方法
func (this *Model) RegisterValidation(tag string, fn func(field validator.FieldLevel) bool, messages... string) {
	if err := this.validate.RegisterValidation(tag, fn); err != nil {
		app.Error(err)
	}
	message := tag + " error"
	if len(messages) >= 1 {
		message = messages[0]
	}
	this.RegisterTranslation(tag, message)
}

// 注册错误提示语
func (this *Model) RegisterTranslation(tag string, message string) {
	if err := this.validate.RegisterTranslation(tag, this.translator, func(ut ut.Translator) error {
		return ut.Add(tag, message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(tag, fe.Field(), fe.Param())
		if err != nil {
			app.Error(err)
		}
		return t
	}); err != nil {
		app.Error(err)
	}
}

// 初始化验证器
func (this *Model) InitValidate() {
	this.translator = this.NewTranslator()
	this.validate = this.NewValidate()
	this.InitTranslation()
}

// 初始化默认错误提示语
func (this *Model) InitTranslation() {
	for tag, message := range map[string]string{
		"required": "{0}不能为空",
		"min": "{0}的长度不能小于{1}个字符",
		"max": "{0}的长度不能大于{1}个字符",
	} {
		this.RegisterTranslation(tag, message)
	}
}

// 往数据模型载入数据
// data 中不存在的键，或值为 null 的成员，读取出来都是 nil
// 反射，只能获取到结构体的公共成员
func (this *Model) setAttributes(data map[string]interface{}) {
	elemValue := this.refValue.Elem()
	types := elemValue.Type()
	for i := 0; i < elemValue.NumField(); i ++ {
		name := utils.SnakeString(types.Field(i).Name)
		value, isSet := data[name]
		if isSet {
			switch types.Field(i).Type.String() {
			case "int64":
				if utils.IsNumeric(value) {
					elemValue.Field(i).SetInt(utils.ToInt64(value))
				}
			case "string":
				if utils.IsString(value) {
					elemValue.Field(i).SetString(utils.ToString(value))
				}
			case "sql.NullInt64":
				if utils.IsNumeric(value) || value == nil {
					elemValue.Field(i).Set(reflect.ValueOf(utils.ToNullInt64(value)))
				}
			default:
				exception.NewException(code.UnknownModelFieldType)
			}
		}
	}
}

// 获取验证器实例
func (this *Model) GetValidate() *validator.Validate {
	return this.validate
}

// 写入错误信息
func (this *Model) SetError(err error) {
	this.error = err
}

// 写入错误提示语
func (this *Model) SetErrorMessages(messages map[string]string) {
	if this.errorMessages == nil {
		this.errorMessages = map[string]string{}
	}
	for key, msg := range messages {
		this.errorMessages[key] = msg
	}
}

// 获取第一条错误提示语
// separators，默认 <br>，适用于 HTML 页面，多条错误提示语的分隔符
func (this *Model) GetError(separators... string) string {
	if this.error == nil {
		return ""
	} else {
		message := this.error.Error()
		switch this.error.(type) {
		case *validator.InvalidValidationError:
			message = this.error.(*validator.InvalidValidationError).Error()
		case validator.ValidationErrors:
			separator := "<br>"
			if len(separators) >= 1 {
				separator = separators[0]
			}
			messages := []string{}
			for _, err := range this.error.(validator.ValidationErrors) {
				if custom := this.errorMessages[err.StructField() + "." + err.Tag()]; len(custom) >= 1 {
					messages = append(messages, custom)
				} else {
					messages = append(messages, err.Translate(this.translator))
				}
			}
			message = strings.Join(messages, separator)
		}
		return message
	}
}

// 是否有错误
func (this *Model) HasError() bool {
	return this.error != nil
}

// 创建数据
func (this *Model) Create() bool {
	if !this.refValue.IsValid() {
		app.Panic("Model.refValue 无效")
	}
	if r := this.GetDb().Create(this.refValue.Interface()); r.Error == nil {
		return true
	} else {
		this.error = r.Error
		return false
	}
}

// 更新数据
func (this *Model) Update() bool {
	if !this.refValue.IsValid() {
		app.Panic("Model.refValue 无效")
	}
	if r := this.GetDb().Save(this.refValue.Interface()); r.Error == nil {
		return true
	} else {
		this.error = r.Error
		return false
	}
}

// 设置数据库连接
func (this *Model) SetDb(db *gorm.DB) {
	this.db = db
}

// 获取数据库连接
func (this *Model) GetDb() *gorm.DB {
	if this.db == nil {
		return app.Db
	} else {
		return this.db
	}
}