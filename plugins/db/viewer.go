package db

import (
	"m3game/meta/metapb"
	"m3game/plugins/log"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	_maxviewnum = 64
)

type Viewer[TM proto.Message, TV Flag] struct {
	root *ViewNode // View树
	mu   sync.Mutex
}

type ViewNode struct {
	field  protoreflect.FieldDescriptor // 节点反射
	wflags [_maxviewnum]bool            // 白名单标记
	bflags [_maxviewnum]bool            // 黑名单标记
	views  [_maxviewnum][]int           // 视图索引
	subs   []*ViewNode                  // 子节点
}

func NewViewer[TM proto.Message, TV Flag]() *Viewer[TM, TV] {
	viewer := &Viewer[TM, TV]{}
	// 收集所有View
	var tv TV
	vflags := make(map[string]int)
	if tv.Descriptor().Values().Len() > _maxviewnum {
		log.Fatal("ViewFlag num %d > maxviewnum %d", tv.Descriptor().Values().Len(), _maxviewnum)
		return nil
	}
	for i := 0; i < tv.Descriptor().Values().Len(); i++ {
		name := tv.Descriptor().Values().Get(i).Name()
		number := tv.Descriptor().Values().Get(i).Number()
		if int(number) >= _maxviewnum {
			log.Fatal("ViewFlag name %s number %d >= maxviewnum %d", string(name), int(number), _maxviewnum)
			return nil
		}
		vflags[string(name)] = int(number)
	}
	// 构建View树
	root := &ViewNode{}
	var tm TM
	for i := 0; i < tm.ProtoReflect().Descriptor().Fields().Len(); i++ {
		field := tm.ProtoReflect().Descriptor().Fields().Get(i)
		vn := BuildViewTree(vflags, field)
		root.subs = append(root.subs, vn)
	}
	// 构建视图
	for _, vflag := range vflags {
		TagViewTree(vflag, root)
	}
	viewer.root = root
	return viewer
}

func BuildViewTree(vflags map[string]int, field protoreflect.FieldDescriptor) *ViewNode {
	vn := &ViewNode{
		field: field,
	}
	// 收集本节点View信息
	if v := proto.GetExtension(field.Options(), metapb.E_ViewfieldOption); !reflect.ValueOf(v).IsNil() {
		viewfieldoption := v.(*metapb.M3ViewFieldOption)
		for _, f := range strings.Split(viewfieldoption.Bflag, ",") {
			if f == "" {
				continue
			}
			if _, ok := vflags[f]; !ok {
				log.Fatal("Field %s has invaild Bflag [%s]", field.Name(), f)
				return nil
			}
			vn.bflags[vflags[f]] = true
		}
		for _, f := range strings.Split(viewfieldoption.Wflag, ",") {
			if f == "" {
				continue
			}
			if _, ok := vflags[f]; !ok {
				log.Fatal("Field %s has invaild Wflag [%s]", field.Name(), f)
				return nil
			}
			vn.wflags[vflags[f]] = true
		}
	}
	// message类型, 收集子节点View信息
	if field.Kind() == protoreflect.MessageKind {
		for i := 0; i < field.Message().Fields().Len(); i++ {
			subfield := field.Message().Fields().Get(i)
			subvn := BuildViewTree(vflags, subfield)
			vn.subs = append(vn.subs, subvn)
		}
	}
	return vn
}

// 标记视图
func TagViewTree(vflag int, vn *ViewNode) bool {
	// 持有B标记
	if vn.bflags[vflag] {
		return false
	}
	// 持有W标记
	if vn.wflags[vflag] {
		return true
	}
	// 检查子节点
	for idex, subvn := range vn.subs {
		if TagViewTree(vflag, subvn) {
			vn.views[vflag] = append(vn.views[vflag], idex)
		}
	}
	// 子节点中存在W标记
	if len(vn.views[vflag]) > 0 {
		return true
	}
	return false
}

func FillViewObj(vflag int, vn *ViewNode, obj proto.Message) proto.Message {
	// 源数据为空
	if reflect.ValueOf(obj).IsNil() {
		return nil
	}
	// 节点为终端节点,全部返回
	if len(vn.views[vflag]) == 0 {
		return obj
	}

	refType := reflect.TypeOf(obj).Elem()
	value := reflect.New(refType)
	nobj := value.Interface().(proto.Message)
	log.Debug("11  [%v] [%v] %v", obj, nobj, obj == nil)
	// 继续检索子节点
	for _, idex := range vn.views[vflag] {
		subvn := vn.subs[idex]
		// 子节点是message 继续检索
		if subvn.field.Kind() == protoreflect.MessageKind {
			if subobj := FillViewObj(vflag, subvn, obj.ProtoReflect().Get(subvn.field).Message().Interface()); subobj != nil {
				log.Debug("22 %v %v", subobj, nobj)
				nobj.ProtoReflect().Set(subvn.field, protoreflect.ValueOfMessage(subobj.ProtoReflect()))
			}
			continue
		}
		// 其他类型 直接拷贝
		nobj.ProtoReflect().Set(subvn.field, obj.ProtoReflect().Get(subvn.field))
	}
	return nobj
}

func (vm *Viewer[TM, TV]) Filter(vflag TV, obj TM) TM {
	log.Debug("Filter %v", vflag)
	nobj := FillViewObj(int(vflag.Number()), vm.root, obj)
	if nobj == nil {
		var tm TM
		return tm
	}
	return nobj.(TM)
}
