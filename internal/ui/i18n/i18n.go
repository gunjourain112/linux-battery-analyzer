package i18n

import "strings"

type Key string

const (
	AppTitle               Key = "app_title"
	ChooseLanguage         Key = "choose_language"
	ChooseRange            Key = "choose_range"
	LanguageKo             Key = "language_ko"
	LanguageEn             Key = "language_en"
	RangeLast7             Key = "range_last_7"
	RangeLast30            Key = "range_last_30"
	RangeCustom            Key = "range_custom"
	SincePrompt            Key = "since_prompt"
	UntilPrompt            Key = "until_prompt"
	EnterToContinue        Key = "enter_to_continue"
	EnterToSelect          Key = "enter_to_select"
	EnterToRunTabToSwitch  Key = "enter_to_run_tab_to_switch"
	InvalidSinceDate       Key = "invalid_since_date"
	InvalidUntilDate       Key = "invalid_until_date"
	LanguageLabel          Key = "language_label"
	SinceLabel             Key = "since_label"
	UntilLabel             Key = "until_label"
	ReportSummary          Key = "report_summary"
	DeviceSpecsSection     Key = "device_specs_section"
	BatteryHealthSection   Key = "battery_health_section"
	AnalysisSummarySection Key = "analysis_summary_section"
	AnalysisPeriodHeader   Key = "analysis_period_header"
	ActualUseHeader        Key = "actual_use_header"
	BatteryStateHeader     Key = "battery_state_header"
	AvgLoadHeader          Key = "avg_load_header"
	TempRangeHeader        Key = "temp_range_header"
	ExpectedRemainHeader   Key = "expected_remain_header"
	ReportSessions         Key = "report_sessions"
	ReportDaily            Key = "report_daily"
	ReportCharging         Key = "report_charging"
	ReportDischargeProfile Key = "report_discharge_profile"
	ReportBatteryHealth    Key = "report_battery_health"
	ReportScenarioEstimate Key = "report_scenario_estimate"
	ReportProcessSummary   Key = "report_process_summary"
	ReportOptimizationTips Key = "report_optimization_tips"
	ReportInsightDashboard Key = "report_insight_dashboard"
	ReportAIContext        Key = "report_ai_context"
	ReportProcessImpacts   Key = "report_process_impacts"
	ReportSystemEvents     Key = "report_system_events"
	ReportSpecs            Key = "report_specs"
	ReportThermals         Key = "report_thermals"
	ReportDischargeTrend   Key = "report_discharge_trend"
	ReportThermalTimeline  Key = "report_thermal_timeline"
	BatteryGraph           Key = "battery_graph"
	UnifiedTimeline        Key = "unified_timeline"
	NoSessions             Key = "no_sessions"
	NoDailyRecords         Key = "no_daily_records"
	NoChargingSessions     Key = "no_charging_sessions"
	NoDischargeProfile     Key = "no_discharge_profile"
	NoBatteryHealthData    Key = "no_battery_health_data"
	NoScenarioEstimateData Key = "no_scenario_estimate_data"
	NoTimelineData         Key = "no_timeline_data"
	NoProcessSummaryData   Key = "no_process_summary_data"
	NoOptimizationTips     Key = "no_optimization_tips"
	NoInsightDashboardData Key = "no_insight_dashboard_data"
	NoAIContextData        Key = "no_ai_context_data"
	NoProcessImpactData    Key = "no_process_impact_data"
	NoSystemEvents         Key = "no_system_events"
	NoHardwareSpecs        Key = "no_hardware_specs"
	NoThermalSamples       Key = "no_thermal_samples"
	ScenarioEstimateNote   Key = "scenario_estimate_note"
	TopDrainLabel          Key = "top_drain_label"
	HeavyLoadLabel         Key = "heavy_load_label"
	PeakTempLabel          Key = "peak_temp_label"
	ProcessSamplesLabel    Key = "process_samples_label"
	TipTopProcess          Key = "tip_top_process"
	TipHeavyLoad           Key = "tip_heavy_load"
	TipThermal             Key = "tip_thermal"
	TipDrainRate           Key = "tip_drain_rate"
	TipNoObviousIssue      Key = "tip_no_obvious_issue"
	AIContextLead          Key = "ai_context_lead"
	AIContextAsk           Key = "ai_context_ask"
	AvgDischarge           Key = "avg_discharge"
	WorstSession           Key = "worst_session"
	CurrentStatusLabel     Key = "current_status_label"
	RangeSummaryLabel      Key = "range_summary_label"
	ActualUseLabel         Key = "actual_use_label"
	TopDrainSessionLabel   Key = "top_drain_session_label"
	StartHeader            Key = "start_header"
	EndHeader              Key = "end_header"
	DurationHeader         Key = "duration_header"
	DrainHeader            Key = "drain_header"
	RateHeader             Key = "rate_header"
	DateHeader             Key = "date_header"
	ActiveHeader           Key = "active_header"
	ChargeHeader           Key = "charge_header"
	AvgWHeader             Key = "avg_w_header"
	PeakWHeader            Key = "peak_w_header"
	DrainWHeader           Key = "drain_w_header"
	BucketHeader           Key = "bucket_header"
	CountHeader            Key = "count_header"
	RatioHeader            Key = "ratio_header"
	ProcessHeader          Key = "process_header"
	LevelHeader            Key = "level_header"
	CPUHeader              Key = "cpu_header"
	MemHeader              Key = "mem_header"
	TimeHeader             Key = "time_header"
	HeatHeader             Key = "heat_header"
	PowerHeader            Key = "power_header"
	EventHeader            Key = "event_header"
	StateHeader            Key = "state_header"
	TypeHeader             Key = "type_header"
	DescriptionHeader      Key = "description_header"
	SamplesHeader          Key = "samples_header"
	MetricHeader           Key = "metric_header"
	ValueHeader            Key = "value_header"
	CurrentHeader          Key = "current_header"
	FullHeader             Key = "full_header"
	DesignCapacityHeader   Key = "design_capacity_header"
	CurrentCapacityHeader  Key = "current_capacity_header"
	HealthHeader           Key = "health_header"
	MinHeader              Key = "min_header"
	MaxHeader              Key = "max_header"
	AvgHeader              Key = "avg_header"
	OSHeader               Key = "os_header"
	DeviceHeader           Key = "device_header"
	RAMHeader              Key = "ram_header"
	BatteryHeader          Key = "battery_header"
	LightLevel             Key = "light_level"
	MediumLevel            Key = "medium_level"
	HeavyLevel             Key = "heavy_level"
	UnknownLevel           Key = "unknown_level"
)

