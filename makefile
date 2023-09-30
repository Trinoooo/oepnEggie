
# 执行测试case，并生成覆盖率报告
test-mistake:
	sudo go test -v -short -cover -coverprofile=./mistake/mistake.out ./mistake && \
	go tool cover -html=./mistake/mistake.out -o ./mistake/mistake-coverage.html && \
	rm -f ./mistake/mistake.out

test-logs:
	sudo go test -v -short -cover -coverprofile=./logs/logs.out ./logs && \
	go tool cover -html=./logs/logs.out -o ./logs/logs-coverage.html && \
	rm -f ./logs/logs.out

test-impl:
	sudo go test -v -short -cover -coverprofile=./impl/impl.out ./impl && \
	go tool cover -html=./impl/impl.out -o ./impl/impl-coverage.html && \
	rm -f ./impl/impl.out

