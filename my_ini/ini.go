package my_ini

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

func LoadIni(fileName string, data interface{}) (err error) {
	t := reflect.TypeOf(data)
	//判断是不是一个指针
	if t.Kind() != reflect.Ptr {
		err = errors.New("data should be a pointer") //格式话输出以后返回一个error类型
		return
	}
	//判断是不是一个结构体类型
	if t.Elem().Kind() != reflect.Struct {
		err = errors.New("data param should be a struct pointer")
		return
	}
	//读取配置文件
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	lineSlince := strings.Split(string(b), "\r\n")
	var structName string
	for idx, line := range lineSlince {
		//去掉字符串首位的空格
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		//如果注释就跳过
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			//如果是[开头就表示是节（section）
			if line[0] != '[' || line[len(line)-1] != ']' {
				err = fmt.Errorf("line:%d syntax error : ", err)
				return
			}
			//把这一行尾的[]去掉，取到中间的内容把首位的空格去掉拿到内容
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			if len(sectionName) == 0 {
				err = fmt.Errorf("line:%d syntax error :", idx+1)
				return
			}
			//根据字符串sectionName去data里面根据反射找到对应的结构体
			for i := 0; i < t.Elem().NumField(); i++ {
				filed := t.Elem().Field(i)
				if sectionName == filed.Tag.Get("ini") {
					//说明找到了对应的嵌套结构体
					structName = filed.Name
					//fmt.Printf("找到%s对应的嵌套结构体是%s\n", sectionName, structName)
				}
			}
		} else {
			//以等号分割这一行，等号左边是key，右边是value
			if strings.Index(line, "=") == -1 || strings.HasPrefix(line, "=") {
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			index := strings.Index(line, "=")
			key := strings.TrimSpace(line[:index])
			value := strings.TrimSpace(line[index+1:])
			//根据structname 去 data 里面把对应的嵌套结构体给取出来
			v := reflect.ValueOf(data)
			sValue := v.Elem().FieldByName(structName)
			sType := sValue.Type()
			if sType.Kind() != reflect.Struct {
				fmt.Println("data is %s is pointer", structName)
				return
			}
			var fileName string
			var fileType reflect.StructField

			//遍历嵌套结构体的每一个字段，判读tag是不是等于key
			for i := 0; i < sValue.NumField(); i++ {
				filed := sType.Field(i)
				fileType = filed
				if filed.Tag.Get("ini") == key {
					//找到对应的字段
					fileName = filed.Name
					break
				}
			}
			if len(fileName) == 0 {
				continue
			}
			fileObj := sValue.FieldByName(fileName)
			//对其赋值
			//fmt.Println(fileName, fileType.Type.Kind())
			switch fileType.Type.Kind() {
			case reflect.String:
				fileObj.SetString(value)
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
				var valueInt int64
				valueInt, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					err = fmt.Errorf("line：%d value syntax type err", idx+1)
					return
				}
				fileObj.SetInt(valueInt)
			case reflect.Float64, reflect.Float32:
				var valueFloat float64
				valueFloat, err = strconv.ParseFloat(value, 64)
				if err != nil {
					err = fmt.Errorf("line:%d value type error", idx+1)
					return
				}
				fileObj.SetFloat(valueFloat)
			}
		}
	}
	return
}
