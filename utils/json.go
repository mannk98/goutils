package utils

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/antchfx/jsonquery"
	xj "github.com/basgys/goxml2json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func JsonGetPathNode(node *jsonquery.Node) (retstr string) {
	if node.Type == jsonquery.TextNode {
		retstr = "/@" + retstr
	} else {
		retstr = "/" + retstr
	}
	tmpnode := node
	for {
		if tmpnode = tmpnode.Parent; tmpnode == nil {
			retstr = strings.Replace(retstr, "//", "/", 1)
			break
		} else {
			retstr = "/" + tmpnode.Data + retstr
		}
	}

	return retstr
}

func JsonSet(jsonstring string, elementPath string, val any) (string, error) {
	return sjson.Set(jsonstring, elementPath, val)
}

func JsonStringFindElements(strjson *string, pathSearch string, valueflag ...bool) (map[string]string, error) {
	var retmap = map[string]string{}
	doc, err := jsonquery.Parse(strings.NewReader(*strjson))
	if err != nil {
		//		fmt.Println("xmlquery.Parse:", err)
		return retmap, err
	}

	nodes, err := jsonquery.QueryAll(doc, pathSearch)
	if err != nil {
		return retmap, err
	}
	if len(nodes) == 0 {
		return retmap, errors.New("missing keypath")
	}
	//found := false
	//fmt.Println("scan nodes", err)
	id := 0
	//numnodes := len(nodes)
	for k := 0; k < len(nodes); k++ {
		node := nodes[k]
		key := JsonGetPathNode(node)
		//fmt.Println("key:", key, v)
		exist := func() bool {
			for i := 0; i < k; i++ {
				if key == JsonGetPathNode(nodes[i]) {
					return true
				}
			}
			return false
		}()

		//found = true
		//fmt.Println("xmlStringFindElement:", v.NamespaceURI, v.InnerText())
		if exist {
			key = key + "[" + strconv.Itoa(id) + "]"
			id = id + 1
		}
		//|| v.Type == xmlquery.AttributeNode
		if node.FirstChild == nil || node.FirstChild.FirstChild == nil {
			retmap[key] = node.InnerText()
		} else {
			xml := strings.NewReader(node.OutputXML())
			json, err := xj.Convert(xml)
			if err != nil {
				retmap[key] = node.OutputXML()
			} else {
				retmap[key] = json.String()
			}
			//retmap[key] = v.OutputXML(true)
		}
		//retmap[strconv.Itoa(id)] = v.InnerText()
	}

	return retmap, nil
}

func JsonStringFindElementsSlide(strjson *string, pathSearch string) ([]string, error) {
	var retslide = []string{}
	doc, err := jsonquery.Parse(strings.NewReader(*strjson))
	if err != nil {
		//		fmt.Println("xmlquery.Parse:", err)
		return retslide, err
	}

	nodes, err := jsonquery.QueryAll(doc, pathSearch)
	if err != nil {
		return retslide, err
	}
	//	numnodes := len(nodes)
	for k := 0; k < len(nodes); k++ {
		v := nodes[k]

		if v.FirstChild == nil || v.FirstChild.FirstChild == nil {
			retslide = append(retslide, v.InnerText())
		} else {
			xml := strings.NewReader(v.OutputXML())
			json, err := xj.Convert(xml)
			if err != nil {
				retslide = append(retslide, v.OutputXML())
			} else {
				retslide = append(retslide, json.String())
			}
			//			retmap[key] = v.OutputXML(true)
		}
		//		retmap[strconv.Itoa(id)] = v.InnerText()
	}
	return retslide, nil
}

func JsonParser(str string) gjson.Result {
	return gjson.Parse(str)
}

func JsonParserBytes(bs []byte) gjson.Result {
	return gjson.ParseBytes(bs)
}

func JsonStringFindElement(strjson *string, pathSearch string) (string, error) {
	if retmap, err := JsonStringFindElements(strjson, pathSearch); err == nil && len(retmap) != 0 {
		keys := make([]string, 0, len(retmap))
		for k := range retmap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return retmap[keys[0]], nil
	} else {
		return "", err
	}
}
