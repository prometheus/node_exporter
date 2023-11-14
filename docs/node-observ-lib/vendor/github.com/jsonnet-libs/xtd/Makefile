.PHONY: test
test:
	@cd test/; \
	jb install; \
	RESULT=0; \
	for f in $$(find . -path './.git' -prune -o -name 'vendor' -prune -o -name '*_test.jsonnet' -print); do \
		echo "$$f"; \
		jsonnet -J vendor -J lib "$$f"; \
		RESULT=$$(($$RESULT + $$?)); \
	done; \
	exit $$RESULT


.PHONY: docs
docs:
	docsonnet main.libsonnet
