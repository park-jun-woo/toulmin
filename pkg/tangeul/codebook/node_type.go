//ff:type feature=tangl type=model
//ff:what NodeType — the ten TANGL Node-packet entity types (codebook §2.1)
package codebook

// NodeType is the entity-type code carried in word 3 of a Node packet
// (TANGEUL §2.1, §5.1). These are Free-area L1 codebook values. The codebook is
// the single source of these code values and is FROZEN append-only: values are
// never reassigned or reused, and a retired code is marked deprecated rather
// than deleted.
type NodeType uint16

// The ten TANGL NodeType codes (§2.1). Assigned 0x0001..0x000A; 0x0000 is
// reserved as "unset" and is never a valid NodeType.
const (
	// NodeNamespace is `this document is x` (0x0001).
	NodeNamespace NodeType = 0x0001
	// NodeCase is `in case of x` (0x0002).
	NodeCase NodeType = 0x0002
	// NodeRule is `x is a general/counter/except rule` (0x0003).
	NodeRule NodeType = 0x0003
	// NodeFunc is an external function reference nodeified, e.g. bank.withdraw (0x0004).
	NodeFunc NodeType = 0x0004
	// NodeEndpoint is `provides x` (0x0005).
	NodeEndpoint NodeType = 0x0005
	// NodeTrigger is `on`/`every`/`until` (0x0006).
	NodeTrigger NodeType = 0x0006
	// NodeTerm is a Definitions `means`/`has` entry (0x0007).
	NodeTerm NodeType = 0x0007
	// NodeField is `x is required [as Type]` (0x0008).
	NodeField NodeType = 0x0008
	// NodeAlias is `see a from path` (0x0009).
	NodeAlias NodeType = 0x0009
	// NodeGroup is an order-preserving execution sequence (`1.`) (0x000A).
	NodeGroup NodeType = 0x000A
)
