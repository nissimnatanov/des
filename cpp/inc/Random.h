#pragma once

#include <algorithm>
#include <array>
#include <random>
#include <vector>
#include <memory>

#include "SudokuValue.h"

class Random
{
  private:
    static std::mt19937 *new_rand();

    std::unique_ptr<std::mt19937> _r;

    std::mt19937 &get() {
        return *_r;
    }

  public:
    Random() : _r(new_rand()) {}

    Random(const Random &) = delete;
    Random &operator=(const Random &) = delete;

    template <typename T>
    void shuffle(std::vector<T> &values)
    {
        std::shuffle(values.begin(), values.end(), get());
    }

    template <typename T, std::size_t _Nm>
    void shuffle(std::array<T, _Nm> &values)
    {
        std::shuffle(values.begin(), values.end(), get());
    }

    template <typename T, std::size_t _Nm>
    void shuffle(std::array<T, _Nm> &values, int start, int count)
    {
        std::shuffle(values.begin() + start, values.begin() + (start + count), get());
    }

    int nextIndex(int size)
    {
        std::uniform_int_distribution<int> distribution(0, size - 1);
        return distribution(get());
    }

    int nextInClosedRange(int minInclusive, int maxInclusive)
    {
        std::uniform_int_distribution<int> distribution(minInclusive, maxInclusive);
        return distribution(get());
    }

    bool percentProbability(int percent)
    {
        return (1 + nextIndex(100)) <= percent;
    }
};
