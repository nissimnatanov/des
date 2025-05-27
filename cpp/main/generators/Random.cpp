#include "Random.h"
#include <chrono>
#include <array>

using namespace std;

// we only use the address of this function
static void seed_function() {}

mt19937 *Random::new_rand()
{
    // Variables used in seeding
    static long long seed_counter = 0;
    int var;
    int *x = new int;

    seed_seq seed{
        // Time
        static_cast<long long>(chrono::high_resolution_clock::now()
                                   .time_since_epoch()
                                   .count()),
        // ASLR
        static_cast<long long>(reinterpret_cast<intptr_t>(&seed_counter)),
        static_cast<long long>(reinterpret_cast<intptr_t>(&var)),
        static_cast<long long>(reinterpret_cast<intptr_t>(x)),
        static_cast<long long>(reinterpret_cast<intptr_t>(&seed_function)),
        static_cast<long long>(reinterpret_cast<intptr_t>(&_Exit)),
        // Thread id (does not compile well on mingw - local var address should be enough for per-thread rand)
        // static_cast<long long>(hash<thread::id>()(this_thread::get_id())),
        // counter
        ++seed_counter,
        ++seed_counter};

    delete x;
    return new mt19937(seed);
}
