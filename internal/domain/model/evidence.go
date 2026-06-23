package model

// SourceLevel 是情报信源等级：S/A/B 可作正式证据，C 只能作背景材料。
type SourceLevel string

// EvidenceRole 区分正式证据与背景材料。
type EvidenceRole string

// EventType 标记证据所对应的事件类型，重大事件需要更严格的多源验证。
type EventType string

const (
	SourceLevelS SourceLevel = "S"
	SourceLevelA SourceLevel = "A"
	SourceLevelB SourceLevel = "B"
	SourceLevelC SourceLevel = "C"

	EvidenceFormal     EvidenceRole = "formal"
	EvidenceBackground EvidenceRole = "background"

	EventNormal        EventType = "normal"
	EventMajorPositive EventType = "major_positive"
	EventMajorNegative EventType = "major_negative"
	EventBuyLogicBreak EventType = "buy_logic_break"
)

func (v SourceLevel) Valid() bool {
	return valid(v, SourceLevelS, SourceLevelA, SourceLevelB, SourceLevelC)
}
func (v EvidenceRole) Valid() bool { return valid(v, EvidenceFormal, EvidenceBackground) }

// FormalAllowed 返回该信源等级是否可进入正式裁决证据链。
func (v SourceLevel) FormalAllowed() bool {
	return v == SourceLevelS || v == SourceLevelA || v == SourceLevelB
}

// HighGrade 返回该信源是否属于重大事件所需的高等级信源。
func (v SourceLevel) HighGrade() bool { return v == SourceLevelS || v == SourceLevelA }

// Evidence 是规则引擎裁决时使用的证据摘要。
type Evidence struct {
	EvidenceID                      string
	SummaryID                       string
	SourceName                      string
	PublishedAt                     string
	CapturedAt                      string
	OriginalURL                     string
	Summary                         string
	ContentHash                     string
	ChunkHash                       string
	TimeWeight                      float64
	RelevanceScore                  float64
	SourceLevel                     SourceLevel
	Role                            EvidenceRole
	EventType                       EventType
	IndependentSourceCount          int
	HighGradeIndependentSourceCount int
}

// EvidenceSet 是一次工作流聚合后的证据集合与验证状态。
type EvidenceSet struct {
	Items              []Evidence
	VerificationStatus VerificationStatus
}
