package str

import (
	"regexp"
	"strings"
)

// 移动到 toolbox/tool/str
// deprecated
func ToLowerCamel(str string) string {
	switch len(str) {
	case 0:
		return str
	case 1:
		return strings.ToLower(str)
	default:
		return strings.ToLower(str[0:1]) + str[1:]
	}
}

/*
// Sanitizes and formats a string to make an appropriate identifier in Go
function format(str)
{
    if (!str)
        return "";
    else if (str.match(/^\d+$/))
        str = "Num" + str;
    else if (str.charAt(0).match(/\d/))
    {
        const numbers = {'0': "Zero_", '1': "One_", '2': "Two_", '3': "Three_",
            '4': "Four_", '5': "Five_", '6': "Six_", '7': "Seven_",
            '8': "Eight_", '9': "Nine_"};
        str = numbers[str.charAt(0)] + str.substr(1);
    }
    return toProperCase(str).replace(/[^a-z0-9]/ig, "") || "NAMING_FAILED";
}
*/

var numExp = regexp.MustCompile("^\\d+$")
var singleNumExp = regexp.MustCompile("\\d")

var numbers = map[string]string{
	"0": "Zero_", "1": "One_", "2": "Two_", "3": "Three_",
	"4": "Four_", "5": "Five_", "6": "Six_", "7": "Seven_",
	"8": "Eight_", "9": "Nine_",
}

var noneWordExp = regexp.MustCompile("[^a-zA-Z0-9]")

// 移动到 toolbox/tool/str
// deprecated
func ToCamel(str string) string {
	if str == "" {
		return str
	}
	if numExp.MatchString(str) {
		str = "Num" + str
	}
	if singleNumExp.MatchString(str[0:1]) {
		str = numbers[str[0:1]] + str[1:]
	}
	upcase := toProperCase(str)
	return noneWordExp.ReplaceAllString(upcase, "")
}

// https://github.com/golang/lint/blob/5614ed5bae6fb75893070bdc0996a68765fdd275/lint.go#L771-L810
var commonInitialisms = map[string]bool{
	"ACL": true, "API": true, "ASCII": true, "CPU": true, "CSS": true, "DNS": true, "EOF": true, "GUID": true, "HTML": true, "HTTP": true,
	"HTTPS": true, "ID": true, "IP": true, "JSON": true, "LHS": true, "QPS": true, "RAM": true, "RHS": true, "RPC": true, "SLA": true,
	"SMTP": true, "SQL": true, "SSH": true, "TCP": true, "TLS": true, "TTL": true, "UDP": true, "UI": true, "UID": true, "UUID": true,
	"URI": true, "URL": true, "UTF8": true, "VM": true, "XML": true, "XMPP": true, "XSRF": true, "XSS": true,
}

var wordExp = regexp.MustCompile("(^|[^a-zA-Z])([a-z]+)")
var bigWordExp = regexp.MustCompile("([A-Z])([a-z]+)")

func toProperCase(str string) string {
	str = wordExp.ReplaceAllStringFunc(str, func(frag string) string {
		frag = strings.Replace(frag, "_", "", -1)
		if _, ok := commonInitialisms[strings.ToUpper(frag)]; ok {
			return strings.ToUpper(frag)
		} else if len(frag) == 1 {
			return strings.ToUpper(frag)
		} else {
			return strings.ToUpper(frag[0:1]) + strings.ToLower(frag[1:])
		}
	})
	return bigWordExp.ReplaceAllStringFunc(str, func(frag string) string {
		if _, ok := commonInitialisms[strings.ToUpper(frag)]; ok {
			return strings.ToUpper(frag)
		} else {
			return frag
		}
	})
}
