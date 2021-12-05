package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type token struct {
	value     string
	tokenType string
}

type node struct {
	childrens       map[string]*node
	childrensArray  []*node
	values          map[string]string
	valuesPrimitive []string
	nodeType        string
}

func lex(str string) []token {
	var res []token
	var t string

	//next
	j := 0
	isString := false
	var next = func() string {
		for i := j; i < len(str); i++ {
			if str[i] == ' ' && isString {
				j = i + 1
				return string(str[i])
			}
			if str[i] == ' ' || str[i] == '\n' || str[i] == '\t' || str[i] == '\r' {
				if i == len(str)-1 {
					j = i + 1
				}
				continue
			} else {
				j = i + 1
				return string(str[i])
			}
		}
		return ""
	}
	var nextnext = func(val string) string {
		for i := j; i < len(str); i++ {
			if str[i] == ' ' || str[i] == '"' || str[i] == '\n' || str[i] == '\t' || str[i] == '\r' {
				continue
			} else {
				if string(str[i]) == val {
					j = i + 1
					return string(str[i])
				} else {
					return ""
				}
			}
		}
		return ""
	}
	//

	for j <= len(str)-1 {
		cChar := next()
		if cChar == "\"" {
			isString = !isString
			continue
		}
		if cChar == "{" {
			res = append(res, token{
				value:     t,
				tokenType: "openObj",
			})
			t = ""
		} else if cChar == ":" && nextnext("{") == "{" {
			res = append(res, token{
				value:     t,
				tokenType: "open",
			})
			t = ""
		} else if cChar == ":" && nextnext("[") == "[" {
			res = append(res, token{
				value:     t,
				tokenType: "openArray",
			})
			t = ""
		} else if cChar == "," || cChar == "}" || cChar == "]" {
			if len(t) > 0 {
				var isPrimitive bool = !strings.Contains(t, ":")
				if isPrimitive {
					res = append(res, token{
						value:     t,
						tokenType: "primitive",
					})
				} else {
					res = append(res, token{
						value:     t,
						tokenType: "value",
					})
				}
				t = ""
			}
			if cChar == "}" || cChar == "]" {
				res = append(res, token{
					value:     string(cChar),
					tokenType: "closed",
				})
			}
		} else {
			t = t + string(cChar)
		}
	}
	return res
}

func serialize(cNode node) string {
	res := ""
	var deepNode func(cNode node, key string, isFinal bool)
	deepNode = func(cNode node, key string, isFinal bool) {
		count := 0
		if len(key) > 0 {
			res = res + "\"" + key + "\"" + ":"
		}
		if cNode.nodeType == "object" {
			res = res + "{"
		} else {
			res = res + "["
		}
		if cNode.nodeType == "object" {
			count = 0
			for key, item := range cNode.values {
				count = count + 1
				res = res + "\"" + key + "\"" + ":" + "\"" + item + "\""
				if count != len(cNode.values) || len(cNode.childrens) != 0 {
					res = res + ","
				}
			}
			count = 0
			for key, item := range cNode.childrens {
				count = count + 1
				deepNode(*item, key, count == len(cNode.childrens))
			}
		} else {
			for i, item := range cNode.valuesPrimitive {
				res = res + item
				if i != len(cNode.valuesPrimitive)-1 || len(cNode.childrensArray) != 0 {
					res = res + ","
				}
			}
			for i, item := range cNode.childrensArray {
				deepNode(*item, "", i == len(cNode.childrensArray)-1)
			}
		}
		if cNode.nodeType == "object" {
			res = res + "}"
		} else {
			res = res + "]"
		}
		if !isFinal {
			res = res + ","
		}
	}
	deepNode(cNode, "", true)
	return res
}

func main() {
	b, err := ioutil.ReadFile("input.json")
	if err != nil {
		fmt.Print(err)
	}
	var stack []*node
	var tObj *node
	var tokens []token = lex(string(b))

	//main
	var initToken string = tokens[0].tokenType
	tokens = tokens[1:]
	if initToken == "openObj" {
		var cNode = node{
			childrens: make(map[string]*node),
			nodeType:  "object",
		}
		stack = append(stack, &cNode)
	} else {
		//todo
	}

	for _, cToken := range tokens {
		if cToken.tokenType == "open" || cToken.tokenType == "openObj" || cToken.tokenType == "openArray" {
			var cNode = node{
				childrens: make(map[string]*node),
				values:    make(map[string]string),
			}
			tObj = stack[len(stack)-1]
			stack = append(stack, &cNode)
			if cToken.tokenType == "openArray" {
				cNode.nodeType = "array"
			} else {
				cNode.nodeType = "object"
			}
			if cToken.tokenType == "openObj" {
				tObj.childrensArray = append(tObj.childrensArray, &cNode)
			} else {
				tObj.childrens[cToken.value] = &cNode
			}
		} else if cToken.tokenType == "value" {
			var c []string = strings.Split(cToken.value, ":")
			if c[1] == "fizz" {
				c[1] = "buzz"
			} else if c[1] == "buzz" {
				c[1] = "fizz"
			}
			tObj = stack[len(stack)-1]
			tObj.values[c[0]] = c[1]
		} else if cToken.tokenType == "primitive" {
			tObj = stack[len(stack)-1]
			tObj.valuesPrimitive = append(tObj.valuesPrimitive, cToken.value)
		} else if cToken.tokenType == "closed" && len(stack) > 1 {
			stack = stack[:len(stack)-1]
		}
	}
	fmt.Print(serialize(*stack[0]))
}