type Translator struct {
	lang string
}

func New(lang string) Translator {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "en":
		return Translator{lang: "en"}
	default:
		return Translator{lang: "ko"}
	}
}

func (t Translator) Get(key Key) string {
	if t.lang == "en" {
		if v, ok := english[key]; ok {
			return v
		}
	}
	if v, ok := korean[key]; ok {
		return v
	}
	return string(key)
}

var english = map[Key]string{
	AppTitle:               "Notebook Battery Analyzer",
	ChooseLanguage:         "Choose language",
	ChooseRange:            "Choose period",
	LanguageKo:             "Korean",
	LanguageEn:             "English",
	RangeLast7:             "Last 7 days",
	RangeLast30:            "Last 30 days",
	RangeCustom:            "Custom range",
	SincePrompt:            "since: ",
	UntilPrompt:            "until: ",
	EnterToContinue:        "enter to continue",
	EnterToSelect:          "enter to select",
	EnterToRunTabToSwitch:  "enter to run, tab to switch",
	InvalidSinceDate:       "invalid since date",
	InvalidUntilDate:       "invalid until date",
	LanguageLabel:          "language",
	SinceLabel:             "since",
	UntilLabel:             "until",
	ReportSummary:          "1. Device & Health Info",
	DeviceSpecsSection:     "Device Specs",
	BatteryHealthSection:   "Battery Health",
	AnalysisSummarySection: "Analysis Summary",
	AnalysisPeriodHeader:   "Analysis Period",
	ActualUseHeader:        "Actual use",
	BatteryStateHeader:     "Battery state",
	AvgLoadHeader:          "Avg load",
	TempRangeHeader:        "Temp range",
	ExpectedRemainHeader:   "Expected remain",
	ReportSessions:         "Sessions",
	ReportDaily:            "Daily",
	ReportCharging:         "Charging",
	ReportDischargeProfile: "Discharge Profile",
	ReportBatteryHealth:    "Battery Health",
	ReportScenarioEstimate: "Scenario Estimate",
	ReportProcessSummary:   "Process Summary",
	ReportOptimizationTips: "Optimization Tips",
	ReportInsightDashboard: "Insight Dashboard",
	ReportAIContext:        "AI Context",
	ReportProcessImpacts:   "Process Impacts",
	ReportSystemEvents:     "System Events",
	ReportSpecs:            "Specs",
	ReportThermals:         "Thermals",
	ReportDischargeTrend:   "Battery Discharge Trend",
	ReportThermalTimeline:  "Hourly Thermal Trend",
	BatteryGraph:           "Battery Graph",
	UnifiedTimeline:        "Unified Timeline",
	NoSessions:             "no sessions",
	NoDailyRecords:         "no daily records",
	NoChargingSessions:     "no charging sessions",
	NoDischargeProfile:     "no discharge profile",
	NoBatteryHealthData:    "no battery health data",
	NoScenarioEstimateData: "no scenario estimate data",
	NoTimelineData:         "no timeline data",
	NoProcessSummaryData:   "no process summary data",
	NoOptimizationTips:     "no optimization tips",
	NoInsightDashboardData: "no insight dashboard data",
	NoAIContextData:        "no ai context data",
	NoProcessImpactData:    "no process impact data",
	NoSystemEvents:         "no system events",
	NoHardwareSpecs:        "no hardware specs",
	NoThermalSamples:       "no thermal samples",
	ScenarioEstimateNote:   "estimate based on observed battery drain over this window, not design capacity",
	TopDrainLabel:          "Top Drain",
	HeavyLoadLabel:         "Heavy Load",
	PeakTempLabel:          "Peak Temp",
	ProcessSamplesLabel:    "Process Samples",
	TipTopProcess:          "Top drain process: %s (%.1fW)",
	TipHeavyLoad:           "Heavy-load bucket is %.0f%% of observed discharge; trim background work",
	TipThermal:             "Peak temperature reached %d°C; sustained heat usually tracks higher drain",
	TipDrainRate:           "Average session drain is %.2f%%/h; power saver or a lighter workload helps",
	TipNoObviousIssue:      "No obvious optimization issue stands out in this window",
	AIContextLead:          "Paste this context into an AI chat to continue analysis.",
	AIContextAsk:           "Focus on drain causes, runtime scenarios, and concrete battery-saving actions.",
	AvgDischarge:           "avg discharge",
	WorstSession:           "worst session",
	CurrentStatusLabel:     "Current",
	RangeSummaryLabel:      "Range",
	ActualUseLabel:         "Actual use",
	TopDrainSessionLabel:   "Top drain session",
	StartHeader:            "Start",
	EndHeader:              "End",
	DurationHeader:         "Duration",
	DrainHeader:            "Drain",
	RateHeader:             "Rate",
	DateHeader:             "Date",
	ActiveHeader:           "Active",
	ChargeHeader:           "Charge",
	AvgWHeader:             "Avg W",
	PeakWHeader:            "Peak W",
	DrainWHeader:           "Drain W",
	BucketHeader:           "Bucket",
	CountHeader:            "Count",
	RatioHeader:            "Ratio",
	ProcessHeader:          "Process",
	LevelHeader:            "Level",
	CPUHeader:              "CPU s",
	MemHeader:              "Mem M",
	TimeHeader:             "Time",
	HeatHeader:             "Temp",
	PowerHeader:            "Power",
	EventHeader:            "Event",
	StateHeader:            "State",
	TypeHeader:             "Type",
	DescriptionHeader:      "Description",
	SamplesHeader:          "Samples",
	MetricHeader:           "Metric",
	ValueHeader:            "Value",
	CurrentHeader:          "Current est.",
	FullHeader:             "Full est.",
	DesignCapacityHeader:   "Design capacity",
	CurrentCapacityHeader:  "Current capacity",
	HealthHeader:           "Health",
	MinHeader:              "Min",
	MaxHeader:              "Max",
	AvgHeader:              "Avg",
	OSHeader:               "OS",
	DeviceHeader:           "Device",
	RAMHeader:              "RAM",
	BatteryHeader:          "Battery",
	LightLevel:             "light",
	MediumLevel:            "medium",
	HeavyLevel:             "heavy",
	UnknownLevel:           "unknown",
}

