package com.darkevilsudoku.utils;

import java.util.List;
import java.util.Random;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicLong;

public final class RandomUtils {

	private static AtomicLong seedRandomizer;
	static {
		UUID uuid = UUID.randomUUID();
		seedRandomizer = new AtomicLong(uuid.getLeastSignificantBits() ^ uuid.getMostSignificantBits());
	}

	public static long nextSeed() {
		// since both inputs are forward only and right one always increments,
		// repetition with sum is impossible
		// do NOT use XOR - it can cause seed repetition if used on two
		// incrementing numbers
		return (System.nanoTime() + seedRandomizer.getAndIncrement());
	}
	
	public static void shuffle(Random r, int[] a, int first, int count) {
		// for i from n - 1 down to 1 do
		// j = random integer with 0 <= j <= i
		// swap a[j] and a[i]

		int last = (count == -1) ? a.length - 1 : first + count - 1;
		if (last <= first) {
			throw new IndexOutOfBoundsException();
		}

		for (int i = last; i > first; --i) {
			int j = first + r.nextInt(i + 1 - first);
			if (i != j) {
				int temp = a[i];
				a[i] = a[j];
				a[j] = temp;
			}
		}
	}

	public static void shuffle(Random r, byte[] a, int first, int count) {
		// for i from n - 1 down to 1 do
		// j = random integer with 0 <= j <= i
		// swap a[j] and a[i]

		int last = (count == -1) ? a.length - 1 : first + count - 1;
		if (last <= first) {
			throw new IndexOutOfBoundsException();
		}

		for (int i = last; i > first; --i) {
			int j = first + r.nextInt(i + 1 - first);
			if (i != j) {
				byte temp = a[i];
				a[i] = a[j];
				a[j] = temp;
			}
		}
	}

	public static <T> void shuffle(Random r, List<T> a, int first, int count) {
		// for i from n - 1 down to 1 do
		// j = random integer with 0 <= j <= i
		// swap a[j] and a[i]

		int last = (count == -1) ? a.size() - 1 : first + count - 1;
		if (last <= first) {
			throw new IndexOutOfBoundsException();
		}

		for (int i = last; i > first; --i) {
			int j = first + r.nextInt(i + 1 - first);
			if (i != j) {
				T temp = a.get(i);
				a.set(i, a.get(j));
				a.set(j, temp);
			}
		}
	}

	public static <T> void shuffle(Random r, T[] a, int first, int count) {
		// for i from n - 1 down to 1 do
		// j = random integer with 0 <= j <= i
		// swap a[j] and a[i]

		int last = (count == -1) ? a.length - 1 : first + count - 1;
		if (last <= first) {
			throw new IndexOutOfBoundsException();
		}

		for (int i = last; i > first; --i) {
			int j = first + r.nextInt(i + 1 - first);
			if (i != j) {
				T temp = a[i];
				a[i] = a[j];
				a[j] = temp;
			}
		}
	}

	private RandomUtils() {
	}
}
