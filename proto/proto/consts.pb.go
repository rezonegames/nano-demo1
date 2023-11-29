// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.4
// source: consts.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 该结构与consts结构一样，客户端服务器共用，只要定义就不能改变
type AccountType int32

const (
	AccountType_DEVICEID AccountType = 0
	AccountType_WX       AccountType = 1
	AccountType_FB       AccountType = 2
	AccountType_GIT      AccountType = 3
)

// Enum value maps for AccountType.
var (
	AccountType_name = map[int32]string{
		0: "DEVICEID",
		1: "WX",
		2: "FB",
		3: "GIT",
	}
	AccountType_value = map[string]int32{
		"DEVICEID": 0,
		"WX":       1,
		"FB":       2,
		"GIT":      3,
	}
)

func (x AccountType) Enum() *AccountType {
	p := new(AccountType)
	*p = x
	return p
}

func (x AccountType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccountType) Descriptor() protoreflect.EnumDescriptor {
	return file_consts_proto_enumTypes[0].Descriptor()
}

func (AccountType) Type() protoreflect.EnumType {
	return &file_consts_proto_enumTypes[0]
}

func (x AccountType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccountType.Descriptor instead.
func (AccountType) EnumDescriptor() ([]byte, []int) {
	return file_consts_proto_rawDescGZIP(), []int{0}
}

// 暂时这样，以后拆出来，游戏内状态和游戏外状态 todo：
type GameState int32

const (
	//  在房间里
	GameState_IDLE GameState = 0
	GameState_WAIT GameState = 1
	//  已分到桌子
	GameState_INGAME GameState = 2
	GameState_LOGOUT GameState = 3
)

// Enum value maps for GameState.
var (
	GameState_name = map[int32]string{
		0: "IDLE",
		1: "WAIT",
		2: "INGAME",
		3: "LOGOUT",
	}
	GameState_value = map[string]int32{
		"IDLE":   0,
		"WAIT":   1,
		"INGAME": 2,
		"LOGOUT": 3,
	}
)

func (x GameState) Enum() *GameState {
	p := new(GameState)
	*p = x
	return p
}

func (x GameState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GameState) Descriptor() protoreflect.EnumDescriptor {
	return file_consts_proto_enumTypes[1].Descriptor()
}

func (GameState) Type() protoreflect.EnumType {
	return &file_consts_proto_enumTypes[1]
}

func (x GameState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GameState.Descriptor instead.
func (GameState) EnumDescriptor() ([]byte, []int) {
	return file_consts_proto_rawDescGZIP(), []int{1}
}

type TableState int32

const (
	TableState_STATE_NONE TableState = 0
	TableState_WAITREADY  TableState = 1
	TableState_CANCEL     TableState = 2
	TableState_CHECK_RES  TableState = 3
	TableState_GAMING     TableState = 4
	TableState_SETTLEMENT TableState = 5
)

// Enum value maps for TableState.
var (
	TableState_name = map[int32]string{
		0: "STATE_NONE",
		1: "WAITREADY",
		2: "CANCEL",
		3: "CHECK_RES",
		4: "GAMING",
		5: "SETTLEMENT",
	}
	TableState_value = map[string]int32{
		"STATE_NONE": 0,
		"WAITREADY":  1,
		"CANCEL":     2,
		"CHECK_RES":  3,
		"GAMING":     4,
		"SETTLEMENT": 5,
	}
)

func (x TableState) Enum() *TableState {
	p := new(TableState)
	*p = x
	return p
}

func (x TableState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TableState) Descriptor() protoreflect.EnumDescriptor {
	return file_consts_proto_enumTypes[2].Descriptor()
}

func (TableState) Type() protoreflect.EnumType {
	return &file_consts_proto_enumTypes[2]
}

func (x TableState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TableState.Descriptor instead.
func (TableState) EnumDescriptor() ([]byte, []int) {
	return file_consts_proto_rawDescGZIP(), []int{2}
}

type RoomType int32

const (
	RoomType_ROOMTYPE_NONE RoomType = 0
	RoomType_QUICK         RoomType = 1
	RoomType_MATCH         RoomType = 2
)

// Enum value maps for RoomType.
var (
	RoomType_name = map[int32]string{
		0: "ROOMTYPE_NONE",
		1: "QUICK",
		2: "MATCH",
	}
	RoomType_value = map[string]int32{
		"ROOMTYPE_NONE": 0,
		"QUICK":         1,
		"MATCH":         2,
	}
)

func (x RoomType) Enum() *RoomType {
	p := new(RoomType)
	*p = x
	return p
}

func (x RoomType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RoomType) Descriptor() protoreflect.EnumDescriptor {
	return file_consts_proto_enumTypes[3].Descriptor()
}

func (RoomType) Type() protoreflect.EnumType {
	return &file_consts_proto_enumTypes[3]
}

func (x RoomType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RoomType.Descriptor instead.
func (RoomType) EnumDescriptor() ([]byte, []int) {
	return file_consts_proto_rawDescGZIP(), []int{3}
}

type TableType int32

const (
	TableType_TABLETYPE_NONE TableType = 0
	TableType_NORMAL         TableType = 1
	TableType_HAPPY          TableType = 2
)

// Enum value maps for TableType.
var (
	TableType_name = map[int32]string{
		0: "TABLETYPE_NONE",
		1: "NORMAL",
		2: "HAPPY",
	}
	TableType_value = map[string]int32{
		"TABLETYPE_NONE": 0,
		"NORMAL":         1,
		"HAPPY":          2,
	}
)

func (x TableType) Enum() *TableType {
	p := new(TableType)
	*p = x
	return p
}

func (x TableType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TableType) Descriptor() protoreflect.EnumDescriptor {
	return file_consts_proto_enumTypes[4].Descriptor()
}

func (TableType) Type() protoreflect.EnumType {
	return &file_consts_proto_enumTypes[4]
}

func (x TableType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TableType.Descriptor instead.
func (TableType) EnumDescriptor() ([]byte, []int) {
	return file_consts_proto_rawDescGZIP(), []int{4}
}

type ActionType int32

const (
	ActionType_ACTIONTYPE_NONE ActionType = 0
	ActionType_FRAME_ONE       ActionType = 1
	ActionType_MOVE            ActionType = 2
	ActionType_ROTATE          ActionType = 3
	ActionType_DROP            ActionType = 4
	ActionType_QUICK_DROP      ActionType = 5
	// 连击
	ActionType_COMBO   ActionType = 8
	ActionType_COMBO_4 ActionType = 9
	ActionType_COMBO_3 ActionType = 10
	// 道具
	ActionType_ITEM_BOOM         ActionType = 11
	ActionType_ITEM_BUFF_DISTURB ActionType = 12
	ActionType_ITEM_ADD_ROW      ActionType = 13
	ActionType_ITEM_DEL_ROW      ActionType = 14
	ActionType_END               ActionType = 100
)

// Enum value maps for ActionType.
var (
	ActionType_name = map[int32]string{
		0:   "ACTIONTYPE_NONE",
		1:   "FRAME_ONE",
		2:   "MOVE",
		3:   "ROTATE",
		4:   "DROP",
		5:   "QUICK_DROP",
		8:   "COMBO",
		9:   "COMBO_4",
		10:  "COMBO_3",
		11:  "ITEM_BOOM",
		12:  "ITEM_BUFF_DISTURB",
		13:  "ITEM_ADD_ROW",
		14:  "ITEM_DEL_ROW",
		100: "END",
	}
	ActionType_value = map[string]int32{
		"ACTIONTYPE_NONE":   0,
		"FRAME_ONE":         1,
		"MOVE":              2,
		"ROTATE":            3,
		"DROP":              4,
		"QUICK_DROP":        5,
		"COMBO":             8,
		"COMBO_4":           9,
		"COMBO_3":           10,
		"ITEM_BOOM":         11,
		"ITEM_BUFF_DISTURB": 12,
		"ITEM_ADD_ROW":      13,
		"ITEM_DEL_ROW":      14,
		"END":               100,
	}
)

func (x ActionType) Enum() *ActionType {
	p := new(ActionType)
	*p = x
	return p
}

func (x ActionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ActionType) Descriptor() protoreflect.EnumDescriptor {
	return file_consts_proto_enumTypes[5].Descriptor()
}

func (ActionType) Type() protoreflect.EnumType {
	return &file_consts_proto_enumTypes[5]
}

func (x ActionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ActionType.Descriptor instead.
func (ActionType) EnumDescriptor() ([]byte, []int) {
	return file_consts_proto_rawDescGZIP(), []int{5}
}

var File_consts_proto protoreflect.FileDescriptor

var file_consts_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x34, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x45, 0x56, 0x49, 0x43, 0x45, 0x49, 0x44,
	0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x57, 0x58, 0x10, 0x01, 0x12, 0x06, 0x0a, 0x02, 0x46, 0x42,
	0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x47, 0x49, 0x54, 0x10, 0x03, 0x2a, 0x37, 0x0a, 0x09, 0x47,
	0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x44, 0x4c, 0x45,
	0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x57, 0x41, 0x49, 0x54, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06,
	0x49, 0x4e, 0x47, 0x41, 0x4d, 0x45, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x4f, 0x47, 0x4f,
	0x55, 0x54, 0x10, 0x03, 0x2a, 0x62, 0x0a, 0x0a, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x4e, 0x4f, 0x4e, 0x45,
	0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x57, 0x41, 0x49, 0x54, 0x52, 0x45, 0x41, 0x44, 0x59, 0x10,
	0x01, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x10, 0x02, 0x12, 0x0d, 0x0a,
	0x09, 0x43, 0x48, 0x45, 0x43, 0x4b, 0x5f, 0x52, 0x45, 0x53, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06,
	0x47, 0x41, 0x4d, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x45, 0x54, 0x54,
	0x4c, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x05, 0x2a, 0x33, 0x0a, 0x08, 0x52, 0x6f, 0x6f, 0x6d,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x4f, 0x4f, 0x4d, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x51, 0x55, 0x49, 0x43, 0x4b,
	0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x10, 0x02, 0x2a, 0x36, 0x0a,
	0x09, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x54, 0x41,
	0x42, 0x4c, 0x45, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x41,
	0x50, 0x50, 0x59, 0x10, 0x02, 0x2a, 0xd8, 0x01, 0x0a, 0x0a, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x13, 0x0a, 0x0f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x46, 0x52, 0x41,
	0x4d, 0x45, 0x5f, 0x4f, 0x4e, 0x45, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4d, 0x4f, 0x56, 0x45,
	0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x4f, 0x54, 0x41, 0x54, 0x45, 0x10, 0x03, 0x12, 0x08,
	0x0a, 0x04, 0x44, 0x52, 0x4f, 0x50, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x51, 0x55, 0x49, 0x43,
	0x4b, 0x5f, 0x44, 0x52, 0x4f, 0x50, 0x10, 0x05, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x4f, 0x4d, 0x42,
	0x4f, 0x10, 0x08, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4d, 0x42, 0x4f, 0x5f, 0x34, 0x10, 0x09,
	0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4d, 0x42, 0x4f, 0x5f, 0x33, 0x10, 0x0a, 0x12, 0x0d, 0x0a,
	0x09, 0x49, 0x54, 0x45, 0x4d, 0x5f, 0x42, 0x4f, 0x4f, 0x4d, 0x10, 0x0b, 0x12, 0x15, 0x0a, 0x11,
	0x49, 0x54, 0x45, 0x4d, 0x5f, 0x42, 0x55, 0x46, 0x46, 0x5f, 0x44, 0x49, 0x53, 0x54, 0x55, 0x52,
	0x42, 0x10, 0x0c, 0x12, 0x10, 0x0a, 0x0c, 0x49, 0x54, 0x45, 0x4d, 0x5f, 0x41, 0x44, 0x44, 0x5f,
	0x52, 0x4f, 0x57, 0x10, 0x0d, 0x12, 0x10, 0x0a, 0x0c, 0x49, 0x54, 0x45, 0x4d, 0x5f, 0x44, 0x45,
	0x4c, 0x5f, 0x52, 0x4f, 0x57, 0x10, 0x0e, 0x12, 0x07, 0x0a, 0x03, 0x45, 0x4e, 0x44, 0x10, 0x64,
	0x42, 0x08, 0x5a, 0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_consts_proto_rawDescOnce sync.Once
	file_consts_proto_rawDescData = file_consts_proto_rawDesc
)

func file_consts_proto_rawDescGZIP() []byte {
	file_consts_proto_rawDescOnce.Do(func() {
		file_consts_proto_rawDescData = protoimpl.X.CompressGZIP(file_consts_proto_rawDescData)
	})
	return file_consts_proto_rawDescData
}

var file_consts_proto_enumTypes = make([]protoimpl.EnumInfo, 6)
var file_consts_proto_goTypes = []interface{}{
	(AccountType)(0), // 0: proto.AccountType
	(GameState)(0),   // 1: proto.GameState
	(TableState)(0),  // 2: proto.TableState
	(RoomType)(0),    // 3: proto.RoomType
	(TableType)(0),   // 4: proto.TableType
	(ActionType)(0),  // 5: proto.ActionType
}
var file_consts_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_consts_proto_init() }
func file_consts_proto_init() {
	if File_consts_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_consts_proto_rawDesc,
			NumEnums:      6,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_consts_proto_goTypes,
		DependencyIndexes: file_consts_proto_depIdxs,
		EnumInfos:         file_consts_proto_enumTypes,
	}.Build()
	File_consts_proto = out.File
	file_consts_proto_rawDesc = nil
	file_consts_proto_goTypes = nil
	file_consts_proto_depIdxs = nil
}
