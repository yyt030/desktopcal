package main

import "testing"

const sampleText = `[2026-10-01]
--------------------------------------------------------------------------------
1. 任务-处理事务1
2. 任务-处理事务2(已完成)
3. 任务-处理事务3（已完成）


[2026-09-30]
--------------------------------------------------------------------------------
1. 任务1-处理事务1
2. ✓ 已完成的任务
3. √ 另一个已完成
`

func TestParseTextEvents(t *testing.T) {
	events, err := parseTextEvents(sampleText)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("期望 2 个日期，得到 %d", len(events))
	}
	if got := events["2026-10-01"]; len(got) != 3 {
		t.Fatalf("期望 2026-10-01 有 3 条任务，得到 %d", len(got))
	} else {
		if got[0].Title != "任务-处理事务1" || got[0].Completed {
			t.Fatalf("第一条任务解析错误：%+v", got[0])
		}
		if got[1].Title != "任务-处理事务2" || !got[1].Completed {
			t.Fatalf("第二条任务(英文括号已完成)解析错误：%+v", got[1])
		}
		if got[2].Title != "任务-处理事务3" || !got[2].Completed {
			t.Fatalf("第三条任务(中文括号已完成)解析错误：%+v", got[2])
		}
	}
	if got := events["2026-09-30"]; len(got) != 3 {
		t.Fatalf("期望 2026-09-30 有 3 条任务，得到 %d", len(got))
	} else {
		if got[1].Title != "已完成的任务" || !got[1].Completed {
			t.Fatalf("✓ 前缀解析错误：%+v", got[1])
		}
		if got[2].Title != "另一个已完成" || !got[2].Completed {
			t.Fatalf("√ 前缀解析错误：%+v", got[2])
		}
	}
}

func TestParseTextEventsInvalid(t *testing.T) {
	if _, err := parseTextEvents("没有任何日期的普通文本"); err == nil {
		t.Fatal("期望无日期文本报错")
	}
	if _, err := parseTextEvents("[2026-13-99]\n----\n1. 任务"); err == nil {
		t.Fatal("期望非法日期报错")
	}
}

func TestParseTextEventsBOMAndCRCRLF(t *testing.T) {
	// 桌面日历老版本文件：UTF-8 BOM 头 + CRCRLF（\\r\\r\\n）行分隔
	text := "\ufeff[2026-10-12]\r\r\n------------------\r\r\n1. 任务A\r\r\n2. 任务B(已完成)\r\r\n"
	events, err := parseTextEvents(text)
	if err != nil {
		t.Fatal(err)
	}
	got := events["2026-10-12"]
	if len(got) != 2 {
		t.Fatalf("期望 2 条任务，得到 %d", len(got))
	}
	if got[0].Title != "任务A" || got[0].Completed {
		t.Fatalf("第一条解析错误：%+v", got[0])
	}
	if got[1].Title != "任务B" || !got[1].Completed {
		t.Fatalf("第二条解析错误：%+v", got[1])
	}
}
