CC      ?= cc
CFLAGS  ?= -std=c11 -Wall -Wextra -Werror -Iinclude
AR      ?= ar
PREFIX  ?= /usr/local

SRC  := $(wildcard src/*.c)
OBJ  := $(SRC:.c=.o)
TSRC := $(wildcard tests/*.c)

.PHONY: build test example clean install

build: librulecat.a

librulecat.a: $(OBJ)
	$(AR) rcs $@ $^

test: librulecat.a
	@mkdir -p build
	$(CC) $(CFLAGS) $(TSRC) librulecat.a -lm -o build/test_rulecat
	./build/test_rulecat

# quickstart 예제 빌드·실행 (README quick start 원본이 실제로 컴파일되는지 확인)
example: librulecat.a
	@mkdir -p build
	$(CC) $(CFLAGS) examples/quickstart.c librulecat.a -lm -o build/quickstart
	./build/quickstart

install: librulecat.a
	install -d $(DESTDIR)$(PREFIX)/lib $(DESTDIR)$(PREFIX)/include/rulecat
	install -m 644 librulecat.a $(DESTDIR)$(PREFIX)/lib/
	install -m 644 include/rulecat/rulecat.h $(DESTDIR)$(PREFIX)/include/rulecat/
	install -d $(DESTDIR)$(PREFIX)/lib/pkgconfig
	sed "s|@PREFIX@|$(PREFIX)|" rulecat.pc.in > $(DESTDIR)$(PREFIX)/lib/pkgconfig/rulecat.pc

clean:
	rm -rf $(OBJ) librulecat.a build
