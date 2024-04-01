package merkledag

import (
	"encoding/json"
	"strings"
)

const STEP = 4 // 步长常量，用于解析二进制数据的步长

// Hash2File 将哈希值映射到文件
// 示例路径：/doc/tmp/temp.txt
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 检查存储中是否存在对应哈希的对象
	flag, _ := store.Has(hash)
	if flag {
		// 获取对应哈希值的二进制对象
		objBinary, _ := store.Get(hash)
		obj := binaryToObj(objBinary)
		// 根据路径获取文件内容并返回
		pathArr := strings.Split(path, "\\")
		cur := 1
		return getFileByDir(obj, pathArr, cur, store)
	}
	return nil
}

// getFileByDir 根据目录获取文件
func getFileByDir(obj *Object, pathArr []string, cur int, store KVStore) []byte {
	if cur >= len(pathArr) {
		return nil
	}
	index := 0
	for i := range obj.Links {
		// 解析对象类型
		objType := string(obj.Data[index : index+STEP])
		index += STEP
		objInfo := obj.Links[i]
		if objInfo.Name != pathArr[cur] {
			continue
		}
		// 根据不同类型获取文件内容
		switch objType {
		case TREE:
			objDirBinary, _ := store.Get(objInfo.Hash)
			objDir := binaryToObj(objDirBinary)
			ans := getFileByDir(objDir, pathArr, cur+1, store)
			if ans != nil {
				return ans
			}
		case BLOB:
			ans, _ := store.Get(objInfo.Hash)
			return ans
		case LIST:
			objLinkBinary, _ := store.Get(objInfo.Hash)
			objList := binaryToObj(objLinkBinary)
			ans := getFileByList(objList, store)
			return ans
		}
	}
	return nil
}

// binaryToObj 将二进制数据解析为对象
func binaryToObj(objBinary []byte) *Object {
	var res Object
	json.Unmarshal(objBinary, &res)
	return &res
}

// getFileByList 根据列表获取文件
func getFileByList(obj *Object, store KVStore) []byte {
	ans := make([]byte, 0)
	index := 0
	for i := range obj.Links {
		curObjType := string(obj.Data[index : index+STEP])
		index += STEP
		curObjLink := obj.Links[i]
		curObjBinary, _ := store.Get(curObjLink.Hash)
		curObj := binaryToObj(curObjBinary)
		// 根据不同类型获取文件内容
		if curObjType == BLOB {
			ans = append(ans, curObjBinary...)
		} else { //List
			tmp := getFileByList(curObj, store)
			ans = append(ans, tmp...)
		}
	}
	return ans
}
