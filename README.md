# Notebook Battery Analyzer

내 리눅스(x86) 노트북의 각종 로그를 참조해서 배터리 사용 기록 분석하는 툴


## What it does

- `upower` 리눅스 history 파일에서 배터리 사용 체크
- `journalctl` sleep, resume, shutdown 같은 놋북 상태 체크
- PC 스펙, 충전, 방전 프로파일, 온도 등.. 리포트 작성
- 추후 OS (윈도우) 추가랑 언어는 `ko` / `en` 확장성 초안 추가


## Run

```bash
go run .
go run . 2026-03-01 2026-03-07
```


## Layout

- `main.go` - config load and app entry
- `internal/infrastructure` - log / history loaders
- `internal/service` - session and summary analysis
- `internal/ui/tui` - language and date input
- `internal/ui/renderer` - report rendering
- `internal/ui/i18n` - ko/en strings

