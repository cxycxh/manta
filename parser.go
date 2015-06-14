package manta

import (
	"bytes"
	"io/ioutil"
	"math"
	"os"

	"github.com/dotabuff/manta/dota"
	"github.com/golang/snappy/snappy"
)

// The first 8 bytes of a replay for Source 1 and Source 2
var magicSource1 = []byte{'P', 'U', 'F', 'D', 'E', 'M', 'S', '\000'}
var magicSource2 = []byte{'P', 'B', 'D', 'E', 'M', 'S', '2', '\000'}

// A replay parser capable of parsing Source 2 replays
type Parser struct {
	Callbacks *Callbacks
	Tick      uint32

	hasClassInfo  bool
	classInfo     map[int32]string
	classIdSize   int
	classBaseline map[int32]map[string]interface{}

	sendTables   *sendTables
	stringTables *stringTables

	reader     *reader
	isStopping bool
}

// Create a new Parser from a file
func NewParserFromFile(path string) (*Parser, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	return NewParser(buf)
}

// Create a new parser from a byte slice
func NewParser(buf []byte) (*Parser, error) {
	// Create a new parser with an internal reader for the given buffer.
	parser := &Parser{
		Callbacks: &Callbacks{},
		Tick:      0,

		reader:     newReader(buf),
		isStopping: false,

		classInfo:     make(map[int32]string),
		classBaseline: make(map[int32]map[string]interface{}),
		stringTables:  newStringTables(),
	}

	// Parse out the header, ensuring that it's valid.
	if magic := parser.reader.readBytes(8); !bytes.Equal(magic, magicSource2) {
		return nil, _errorf("unexpected magic: expected %s, got %s", magicSource2, magic)
	}

	// Skip the next 8 bytes, which appear to be two int32s
	parser.reader.seekBytes(8)

	// Register callbacks

	// CDemoPacket outer messages have a inner handler
	parser.Callbacks.OnCDemoPacket(parser.onCDemoPacket)
	parser.Callbacks.OnCDemoSignonPacket(parser.onCDemoPacket)
	parser.Callbacks.OnCDemoFullPacket(parser.onCDemoFullPacket)

	// Packet entities, send tables and string tables are also low-level and
	// require internal handlers.
	parser.Callbacks.OnCSVCMsg_PacketEntities(parser.onCSVCMsg_PacketEntities)
	parser.Callbacks.OnCDemoSendTables(parser.onCDemoSendTables)
	parser.Callbacks.OnCDemoStringTables(parser.onCDemoStringTables)
	parser.Callbacks.OnCSVCMsg_CreateStringTable(parser.onCSVCMsg_CreateStringTable)
	parser.Callbacks.OnCSVCMsg_UpdateStringTable(parser.onCSVCMsg_UpdateStringTable)
	parser.Callbacks.OnCSVCMsg_SendTable(parser.onCSVCMsg_SendTable)

	parser.Callbacks.OnCSVCMsg_GameEvent(func(m *dota.CSVCMsg_GameEvent) error {
		_dump("gameevent", m)
		return nil
	})

	parser.Callbacks.OnCDemoClassInfo(parser.onCDemoClassInfo)

	parser.Callbacks.OnCSVCMsg_ServerInfo(func(m *dota.CSVCMsg_ServerInfo) error {
		parser.classIdSize = int(math.Log(float64(m.GetMaxClasses()))/math.Log(2)) + 1
		return nil
	})

	// Maintains the value of parser.Tick
	parser.Callbacks.OnCNETMsg_Tick(func(m *dota.CNETMsg_Tick) error {
		parser.Tick = m.GetTick()
		return nil
	})

	// Stops parsing when we reach the end of the replay.
	parser.Callbacks.OnCDemoStop(func(m *dota.CDemoStop) error {
		parser.Stop()
		return nil
	})

	// TODO
	parser.Callbacks.OnCDemoSpawnGroups(func(m *dota.CDemoSpawnGroups) error {
		return nil
	})

	// TODO
	parser.Callbacks.OnCNETMsg_SpawnGroup_Load(func(m *dota.CNETMsg_SpawnGroup_Load) error {
		return nil
	})

	// TODO
	parser.Callbacks.OnCDemoUserCmd(func(m *dota.CDemoUserCmd) error {
		return nil
	})

	return parser, nil
}

// Start parsing the replay. Will stop processing new events after Stop() is called.
func (p *Parser) Start() error {
	var msg Message
	var err error

	for !p.isStopping {
		if msg, err = p.read(); err != nil {
			return err
		}

		if err = p.CallByDemoType(int32(msg.Type), msg.data); err != nil {
			return err
		}
	}

	return nil
}

// Stop parsing the replay, causing the parser to stop processing new events.
func (p *Parser) Stop() {
	p.isStopping = true
}

// An outer message, right off the wire.
type Message struct {
	Compressed bool
	Tick       uint32
	Type       dota.EDemoCommands
	data       []byte
	Size       uint32
}

// Read the next outer message from the buffer.
func (p *Parser) read() (Message, error) {
	binType := p.reader.readVarUint32()
	binTick := p.reader.readVarUint32()
	binSize := p.reader.readVarUint32()

	msg := Message{
		Tick: binTick,
		Size: binSize,
	}

	command := dota.EDemoCommands(binType)
	msg.Compressed = (command & dota.EDemoCommands_DEM_IsCompressed) == dota.EDemoCommands_DEM_IsCompressed
	msg.Type = command & ^dota.EDemoCommands_DEM_IsCompressed

	buf := p.reader.readBytes(int(msg.Size))

	if msg.Compressed {
		decodedLen, err := snappy.DecodedLen(buf)
		if err != nil {
			return msg, err
		}

		if decodedLen > 0x100000 {
			return msg, _errorf("decompressed size too big: %d", decodedLen)
		}

		out, err := snappy.Decode(nil, buf)
		if err != nil {
			return msg, err
		}
		msg.data = out
		msg.Size = uint32(decodedLen)
	} else {
		msg.data = buf
	}

	return msg, nil
}
