TARGET_DIR = ../bin/cpp
TARGET = $(TARGET_DIR)/libsudoku.a
TMP_DIR = ../tmp/cpp
CXX = g++

INCLUDES := $(wildcard ./inc/*.h)
INCLUDES_FLAGS := \
	-iquote ./inc \

CXX_FLAGS := -std=c++14

BOARD_OBJS := $(patsubst %.cpp,$(TMP_DIR)/%.o,$(wildcard board/*.cpp))
$(BOARD_OBJS): $(INCLUDES) | $(TMP_DIR)/board
$(TMP_DIR)/board/%.o: board/%.cpp
	$(CXX) $(CXX_FLAGS) $(INCLUDES_FLAGS) -c $< -o $@
$(TMP_DIR)/board: | $(TMP_DIR)
	mkdir -p $(TMP_DIR)/board

GENERATORS_OBJS := $(patsubst %.cpp,$(TMP_DIR)/%.o,$(wildcard generators/*.cpp))
$(GENERATORS_OBJS): $(INCLUDES) | $(TMP_DIR)/generators
$(TMP_DIR)/generators/%.o: generators/%.cpp
	$(CXX) $(CXX_FLAGS) $(INCLUDES_FLAGS) -c $< -o $@
$(TMP_DIR)/generators: | $(TMP_DIR)
	mkdir -p $(TMP_DIR)/generators

SOLVER_OBJS := $(patsubst %.cpp,$(TMP_DIR)/%.o,$(wildcard solver/*.cpp))
$(SOLVER_OBJS): $(INCLUDES) | $(TMP_DIR)/solver
$(TMP_DIR)/solver/%.o: solver/%.cpp
	$(CXX) $(CXX_FLAGS) $(INCLUDES_FLAGS) -c $< -o $@
$(TMP_DIR)/solver: | $(TMP_DIR)
	mkdir -p $(TMP_DIR)/solver

$(TMP_DIR):
	mkdir -p $(TMP_DIR)

$(TARGET): $(BOARD_OBJS) $(GENERATORS_OBJS) $(SOLVER_OBJS)
	ar rcvs $(TARGET) $^

.PHONY: build
build: $(TARGET)

.PHONY: clean
clean:
	rm -f $(TARGET)
	rm -f -r $(TMP_DIR)
