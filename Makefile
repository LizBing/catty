.PHONY: check docs-clean

# 机械校验门（Mechanical Constitution, 协议 §16）
# 宣告任务完成前必须通过。Agent 无权宣布失败"没关系"。
check:
	@bash tools/check_docs.sh

docs-clean:
	@echo "no generated docs to clean"
