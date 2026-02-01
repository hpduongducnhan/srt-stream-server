package app

// =======================
// COMMON
// =======================

type BaseResponse struct {
	Code    int    `json:"code"`
	Server  string `json:"server"`
	Service string `json:"service"`
	PID     string `json:"pid"`
}

type Kbps struct {
	Recv30s int `json:"recv_30s"`
	Send30s int `json:"send_30s"`
}

// =======================
// STREAM
// =======================

type PublishInfo struct {
	Active bool   `json:"active"`
	CID    string `json:"cid"`
}

type VideoInfo struct {
	Codec   string `json:"codec"`
	Profile string `json:"profile"`
	Level   string `json:"level"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

type AudioInfo struct {
	Codec      string `json:"codec"`
	SampleRate int    `json:"sample_rate"`
	Channel    int    `json:"channel"`
	Profile    string `json:"profile"`
}

type Stream struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Vhost     string      `json:"vhost"`
	App       string      `json:"app"`
	TcURL     string      `json:"tcUrl"`
	URL       string      `json:"url"`
	LiveMS    int64       `json:"live_ms"`
	Clients   int         `json:"clients"`
	Frames    int         `json:"frames"`
	SendBytes int64       `json:"send_bytes"`
	RecvBytes int64       `json:"recv_bytes"`
	Kbps      Kbps        `json:"kbps"`
	Publish   PublishInfo `json:"publish"`
	Video     *VideoInfo  `json:"video,omitempty"`
	Audio     *AudioInfo  `json:"audio,omitempty"`
}

// List streams
type StreamsResponse struct {
	BaseResponse
	Streams []Stream `json:"streams"`
}

// Single stream
type StreamResponse struct {
	BaseResponse
	Stream Stream `json:"stream"`
}

// =======================
// CLIENT
// =======================

type Client struct {
	ID        string  `json:"id"`
	Vhost     string  `json:"vhost"`
	Stream    string  `json:"stream"`
	IP        string  `json:"ip"`
	PageURL   string  `json:"pageUrl"`
	SwfURL    string  `json:"swfUrl"`
	TcURL     string  `json:"tcUrl"`
	URL       string  `json:"url"`
	Name      string  `json:"name"`
	Type      string  `json:"type"` // srt-publish | srt-play | rtmp-*
	Publish   bool    `json:"publish"`
	Alive     float64 `json:"alive"`
	SendBytes int64   `json:"send_bytes"`
	RecvBytes int64   `json:"recv_bytes"`
	Kbps      Kbps    `json:"kbps"`
}

// List clients
type ClientsResponse struct {
	BaseResponse
	Clients []Client `json:"clients"`
}

// Single client
type ClientResponse struct {
	BaseResponse
	Client Client `json:"client"`
}

// =======================
// SUMMARIES
// =======================

type SummarySelf struct {
	Version    string  `json:"version"`
	PID        int     `json:"pid"`
	PPID       int     `json:"ppid"`
	Argv       string  `json:"argv"`
	Cwd        string  `json:"cwd"`
	MemKB      int     `json:"mem_kbyte"`
	MemPercent float64 `json:"mem_percent"`
	CPUPercent float64 `json:"cpu_percent"`
	SrsUptime  int64   `json:"srs_uptime"`
}

type SummarySystem struct {
	CPUPercent      float64 `json:"cpu_percent"`
	DiskReadKBps    int64   `json:"disk_read_KBps"`
	DiskWriteKBps   int64   `json:"disk_write_KBps"`
	DiskBusyPercent float64 `json:"disk_busy_percent"`

	MemRamKB       int64   `json:"mem_ram_kbyte"`
	MemRamPercent  float64 `json:"mem_ram_percent"`
	MemSwapKB      int64   `json:"mem_swap_kbyte"`
	MemSwapPercent float64 `json:"mem_swap_percent"`

	CPUs       int `json:"cpus"`
	CPUsOnline int `json:"cpus_online"`

	Uptime   float64 `json:"uptime"`
	IdleTime float64 `json:"ilde_time"` // typo từ API – bạn xử lý đúng 👍

	Load1m  float64 `json:"load_1m"`
	Load5m  float64 `json:"load_5m"`
	Load15m float64 `json:"load_15m"`

	// Network (current sample)
	NetSampleTime int64 `json:"net_sample_time"`
	NetRecvBytes  int64 `json:"net_recv_bytes"`
	NetSendBytes  int64 `json:"net_send_bytes"`

	// Network (incremental / i)
	NetRecvIBytes int64 `json:"net_recvi_bytes"`
	NetSendIBytes int64 `json:"net_sendi_bytes"`

	// SRS traffic
	SrsSampleTime int64 `json:"srs_sample_time"`
	SrsRecvBytes  int64 `json:"srs_recv_bytes"`
	SrsSendBytes  int64 `json:"srs_send_bytes"`

	// Connections
	ConnSys    int `json:"conn_sys"`
	ConnSysET  int `json:"conn_sys_et"`
	ConnSysTW  int `json:"conn_sys_tw"`
	ConnSysUDP int `json:"conn_sys_udp"`
	ConnSrs    int `json:"conn_srs"`
}

type SummaryData struct {
	OK     bool          `json:"ok"`
	NowMS  int64         `json:"now_ms"`
	Self   SummarySelf   `json:"self"`
	System SummarySystem `json:"system"`
}

type SummaryResponse struct {
	BaseResponse
	Data SummaryData `json:"data"`
}
