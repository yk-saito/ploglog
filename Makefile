STORE_DIR := ./WriteLogPackage
PROTOBUF_DIR := ./StructureDataWithProtobuf

.PHONY:	test
test:	store protobuf

.PHONY:	store
store:
			$(MAKE) -C ${STORE_DIR} test

.PHONY: protobuf
protobuf:
			${MAKE} -C ${PROTOBUF_DIR} test