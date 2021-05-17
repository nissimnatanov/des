TARGET_DIR = ../../bin/cpp
TARGET = $(TARGET_DIR)/libsudoku_test
TMP_DIR = ../../tmp/cpp/test
CXX = g++

INCLUDES := $(wildcard ../inc/*.h)
INCLUDES_FLAGS := \
	-iquote ../inc \

CORE_LIBS = -L $(TARGET_DIR) -l sudoku

CXX_FLAGS := -std=c++14
CXX_LINK_FLAGS := -O3

MAIN_OBJS := $(patsubst %.cpp,$(TMP_DIR)/%.o,$(wildcard *.cpp))
$(MAIN_OBJS): $(INCLUDES) | $(TMP_DIR)
$(TMP_DIR)/%.o: %.cpp
	$(CXX) $(CXX_FLAGS) $(INCLUDES_FLAGS) -c $< -o $@

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

$(TARGET): $(MAIN_OBJS) $(BOARD_OBJS) $(GENERATORS_OBJS) $(SOLVER_OBJS)
	$(CXX) $(CXX_FLAGS) $(CXX_LINK_FLAGS) $(CORE_LIBS) $^ -o $(TARGET)

.PHONY: build
build: $(TARGET)

.PHONY: clean
clean:
	rm -f $(TARGET)
	rm -f -r $(TMP_DIR)
