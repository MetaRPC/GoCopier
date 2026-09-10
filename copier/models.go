package copier

type Account struct {
	Type     string `json:"type"`
	User     uint64 `json:"user"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Name     string `json:"name"`
}

type StartRequest struct {
	UserKey            string   `json:"user_key"`
	ManagerKey         string   `json:"manager_key"`
	Master             *Account `json:"master"`
	Slave              *Account `json:"slave"`
	RiskType           string   `json:"risk_type"`
	RiskValue          string   `json:"risk_value"`
	FixedMasterBalance string   `json:"fixed_master_balance"`
	CopySl             bool     `json:"copy_sl"`
	CopyTp             bool     `json:"copy_tp"`
	CopyPendingOrders  bool     `json:"copy_pending_orders"`
	ReverseCopy        bool     `json:"reverse_copy"`
}

type StartReply struct {
	Ok       bool   `json:"ok"`
	CopierId string `json:"copier_id"`
	Error    string `json:"error"`
}

type CopierSummary struct {
	Id          string `json:"id"`
	MasterType  string `json:"master_type"`
	MasterUser  uint64 `json:"master_user"`
	MasterServer string `json:"master_server"`
	SlaveType   string `json:"slave_type"`
	SlaveUser   uint64 `json:"slave_user"`
	SlaveServer string `json:"slave_server"`
	RiskType    string `json:"risk_type"`
	RiskValue   string `json:"risk_value"`
	Paused      bool   `json:"paused"`
	PauseReason string `json:"pause_reason"`
}

type ListReply struct {
	Ok      bool            `json:"ok"`
	Copiers []*CopierSummary `json:"copiers"`
	Error   string          `json:"error"`
}

type SimpleReply struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type GuiDemoOpenAccountRequest struct {
	Company        string `json:"company"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Server         string `json:"server"`
	TimeoutSeconds int32  `json:"timeout_seconds"`
}

type GuiDemoOpenAccountReply struct {
	ResultCode int32  `json:"result_code"`
	Login      uint64 `json:"login"`
	Password   string `json:"password"`
	Investor   string `json:"investor"`
	Server     string `json:"server"`
}
