.PHONY: check docs-clean fuzz

# 机械校验门（Mechanical Constitution, 协议 §16）
# 宣告任务完成前必须通过。Agent 无权宣布失败"没关系"。
check:
	@bash tools/check_docs.sh

docs-clean:
	@echo "no generated docs to clean"

# DEBT-0001 缓解：解析器模糊测试（有界会话；发现 panic 即为 bug）
fuzz:
	@go test ./internal/classfile -run '^$$' -fuzz FuzzParse -fuzztime 30s