var korean = map[Key]string{
	AppTitle:               "Notebook Battery Analyzer",
	ChooseLanguage:         "언어 선택",
	ChooseRange:            "기간 선택",
	LanguageKo:             "한국어",
	LanguageEn:             "영어",
	RangeLast7:             "최근 7일",
	RangeLast30:            "최근 30일",
	RangeCustom:            "직접 입력",
	SincePrompt:            "시작일: ",
	UntilPrompt:            "종료일: ",
	EnterToContinue:        "Enter: 다음",
	EnterToSelect:          "Enter: 선택",
	EnterToRunTabToSwitch:  "Enter: 실행, Tab: 전환",
	InvalidSinceDate:       "시작일이 올바르지 않습니다",
	InvalidUntilDate:       "종료일이 올바르지 않습니다",
	LanguageLabel:          "언어",
	SinceLabel:             "시작",
	UntilLabel:             "종료",
	ReportSummary:          "1. 기기 및 헬스 정보",
	DeviceSpecsSection:     "기기 사양",
	BatteryHealthSection:   "배터리 헬스",
	AnalysisSummarySection: "분석 세션 요약",
	AnalysisPeriodHeader:   "분석 기간",
	ActualUseHeader:        "실제 사용",
	BatteryStateHeader:     "배터리 상태",
	AvgLoadHeader:          "평균 방전",
	TempRangeHeader:        "온도 범위",
	ExpectedRemainHeader:   "예상 잔여 시간",
	ReportSessions:         "세션",
	ReportDaily:            "일별",
	ReportCharging:         "충전",
	ReportDischargeProfile: "방전 프로파일",
	ReportBatteryHealth:    "배터리 상태",
	ReportScenarioEstimate: "시나리오 추정",
	ReportProcessSummary:   "프로세스 요약",
	ReportOptimizationTips: "최적화 팁",
	ReportInsightDashboard: "인사이트 대시보드",
	ReportAIContext:        "AI 컨텍스트",
	ReportProcessImpacts:   "프로세스 영향",
	ReportSystemEvents:     "시스템 이벤트",
	ReportSpecs:            "사양",
	ReportThermals:         "온도",
	ReportDischargeTrend:   "배터리 방전 추이",
	ReportThermalTimeline:  "시간대별 온도 추이",
	BatteryGraph:           "배터리 그래프",
	UnifiedTimeline:        "통합 타임라인",
	NoSessions:             "세션 없음",
	NoDailyRecords:         "일별 기록 없음",
	NoChargingSessions:     "충전 세션 없음",
	NoDischargeProfile:     "방전 프로파일 없음",
	NoBatteryHealthData:    "배터리 상태 데이터 없음",
	NoScenarioEstimateData: "시나리오 추정 데이터 없음",
	NoTimelineData:         "타임라인 데이터 없음",
	NoProcessSummaryData:   "프로세스 요약 데이터 없음",
	NoOptimizationTips:     "최적화 팁 없음",
	NoInsightDashboardData: "인사이트 대시보드 데이터 없음",
	NoAIContextData:        "AI 컨텍스트 데이터 없음",
	NoProcessImpactData:    "프로세스 영향 데이터 없음",
	NoSystemEvents:         "시스템 이벤트 없음",
	NoHardwareSpecs:        "하드웨어 사양 없음",
	NoThermalSamples:       "온도 샘플 없음",
	ScenarioEstimateNote:   "이 추정은 이 구간의 배터리 방전 추세를 기준으로 계산했으며, 설계 용량은 사용하지 않습니다",
	TopDrainLabel:          "상위 방전",
	HeavyLoadLabel:         "고부하",
	PeakTempLabel:          "최고 온도",
	ProcessSamplesLabel:    "프로세스 샘플",
	TipTopProcess:          "가장 많이 방전시킨 프로세스: %s (%.1fW)",
	TipHeavyLoad:           "고부하 구간이 관측 방전의 %.0f%%입니다. 백그라운드 작업을 줄여보세요",
	TipThermal:             "최고 온도가 %d°C입니다. 지속 발열은 방전 증가와 자주 같이 옵니다",
	TipDrainRate:           "평균 세션 방전율이 %.2f%%/h입니다. 절전 모드나 작업량 조절이 도움이 됩니다",
	TipNoObviousIssue:      "이 구간에서는 뚜렷한 최적화 포인트가 보이지 않습니다",
	AIContextLead:          "이 컨텍스트를 AI 채팅에 붙여넣어 추가 분석을 진행하세요.",
	AIContextAsk:           "방전 원인, 사용 시나리오, 배터리 절감 행동을 중심으로 봐주세요.",
	AvgDischarge:           "평균 방전",
	WorstSession:           "가장 빠른 방전 세션",
	CurrentStatusLabel:     "현재",
	RangeSummaryLabel:      "범위",
	ActualUseLabel:         "실제 사용",
	TopDrainSessionLabel:   "최고 소모 세션",
	StartHeader:            "시작",
	EndHeader:              "종료",
	DurationHeader:         "지속시간",
	DrainHeader:            "방전",
	RateHeader:             "비율",
	DateHeader:             "날짜",
	ActiveHeader:           "활동",
	ChargeHeader:           "충전",
	AvgWHeader:             "평균 W",
	PeakWHeader:            "최대 W",
	DrainWHeader:           "방전 W",
	BucketHeader:           "구간",
	CountHeader:            "개수",
	RatioHeader:            "비율",
	ProcessHeader:          "프로세스",
	LevelHeader:            "레벨",
	CPUHeader:              "CPU 초",
	MemHeader:              "메모리 M",
	TimeHeader:             "시간",
	HeatHeader:             "온도",
	PowerHeader:            "전력",
	EventHeader:            "이벤트",
	StateHeader:            "상태",
	TypeHeader:             "유형",
	DescriptionHeader:      "설명",
	SamplesHeader:          "샘플",
	MetricHeader:           "항목",
	ValueHeader:            "값",
	CurrentHeader:          "현재 추정",
	FullHeader:             "완충 추정",
	DesignCapacityHeader:   "설계 용량",
	CurrentCapacityHeader:  "현재 용량",
	HealthHeader:           "헬스",
	MinHeader:              "최소",
	MaxHeader:              "최대",
	AvgHeader:              "평균",
	OSHeader:               "OS",
	DeviceHeader:           "장치",
	RAMHeader:              "RAM",
	BatteryHeader:          "배터리",
	LightLevel:             "낮음",
	MediumLevel:            "보통",
	HeavyLevel:             "높음",
	UnknownLevel:           "알 수 없음",
}
