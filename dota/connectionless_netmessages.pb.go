// Code generated by protoc-gen-go.
// source: connectionless_netmessages.proto
// DO NOT EDIT!

package dota

import proto "github.com/golang/protobuf/proto"
import math "math"

// discarding unused import google_protobuf "github.com/dotabuff/manta/dota/google/protobuf"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type ConnectionLessMessageIds int32

const (
	ConnectionLessMessageIds_C2S_CONNECT_id ConnectionLessMessageIds = 80
)

var ConnectionLessMessageIds_name = map[int32]string{
	80: "C2S_CONNECT_id",
}
var ConnectionLessMessageIds_value = map[string]int32{
	"C2S_CONNECT_id": 80,
}

func (x ConnectionLessMessageIds) Enum() *ConnectionLessMessageIds {
	p := new(ConnectionLessMessageIds)
	*p = x
	return p
}
func (x ConnectionLessMessageIds) String() string {
	return proto.EnumName(ConnectionLessMessageIds_name, int32(x))
}
func (x *ConnectionLessMessageIds) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ConnectionLessMessageIds_value, data, "ConnectionLessMessageIds")
	if err != nil {
		return err
	}
	*x = ConnectionLessMessageIds(value)
	return nil
}

type C2S_CONNECT_Message struct {
	HostVersion       *uint32                       `protobuf:"varint,1,opt,name=host_version" json:"host_version,omitempty"`
	AuthProtocol      *uint32                       `protobuf:"varint,2,opt,name=auth_protocol" json:"auth_protocol,omitempty"`
	ChallengeNumber   *uint32                       `protobuf:"varint,3,opt,name=challenge_number" json:"challenge_number,omitempty"`
	ReservationCookie *uint64                       `protobuf:"fixed64,4,opt,name=reservation_cookie" json:"reservation_cookie,omitempty"`
	LowViolence       *bool                         `protobuf:"varint,5,opt,name=low_violence" json:"low_violence,omitempty"`
	EncryptedPassword []byte                        `protobuf:"bytes,6,opt,name=encrypted_password" json:"encrypted_password,omitempty"`
	Splitplayers      []*CCLCMsg_SplitPlayerConnect `protobuf:"bytes,7,rep,name=splitplayers" json:"splitplayers,omitempty"`
	AuthSteam         []byte                        `protobuf:"bytes,8,opt,name=auth_steam" json:"auth_steam,omitempty"`
	XXX_unrecognized  []byte                        `json:"-"`
}

func (m *C2S_CONNECT_Message) Reset()         { *m = C2S_CONNECT_Message{} }
func (m *C2S_CONNECT_Message) String() string { return proto.CompactTextString(m) }
func (*C2S_CONNECT_Message) ProtoMessage()    {}

func (m *C2S_CONNECT_Message) GetHostVersion() uint32 {
	if m != nil && m.HostVersion != nil {
		return *m.HostVersion
	}
	return 0
}

func (m *C2S_CONNECT_Message) GetAuthProtocol() uint32 {
	if m != nil && m.AuthProtocol != nil {
		return *m.AuthProtocol
	}
	return 0
}

func (m *C2S_CONNECT_Message) GetChallengeNumber() uint32 {
	if m != nil && m.ChallengeNumber != nil {
		return *m.ChallengeNumber
	}
	return 0
}

func (m *C2S_CONNECT_Message) GetReservationCookie() uint64 {
	if m != nil && m.ReservationCookie != nil {
		return *m.ReservationCookie
	}
	return 0
}

func (m *C2S_CONNECT_Message) GetLowViolence() bool {
	if m != nil && m.LowViolence != nil {
		return *m.LowViolence
	}
	return false
}

func (m *C2S_CONNECT_Message) GetEncryptedPassword() []byte {
	if m != nil {
		return m.EncryptedPassword
	}
	return nil
}

func (m *C2S_CONNECT_Message) GetSplitplayers() []*CCLCMsg_SplitPlayerConnect {
	if m != nil {
		return m.Splitplayers
	}
	return nil
}

func (m *C2S_CONNECT_Message) GetAuthSteam() []byte {
	if m != nil {
		return m.AuthSteam
	}
	return nil
}

func init() {
	proto.RegisterEnum("dota.ConnectionLessMessageIds", ConnectionLessMessageIds_name, ConnectionLessMessageIds_value)
}