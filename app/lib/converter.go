package lib

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ConvertToMD5 func
func ConvertToMD5(value *int) string {
	str := IntToStr(*value)
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ConvertStrToMD5 func
func ConvertStrToMD5(value *string) string {
	var str string = *value
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ConvertToSHA1 func
func ConvertToSHA1(value string) string {
	sha := sha1.New()
	sha.Write([]byte(value))
	encrypted := sha.Sum(nil)
	encryptedString := fmt.Sprintf("%x", encrypted)
	return encryptedString
}

// ConvertToSHA256 func
func ConvertToSHA256(value string) string {
	hash := sha256.Sum256([]byte(value))
	res := fmt.Sprintf("%x", hash)
	return res
}

// IntToStr func
func IntToStr(value int) string {
	return strconv.Itoa(value)
}

// StrToInt func
func StrToInt(value string) int {
	valueInt, _ := strconv.Atoi(value)
	return valueInt
}

// StrToInt64 func
func StrToInt64(value string) int64 {
	valueInt, _ := strconv.ParseInt(value, 10, 64)
	return valueInt
}

// StrToFloat func
func StrToFloat(value string) float64 {
	valueFloat, _ := strconv.ParseFloat(value, 64)
	return valueFloat
}

// StrToBool func
func StrToBool(value string) bool {
	valueBool, _ := strconv.ParseBool(value)
	return valueBool
}

// FloatToStr func
func FloatToStr(inputNum float64, prec ...int) string {
	precision := 0 // Default precision is 0 if not specified
	if len(prec) > 0 {
		precision = prec[0]
	}
	return strconv.FormatFloat(inputNum, 'f', precision, 64)
}

// ConvertJsonToStr func
func ConvertJsonToStr(payload interface{}) string {
	jsonData, _ := JSONMarshal(payload)
	return string(jsonData)
}

// ConvertStrToObj func
func ConvertStrToObj(value string) map[string]interface{} {
	var output map[string]interface{}
	JSONUnmarshal([]byte(value), &output)
	return output
}

// ConvertStrToArrObj func
func ConvertStrToArrObj(value string) []map[string]interface{} {
	var output []map[string]interface{}
	JSONUnmarshal([]byte(value), &output)
	return output
}

// ConvertStrToJson func
func ConvertStrToJson(value string) interface{} {
	var output interface{}
	JSONUnmarshal([]byte(value), &output)
	return output
}

// ConvertStrToTime func
func ConvertStrToTime(value string) *time.Time {
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, value)
	return &t
}

// ConvertSliceIntToStr func
// Source: https://stackoverflow.com/a/37533144
func ConvertSliceIntToStr(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

// ConvertSliceStrToStr func
// Source: https://stackoverflow.com/a/37533144
func ConvertSliceStrToStr(a []string, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

// ConvertSliceUUIDToStr func
// Source: https://stackoverflow.com/a/37533144
func ConvertSliceUUIDToStr(a []uuid.UUID, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

// StrLeadingZerosRemove func
func StrLeadingZerosRemove(str string) (result string) {
	if len(str) > 0 {
		result = strings.TrimLeft(str, "0")
	}
	return
}
