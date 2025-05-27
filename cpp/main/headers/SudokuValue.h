#pragma once

#include <vector>
#include <iostream>

enum class SudokuValue
{
    EMPTY,
    ONE,
    TWO,
    THREE,
    FOUR,
    FIVE,
    SIX,
    SEVEN,
    EIGHT,
    NINE
};

class SudokuValueSet
{
  private:
    static const int FULL_MASK = 0x1FF;

  public:
    static const SudokuValueSet none()
    {
        return SudokuValueSet(0);
    }

    static const SudokuValueSet all()
    {
        return SudokuValueSet(FULL_MASK);
    }

    SudokuValueSet() noexcept : _mask(0) {}

    SudokuValueSet(const SudokuValue &value) noexcept : _mask(maskOf(value)) {}

    SudokuValueSet(const SudokuValue &value1, const SudokuValue &value2) noexcept
        : _mask(maskOf(value1) | maskOf(value2)) {}

    SudokuValueSet(const SudokuValue &value1, const SudokuValue &value2, const SudokuValue &value3) noexcept
        : _mask(maskOf(value1) | maskOf(value2) | maskOf(value3)) {}

    SudokuValueSet(const SudokuValueSet &other) noexcept : _mask(other._mask) {}

    SudokuValueSet operator+(const SudokuValue &value) const noexcept
    {
        return operator+(SudokuValueSet(value));
    }

    SudokuValueSet &operator+=(const SudokuValue &value) noexcept
    {
        return operator+=(SudokuValueSet(value));
    }

    SudokuValueSet operator+(const SudokuValueSet &other) const noexcept
    {
        return SudokuValueSet(_mask | other._mask);
    }

    SudokuValueSet &operator+=(const SudokuValueSet &other) noexcept
    {
        _mask |= other._mask;
        return *this;
    }

    SudokuValueSet operator&(const SudokuValueSet &other) const noexcept
    {
        return SudokuValueSet(_mask & other._mask);
    }

    SudokuValueSet &operator&=(const SudokuValueSet &other) noexcept
    {
        _mask &= other._mask;
        return *this;
    }

    SudokuValueSet operator~() const noexcept
    {
        return SudokuValueSet(FULL_MASK & ~_mask);
    }

    SudokuValueSet operator-(const SudokuValue &value) const noexcept
    {
        return operator-(SudokuValueSet(value));
    }

    SudokuValueSet &operator-=(const SudokuValue &value) noexcept
    {
        return operator-=(SudokuValueSet(value));
    }

    SudokuValueSet operator-(const SudokuValueSet &other) const noexcept
    {
        return SudokuValueSet(_mask & (~other)._mask);
    }

    SudokuValueSet &operator-=(const SudokuValueSet &other) noexcept
    {
        _mask &= (~other)._mask;
        return *this;
    }

    SudokuValueSet &operator=(const SudokuValueSet &other) noexcept
    {
        _mask = other._mask;
        return *this;
    }

    bool operator==(const SudokuValueSet &other) const noexcept
    {
        return _mask == other._mask;
    }

    bool operator!=(const SudokuValueSet &other) const noexcept
    {
        return !(operator==(other));
    }

    bool containsAll(const SudokuValueSet &other) const noexcept
    {
        return (_mask | other._mask) == _mask;
    }

    bool contains(SudokuValue value) const noexcept
    {
        return containsAll(SudokuValueSet(value));
    }

    int size() const noexcept;

    SudokuValue operator[](int index) const
    {
        return values().at(index);
    }

    auto begin() const noexcept
    {
        return values().cbegin();
    }
    auto end() const noexcept
    {
        return values().cend();
    }
    auto cbegin() const noexcept
    {
        return values().cbegin();
    }
    auto cend() const noexcept
    {
        return values().cend();
    }
    auto crbegin() const noexcept
    {
        return values().crbegin();
    }
    auto crend() const noexcept
    {
        return values().crend();
    }

    int combined() const noexcept;

  private:
    SudokuValueSet(int mask) : _mask(mask) {}

    const std::vector<SudokuValue> &values() const noexcept;

    static inline int maskOf(SudokuValue value)
    {
        return (value == SudokuValue::EMPTY) ? 0 : 1 << (static_cast<int>(value) - 1);
    }

    int _mask;
};

inline std::ostream &operator<<(std::ostream &os, const SudokuValue &v)
{
    os << static_cast<int>(v);
    return os;
}

inline std::ostream &operator<<(std::ostream &os, const SudokuValueSet &vs)
{
    os << vs.combined();
    return os;
}
