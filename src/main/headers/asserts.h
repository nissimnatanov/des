#pragma once

#include <stdexcept>
#include <iostream>
#include <string>

class assertion_error : public std::exception
{
  private:
    const char *_testCase;
    const char *_file;
    const int _line;

  public:
    assertion_error(const char *testCase, const char *file, int line) : _testCase(testCase), _file(file), _line(line) {}

    const char *testCase() const throw() { return _testCase; }
    const char *file() const throw() { return _file; }
    int line() const throw() { return _line; }
};

template <typename T>
class Subject
{
  private:
    T _value;
    const char *_testCase;
    const char *_file;
    const int _line;

    void fail()
    {
        throw assertion_error(_testCase, _file, _line);
    }

  public:
    Subject(const T &value, const char *testCase, const char *file, int line) : _value(value), _testCase(testCase), _file(file), _line(line) {}
    Subject(T &&value, const char *testCase, int line) : _value(std::move(value)), _testCase(testCase), _line(line) {}

    void isTrue()
    {
        if (!_value)
        {
            fail();
        }
    }

    void isFalse()
    {
        if (_value)
        {
            fail();
        }
    }

    void isEqualTo(const T &other)
    {
        if (_value != other)
        {
            std::cerr << "Expected: " << other << std::endl;
            std::cerr << "Actual: " << _value << std::endl;
            fail();
        }
    }

    void isNotEqualTo(const T &other)
    {
        if (_value == other)
        {
            std::cerr << "Expected not equal to: " << other << std::endl;
            std::cerr << "Actual: " << _value << std::endl;
            fail();
        }
    }
};

template <typename T>
auto __assertThat(T &&value, const char *testCase, const char *file, int line)
{
    return Subject<std::remove_reference_t<T>>{value, testCase, file, line};
}

#define assertThat(value) __assertThat(value, __FUNCTION__, __FILE__, __LINE__)
