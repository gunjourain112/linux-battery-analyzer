package i18n

import "strings"

type Key string

const (
	AppTitle               Key = "app_title"
	ChooseLanguage         Key = "choose_language"
	LanguageKo             Key = "language_ko"
	LanguageEn             Key = "language_en"
	SincePrompt            Key = "since_prompt"
	UntilPrompt            Key = "until_prompt"
	EnterToContinue        Key = "enter_to_continue"
	EnterToRunTabToSwitch  Key = "enter_to_run_tab_to_switch"
	InvalidSinceDate       Key = "invalid_since_date"
	InvalidUntilDate       Key = "invalid_until_date"
	LanguageLabel          Key = "language_label"
	SinceLabel             Key = "since_label"
	UntilLabel             Key = "until_label"
	ReportSummary          Key = "report_summary"
	ReportSessions         Key = "report_sessions"
	ReportDaily            Key = "report_daily"
	ReportCharging         Key = "report_charging"
	ReportDischargeProfile Key = "report_discharge_profile"
	ReportProcessImpacts   Key = "report_process_impacts"
	ReportSystemEvents     Key = "report_system_events"
	ReportSpecs            Key = "report_specs"
	ReportThermals         Key = "report_thermals"
	NoSessions             Key = "no_sessions"
	NoDailyRecords         Key = "no_daily_records"
	NoChargingSessions     Key = "no_charging_sessions"
	NoDischargeProfile     Key = "no_discharge_profile"
	NoProcessImpactData    Key = "no_process_impact_data"
	NoSystemEvents         Key = "no_system_events"
	NoHardwareSpecs        Key = "no_hardware_specs"
	NoThermalSamples       Key = "no_thermal_samples"
	AvgDischarge           Key = "avg_discharge"
	WorstSession           Key = "worst_session"
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
	TypeHeader             Key = "type_header"
	DescriptionHeader      Key = "description_header"
	SamplesHeader          Key = "samples_header"
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
	LanguageKo:             "Korean",
	LanguageEn:             "English",
	SincePrompt:            "since: ",
	UntilPrompt:            "until: ",
	EnterToContinue:        "enter to continue",
	EnterToRunTabToSwitch:  "enter to run, tab to switch",
	InvalidSinceDate:       "invalid since date",
	InvalidUntilDate:       "invalid until date",
	LanguageLabel:          "language",
	SinceLabel:             "since",
	UntilLabel:             "until",
	ReportSummary:          "Summary",
	ReportSessions:         "Sessions",
	ReportDaily:            "Daily",
	ReportCharging:         "Charging",
	ReportDischargeProfile: "Discharge Profile",
	ReportProcessImpacts:   "Process Impacts",
	ReportSystemEvents:     "System Events",
	ReportSpecs:            "Specs",
	ReportThermals:         "Thermals",
	NoSessions:             "no sessions",
	NoDailyRecords:         "no daily records",
	NoChargingSessions:     "no charging sessions",
	NoDischargeProfile:     "no discharge profile",
	NoProcessImpactData:    "no process impact data",
	NoSystemEvents:         "no system events",
	NoHardwareSpecs:        "no hardware specs",
	NoThermalSamples:       "no thermal samples",
	AvgDischarge:           "avg discharge",
	WorstSession:           "worst session",
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
	TypeHeader:             "Type",
	DescriptionHeader:      "Description",
	SamplesHeader:          "Samples",
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
	LanguageKo:             "한국어",
	LanguageEn:             "영어",
	SincePrompt:            "시작일: ",
	UntilPrompt:            "종료일: ",
	EnterToContinue:        "Enter: 다음",
	EnterToRunTabToSwitch:  "Enter: 실행, Tab: 전환",
	InvalidSinceDate:       "시작일이 올바르지 않습니다",
	InvalidUntilDate:       "종료일이 올바르지 않습니다",
	LanguageLabel:          "언어",
	SinceLabel:             "시작",
	UntilLabel:             "종료",
	ReportSummary:          "요약",
	ReportSessions:         "세션",
	ReportDaily:            "일별",
	ReportCharging:         "충전",
	ReportDischargeProfile: "방전 프로파일",
	ReportProcessImpacts:   "프로세스 영향",
	ReportSystemEvents:     "시스템 이벤트",
	ReportSpecs:            "사양",
	ReportThermals:         "온도",
	NoSessions:             "세션 없음",
	NoDailyRecords:         "일별 기록 없음",
	NoChargingSessions:     "충전 세션 없음",
	NoDischargeProfile:     "방전 프로파일 없음",
	NoProcessImpactData:    "프로세스 영향 데이터 없음",
	NoSystemEvents:         "시스템 이벤트 없음",
	NoHardwareSpecs:        "하드웨어 사양 없음",
	NoThermalSamples:       "온도 샘플 없음",
	AvgDischarge:           "평균 방전",
	WorstSession:           "가장 빠른 방전 세션",
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
	TypeHeader:             "유형",
	DescriptionHeader:      "설명",
	SamplesHeader:          "샘플",
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
