package ipTable

import (
	_ "fmt"
	_ "net"
	"utils/logging"
)

/*
#cgo CFLAGS: -I../../../netfilter/libiptables/include -I../../../netfilter/iptables/include
#cgo LDFLAGS: -L../../../netfilter/libiptables/lib -lip4tc
#include "ipTable.h"
*/
import "C"

const (
	ALL_RULE_STR = "all"

	// Error Messages
	INSERTING_RULE_ERROR = "adding ip rule to iptables failed: "
	DELETING_RULE_ERROR  = "deleting rule failed: "
)

type SysdIpTableHandler struct {
	logger   *logging.Writer
	ruleInfo map[string]C.ipt_config_t
}
