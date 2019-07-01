package abi_serializer

import (
	"reflect"

	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
)

type abiTraverseContext struct {
	maxSerializationTime common.Microseconds
	deadline             common.TimePoint
	recursionDepth       common.SizeT
}

func newAbiTraverseContext(maxSerializationTime common.Microseconds) abiTraverseContext {
	return abiTraverseContext{maxSerializationTime: maxSerializationTime, deadline: common.Now() + common.TimePoint(maxSerializationTime)}
}

func newAbiTraverseContextWithDeadline(maxSerializationTime common.Microseconds, deadline common.TimePoint) abiTraverseContext {
	return abiTraverseContext{maxSerializationTime: maxSerializationTime, deadline: deadline}
}

func outputName(s string, str string, shorten bool, maxLength common.SizeT) {
	minNumCharactersAtEnds := common.SizeT(4)
	preferredNumTailEndCharacters := common.SizeT(6)
	fillIn := "..."
	Assert(minNumCharactersAtEnds <= preferredNumTailEndCharacters, "preferred number of tail end characters cannot be less than the imposed absolute minimum")
	fillInLength := len(fillIn)
	minLength := fillInLength + 2*minNumCharactersAtEnds
	preferredMinLength := fillInLength + 2*preferredNumTailEndCharacters

	maxLength = common.SizeT(common.Max(uint64(maxLength), uint64(minLength)))
	if !shorten || len(str) <= maxLength {
		s += str
		return
	}

	actualNumTailEndCharacters := preferredNumTailEndCharacters
	if maxLength < preferredMinLength {
		actualNumTailEndCharacters = minNumCharactersAtEnds + (maxLength-minLength)/2
	}
	s += string([]byte(str)[:maxLength-fillInLength-actualNumTailEndCharacters])
	s += fillIn
	s += string([]byte(str + string(len(str)-actualNumTailEndCharacters))[:actualNumTailEndCharacters])
}

func (a abiTraverseContext) checkDeadline() {
	EosAssert(common.Now() < a.deadline, &AbiSerializationDeadlineException{}, "serialization time limit %v us exceeded", a.maxSerializationTime)
}

func (a abiTraverseContext) enterScope() func() {
	oldRecursionDepth := a.recursionDepth
	callBack := func() {
		a.recursionDepth = oldRecursionDepth
	}
	a.recursionDepth++
	EosAssert(a.recursionDepth < maxRecursionDepth, &AbiRecursionDepthException{}, "recursive definition, max_recursion_depth %d", maxRecursionDepth)
	a.checkDeadline()
	return callBack
}

type defMapItr struct {
	first  typeName
	second interface{}
}

type emptyPathRoot struct{}

type arrayTypePathRoot struct{}

type structTypePathRoot struct {
	structItr defMapItr
}

type variantTypePathRoot struct {
	variantItr defMapItr
}

type pathRoot = common.StaticVariant

type emptyPathItem struct{}

type arrayIndexPathItem struct {
	typeHint   pathRoot
	arrayIndex uint32
}

type fieldPathItem struct {
	parentStructItr defMapItr
	fieldOrdinal    uint32
}

type variantPathItem struct {
	variantItr   defMapItr
	fieldOrdinal uint32
}

type pathItem = common.StaticVariant

type abiTraverseContextWithPath struct {
	abiTraverseContext
	abis       *AbiSerializer
	rootOfPath pathRoot
	path       []pathItem
	ShortPath  bool
}

func newAbiTraverseContextWithPathByAbis(abis AbiSerializer, maxSerizalizationTime common.Microseconds, typeName typeName) abiTraverseContextWithPath {
	aw := abiTraverseContextWithPath{abiTraverseContext: newAbiTraverseContext(maxSerizalizationTime), abis: &abis}
	aw.setPathRoot(typeName)
	return aw
}

func newAbiTraverseContextWithPathByAbisDeadline(abis AbiSerializer, maxSerizalizationTime common.Microseconds, deadline common.TimePoint, typeName typeName) abiTraverseContextWithPath {
	aw := abiTraverseContextWithPath{abiTraverseContext: newAbiTraverseContextWithDeadline(maxSerizalizationTime, deadline), abis: &abis}
	aw.setPathRoot(typeName)
	return aw
}

func newAbiTraverseContextWithPathByCtx(abis AbiSerializer, ctx abiTraverseContext, typeName typeName) abiTraverseContextWithPath {
	aw := abiTraverseContextWithPath{abiTraverseContext: ctx, abis: &abis}
	aw.setPathRoot(typeName)
	return aw
}

