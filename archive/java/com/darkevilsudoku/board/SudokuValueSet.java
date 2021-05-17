package com.darkevilsudoku.board;

import java.util.EnumSet;
import java.util.HashMap;

public final class SudokuValueSet {
	private static HashMap<Integer, SudokuValueSet> cache;
	
	static {
		cache = new HashMap<Integer, SudokuValueSet>();
		EnumSet<SudokuValue> values = EnumSet.noneOf(SudokuValue.class);
		cache.put(0, construct(values, 0));
		
		populateCacheRec(values, SudokuValue.ONE, 0, 0);
	}

	private static void populateCacheRec(EnumSet<SudokuValue> values, SudokuValue start, int baseMask, int recDepth) {
		SudokuValue v = start;
		do
		{
			int currentMask = SudokuValue.combineMask(baseMask, v.getMask());
			if (!values.add(v))
				throw new RuntimeException();

			cache.put(currentMask, construct(values, currentMask));
			
			if (!v.isLast() && recDepth <= 4) {
				populateCacheRec(values, v.next(), currentMask, recDepth + 1);
			}

			if (!values.remove(v))
				throw new RuntimeException();

			if (v.isLast())
				break;
			
			v = v.next();
		} while (true); // breaks in the middle
	}


	private static SudokuValueSet construct(EnumSet<SudokuValue> values, int mask) {
		SudokuValueSet valueSet = new SudokuValueSet(values, mask);

		EnumSet<SudokuValue> complementValues = EnumSet.complementOf(values);
		int complementMask = SudokuValue.complementMask(mask);

		SudokuValueSet complementSet = new SudokuValueSet(complementValues, complementMask);
		valueSet.setComplement(complementSet);
		complementSet.setComplement(valueSet);
		return valueSet;
	}

	public static SudokuValueSet buildFromMask(int mask) {
		if (cache.containsKey(mask)) {
			return cache.get(mask);
		}

		mask = SudokuValue.complementMask(mask);
		if (cache.containsKey(mask)) {
			return cache.get(mask).getComplement();
		}

		throw new IndexOutOfBoundsException();
	}

	public static SudokuValueSet[] buildArray(int size) {
		return buildArray(size, null);
	}

	public static SudokuValueSet[] buildArray(int size, SudokuValueSet vs) {
		if (vs == null) {
			vs = SudokuValueSet.buildEmpty();
		}

		SudokuValueSet[] vsArray = new SudokuValueSet[size];
		for (int i = 0; i < size; i++) {
			vsArray[i] = vs;
		}

		return vsArray;
	}

	private SudokuValue[] values;
	private int mask;
	private int combined;
	private SudokuValueSet complement;

	private SudokuValueSet(EnumSet<SudokuValue> values, int mask) {
		this.values = new SudokuValue[values.size()];
		values.toArray(this.values);
		this.mask = mask;

		int combined = 0;
		for (SudokuValue value : values) {
			combined = combined * 10 + value.getValue();
		}

		this.combined = combined;
	}

	private void setComplement(SudokuValueSet complement) {
		this.complement = complement;
	}

	static public SudokuValueSet buildEmpty() {
		return buildFromMask(0);
	}

	static public SudokuValueSet buildFull() {
		return buildFromMask(SudokuValue.FULL_MASK);
	}

	static public SudokuValueSet build(SudokuValue v1) {
		return buildFromMask(v1.getMask());
	}

	static public SudokuValueSet build(SudokuValue v1, SudokuValue v2) {
		return buildFromMask(v1.getMask() | v2.getMask());
	}

	static public SudokuValueSet add(SudokuValueSet vs, SudokuValue value) {
		if (vs.contains(value))
			return vs;

		int mask = SudokuValue.combineMask(vs.getMask(), value.getMask());
		return buildFromMask(mask);
	}

	static public SudokuValueSet add(SudokuValueSet vs1, SudokuValueSet vs2) {
		int mask = SudokuValue.combineMask(vs1.getMask(), vs2.getMask());
		return buildFromMask(mask);
	}

	static public SudokuValueSet add(SudokuValueSet vs1, int vs2Mask) {
		int mask = SudokuValue.combineMask(vs1.getMask(), vs2Mask);
		return buildFromMask(mask);
	}

	static public SudokuValueSet remove(SudokuValueSet v1, SudokuValue value) {
		if (!v1.contains(value))
			return v1;

		int mask = SudokuValue.removeMask(v1.getMask(), value.getMask());
		return buildFromMask(mask);
	}

	static public SudokuValueSet remove(SudokuValueSet vs, SudokuValueSet vsToRemove) {
		int mask = SudokuValue.removeMask(vs.getMask(), vsToRemove.getMask());
		return buildFromMask(mask);
	}

	static public SudokuValueSet remove(SudokuValueSet vs, int maskToRemove) {
		int mask = SudokuValue.removeMask(vs.getMask(), maskToRemove);
		return buildFromMask(mask);
	}

	public int size() {
		return values.length;
	}
	
	public SudokuValue getAt(int i) {
		return values[i];
	}

	public boolean contains(SudokuValue value) {
		int valueMask = value.getMask();
		return (this.mask & valueMask) != 0;
	}
	
	public int getMask() {
		return mask;
	}

	public int getCombined() {
		return combined;
	}

	public SudokuValueSet getComplement() {
		return complement;
	}

	@Override
	public boolean equals(Object other) {
		if (this == other)
			return true;

		if (other == null)
			return false;

		if (other.getClass() != SudokuValueSet.class)
			return false;

		return this.mask == ((SudokuValueSet) other).mask;
	}

	public boolean equals(SudokuValueSet other) {
		if (this == other)
			return true;

		if (other == null)
			return false;

		return this.mask == other.mask;
	}

	@Override
	public int hashCode() {
		return values.hashCode();
	}

	@Override
	public String toString() {
		return Integer.toString(combined);
	}
}