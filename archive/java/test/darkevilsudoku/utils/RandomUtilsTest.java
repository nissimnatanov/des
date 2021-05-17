package test.darkevilsudoku.utils;

import static org.junit.Assert.*;

import java.util.Base64;
import java.util.HashSet;
import java.util.Random;
import java.util.logging.Logger;

import org.junit.Test;

import com.darkevilsudoku.utils.RandomUtils;

public class RandomUtilsTest {
	private Logger logger = Logger.getGlobal();

	//@Test
	@Test(timeout = 1000)
	public void testShuffle() {
		Random r = new Random();
		HashSet<String> seenValues = new HashSet<>();

		byte[] bytes = new byte[100];
		boolean[] found = new boolean[bytes.length];

		for (int i = 0; i < 100000; i++) {
			for (int j = 0; j < bytes.length; j++) {
				bytes[j] = (byte) j;
				found[j] = false;
			}
			RandomUtils.shuffle(r, bytes, 0, bytes.length);
			for (int j = 0; j < bytes.length; j++) {
				byte b = bytes[j];
				assertFalse(found[b]);
				found[b] = true;
			}
			String randString = Base64.getEncoder().encodeToString(bytes);
			if (seenValues.contains(randString)) {
				logger.severe("Repeated sequence on " + i + ": " + randString);
				assertFalse(seenValues.contains(randString));
			}
			seenValues.add(randString);
		}
	}
}