func (aw abiTraverseContextWithPath) setPathRoot(typename typeName) {
	rtype := aw.abis.ResolveType(typename)
	if aw.abis.IsArray(rtype) {
		aw.rootOfPath = arrayTypePathRoot{}
	} else {
		structDef, ok := aw.abis.structs[rtype]
		if ok {
			aw.rootOfPath = structTypePathRoot{structItr: defMapItr{first: rtype, second: structDef}}
		} else {
			variantDef, ok := aw.abis.variants[rtype]
			if ok {
				aw.rootOfPath = variantTypePathRoot{variantItr: defMapItr{first: rtype, second: variantDef}}
			}
		}
	}
}

func (aw abiTraverseContextWithPath) pushToPath(item pathItem) func() {
	callBack := func() {
		EosAssert(len(aw.path) > 0, &AbiException{}, "invariant failure in variant_to_binary_context: path is empty on scope exit")
		aw.path = aw.path[:len(aw.path)-1]
	}
	aw.path = append(aw.path, item)
	return callBack
}

func (aw abiTraverseContextWithPath) setArrayIndexOfPathBack(i uint32) {
	EosAssert(len(aw.path) > 0, &AbiException{}, "path is empty")
	b := aw.path[len(aw.path)-1]
	EosAssert(reflect.TypeOf(b) == reflect.TypeOf(arrayIndexPathItem{}), &AbiException{}, "trying to set array index without first pushing new array index item")

	arrayItem := b.(arrayIndexPathItem)
	arrayItem.arrayIndex = i
	b = arrayItem
}

func (aw abiTraverseContextWithPath) hintArrayTypeIfInArray() {
	len := len(aw.path)
	if len == 0 || reflect.TypeOf(aw.path[len-1]) != reflect.TypeOf(arrayIndexPathItem{}) {
		return
	}

	arrayItem := aw.path[len-1].(arrayIndexPathItem)
	arrayItem.typeHint = arrayTypePathRoot{}
	aw.path[len-1] = arrayItem
}

func (aw abiTraverseContextWithPath) hintStructTypeIfInArray(itr defMapItr) {
	len := len(aw.path)
	if len == 0 || reflect.TypeOf(aw.path[len-1]) != reflect.TypeOf(arrayIndexPathItem{}) {
		return
	}

	arrayItem := aw.path[len-1].(arrayIndexPathItem)
	arrayItem.typeHint = structTypePathRoot{structItr: itr}
	aw.path[len-1] = arrayItem
}

func (aw abiTraverseContextWithPath) hintVariantTypeIfInArray(itr defMapItr) {
	len := len(aw.path)
	if len == 0 || reflect.TypeOf(aw.path[len-1]) != reflect.TypeOf(arrayIndexPathItem{}) {
		return
	}

	arrayItem := aw.path[len-1].(arrayIndexPathItem)
	arrayItem.typeHint = variantTypePathRoot{variantItr: itr}
	aw.path[len-1] = arrayItem
}

func (aw abiTraverseContextWithPath) maybeShorten(str string) string {
	if !aw.ShortPath {
		return str
	}
	s := ""
	outputName(s, str, true, common.SizeT(64))
	return s
}

type pathItemTypeVisitor struct {
	s            string
	shortenNames bool
}

func (visitor pathItemTypeVisitor) visit(item interface{}) {
	switch v := item.(type) {
	case emptyPathItem:
		visitor.visitEmptyPathItem(v)
	case arrayIndexPathItem:
		visitor.visitArrayIndexPathItem(v)
	case fieldPathItem:
		visitor.visitFieldPathItem(v)
	case variantPathItem:
		visitor.visitVariantPathItem(v)
	}
}

func (visitor pathItemTypeVisitor) visitEmptyPathItem(item emptyPathItem) {}

func (visitor pathItemTypeVisitor) visitArrayIndexPathItem(item arrayIndexPathItem) {
	th := item.typeHint
	if reflect.TypeOf(th) == reflect.TypeOf(structTypePathRoot{}) {
		str := th.(structTypePathRoot).structItr.first
		outputName(visitor.s, str, visitor.shortenNames, common.SizeT(64))
	} else if reflect.TypeOf(th) == reflect.TypeOf(variantTypePathRoot{}) {
		str := th.(variantTypePathRoot).variantItr.first
		outputName(visitor.s, str, visitor.shortenNames, common.SizeT(64))
	} else if reflect.TypeOf(th) == reflect.TypeOf(arrayTypePathRoot{}) {
		visitor.s += "ARRAY"
	} else {
		visitor.s += "UNKNOWN"
	}
}

func (visitor pathItemTypeVisitor) visitFieldPathItem(item fieldPathItem) {
	str := item.parentStructItr.second.(StructDef).Fields[item.fieldOrdinal].Type
	outputName(visitor.s, str, visitor.shortenNames, common.SizeT(64))
}

func (visitor pathItemTypeVisitor) visitVariantPathItem(item variantPathItem) {
	str := item.variantItr.second.(VariantDef).Types[item.fieldOrdinal]
	outputName(visitor.s, str, visitor.shortenNames, common.SizeT(64))
}
