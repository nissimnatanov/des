package solver_test

import "github.com/nissimnatanov/des/go/solver"

type testBoard struct {
	name     string
	board    string
	solution string
	// leave default for StatusSucceeded
	expected           solver.Status
	expectedLevel      solver.Level
	expectedComplexity solver.StepComplexity
	failToLog          bool
}

var benchBoards = []testBoard{
	{
		/*
			╔═══════╦═══════╦═══════╗
			║ 6 0.0.║ 0.0.8 ║ 9 4 0.║
			║ 9 0.0.║ 0.0.6 ║ 1 0.0.║
			║ 0.7 0.║ 0.4 0.║ 0.0.0.║
			╠═══════╬═══════╬═══════╣
			║ 2 0.0.║ 6 1 0.║ 0.0.0.║
			║ 0.0.0.║ 0.0.0.║ 2 0.0.║
			║ 0.8 9 ║ 0.0.2 ║ 0.0.0.║
			╠═══════╬═══════╬═══════╣
			║ 0.0.0.║ 0.6 0.║ 0.0.5 ║
			║ 0.0.0.║ 0.0.0.║ 0.3 0.║
			║ 8 0.0.║ 0.0.1 ║ 6 0.0.║
			╚═══════╩═══════╩═══════╝
		*/
		name:               "Hardest 28",
		board:              "6D894A9D61C7B4D2B61J2C89B2G6C5G3A8D16B",
		solution:           "62_5_1_7_8943_94_8_3_2_615_7_3_71_9_45_8_6_2_25_7_619_3_8_4_4_6_3_5_8_7_29_1_1_894_3_25_7_6_7_9_2_8_63_4_1_55_1_6_2_9_4_7_38_83_4_7_5_162_9_",
		expectedLevel:      solver.LevelBlackHole,
		expectedComplexity: 230194,
	},
	{
		name:               "al escargot",
		board:              "1D7A9B3B2C8B96B5D53B9C1B8C26D4C3F1B4F7B7C3B",
		solution:           "16_2_8_5_74_93_5_34_1_29_6_7_87_8_964_3_52_1_4_7_531_2_98_6_9_13_5_86_7_4_262_8_7_9_41_3_5_35_6_4_7_8_2_19_2_41_9_3_5_8_6_78_9_72_6_1_35_4_",
		expectedLevel:      solver.LevelBlackHole,
		expectedComplexity: 145014,
	},
	{
		name:     "Arto Inkala",
		board:    "B53E8F2B7B1A5B4D53C1B7C6B32C8B6A5D9B4D3F97B",
		solution: "1_4_532_7_6_9_8_83_9_6_5_4_1_27_6_72_9_18_54_3_49_6_1_8_537_2_2_18_4_73_9_5_67_5_329_6_4_81_3_67_54_2_8_1_99_8_47_6_1_2_35_5_2_1_8_3_976_4_",
	},
	{
		name:     "hardest 1",
		board:    "8J36F7B9A2C5C7G457E1C3C1D68B85C1B9D4B",
		solution: "81_2_7_5_3_6_4_9_9_4_368_2_1_7_5_6_75_4_91_28_3_1_54_2_3_78_9_6_3_6_9_8_4572_1_2_8_7_16_9_5_34_5_2_19_7_4_3_684_3_852_6_9_17_7_96_3_1_8_45_2_",
	},
}

var otherBoards = []testBoard{
	{
		name:     "less than 19 (1)",
		board:    "E6B4B9G7C4B3F5F2B1B46B7E8A1B9E5C76F5A",
		solution: "8_3_1_5_7_69_2_44_6_92_3_1_8_7_5_2_75_8_9_41_6_31_2_7_3_6_9_54_8_9_8_3_4_25_7_16_5_461_8_72_3_9_7_5_86_13_4_92_3_1_4_9_52_6_8_769_2_7_4_8_3_51_",
	},
	{
		name:     "less than 19 (2)",
		board:    "C1D7B2A8I3N4A52C6F91B5C8B1A39E2E56A",
		solution: "5_3_6_12_9_4_8_77_9_24_86_1_3_5_8_4_1_5_7_39_2_6_3_2_8_7_9_1_6_5_4_9_1_46_523_7_8_65_7_3_4_8_2_914_6_52_3_7_81_9_18_396_5_7_4_2_27_9_8_1_4_563_",
	},
	{
		name:     "less than 19 (3)",
		board:    "E4D59E2A6E7D1D9D5A6B3D74E2C8D9D54A7F",
		solution: "2_7_3_5_8_41_9_6_8_597_6_1_3_4_21_64_3_2_9_5_78_7_4_6_13_2_8_5_99_1_2_4_58_63_7_38_5_6_9_742_1_5_3_1_27_6_9_84_6_2_8_94_3_7_1_549_78_1_5_2_6_3_",
	},
	{
		name:     "Hardest Nightmare from C++",
		board:    "G3F761A1A76D2C35C94B9B3G8B67A1C2B6B1B9A339B7C1",
		solution: "8_6_9_5_1_2_4_37_5_2_3_4_9_7618_14_768_3_5_9_22_1_6_354_7_8_947_8_96_1_32_5_9_3_5_7_2_81_4_675_18_3_9_26_4_68_2_14_5_97_3394_2_76_8_5_1",
	},
	{
		name:     "1",
		board:    "A12G5A1A6A4D32B6C7D58C8B61A9K2C83A4B9B7D4A3C",
		solution: "6_129_8_4_7_3_5_3_58_17_62_49_7_9_4_325_8_61_1_3_76_4_2_9_584_2_5_83_9_617_98_6_7_5_1_3_2_4_5_6_9_21_7_4_832_43_5_98_1_76_8_7_1_46_35_9_2_",
	},
	{
		name:     "2",
		board:    "A8C4B3B2C7B5C8B64E91B2B1E3C28B6B78A3B9E6B5B5C4B",
		solution: "6_89_5_7_42_1_34_3_29_6_1_75_8_57_1_3_82_9_647_5_8_6_3_914_2_29_6_14_5_3_8_7_31_4_7_285_9_61_4_785_36_2_99_2_3_4_1_68_7_58_6_52_9_7_43_1_",
	},
	{
		/*
			>>> Level: Nightmare
			>>> Complexity: 41428
		*/
		name:               "3",
		board:              "B4B7C5D9D1A6B5B6B1A8A4A38F1D2E37D25B89C63E29B",
		solution:           "8_9_43_5_76_1_2_56_3_2_1_97_8_4_7_12_68_4_53_9_62_5_19_83_47_389_7_4_6_2_5_14_7_1_5_23_8_9_6_9_378_6_1_4_252_4_897_5_1_631_5_6_4_3_297_8_",
		expectedLevel:      solver.LevelDarkEvil,
		expectedComplexity: 43141,
	},
	{
		// hardest so far created by me (with C++ solver)
		name:               "my board 1",
		board:              "B1A6A5B4A8B16C5E12B4E3C9F3B7B6A2C3C1A156C7A8B7E",
		solution:           "7_2_13_68_54_9_49_82_5_163_7_6_53_4_9_7_8_125_7_48_2_6_1_9_31_8_6_94_3_7_2_5_9_32_1_75_4_68_26_7_5_34_9_8_13_1568_9_2_74_84_9_71_2_3_5_6_",
		expectedLevel:      solver.LevelBlackHole,
		expectedComplexity: 176657,
	},
	{
		name:     "my board 2",
		board:    "I1B38A5C9B5A31A7B54B6B36C4F6C8I5C3B41A7B1683A",
		solution: "6_5_3_1_7_9_2_8_4_14_7_382_59_6_2_98_6_54_317_72_9_548_1_63_8_367_2_1_45_9_4_1_5_9_63_7_2_83_8_1_4_9_5_6_7_2_56_2_8_37_9_419_74_2_16835_",
	},
	{
		name:               "my board 3",
		board:              "B3A85C6B1G27C5E17A8B6B3A51F6B7A9A3A86B1C7F5C8B",
		solution:           "7_1_34_852_6_9_65_8_19_2_7_4_3_9_4_276_3_1_58_5_3_9_2_174_86_4_67_8_39_512_8_2_1_5_4_69_3_72_94_35_867_1_18_6_9_74_3_2_5_3_7_56_2_1_89_4_",
		expectedLevel:      solver.LevelBlackHole,
		expectedComplexity: 152913,
	},
	{
		name:               "my board 4",
		board:              "F5E7A8A9I62A7A4A6B9F71A6A3B95B2495A8F2D1C8D",
		solution:           "7_9_3_6_1_4_52_8_6_5_2_73_81_94_8_4_1_2_9_5_7_3_621_75_49_68_3_93_5_8_6_2_4_714_68_37_1_952_3_24956_81_7_5_8_6_1_27_3_4_9_17_9_4_83_2_6_5_",
		expectedLevel:      solver.LevelNightmare,
		expectedComplexity: 99387,
	},
	{
		/*
			>>> Level: Nightmare
		*/
		name:               "my board 5",
		board:              "B3A1B2D7A8A9E9A7D7B96J14A8D5C4A56C5B12A3B1G5",
		solution:           "7_9_36_14_5_28_6_5_2_73_81_94_8_4_1_2_95_73_6_2_1_75_4_968_3_9_3_5_8_6_2_4_7_146_83_7_1_9_52_3_2_49_568_1_7_58_6_127_34_9_17_9_4_8_3_2_6_5",
		expectedComplexity: 69728,
	},
	{
		/*
			>>> Level: Nightmare
		*/
		name:               "my board 6",
		board:              "A936B5C5B381B8H2A7C6C3F14F5E5B1C6A27A49C4E",
		solution:           "7_9361_4_52_8_6_52_7_3819_4_84_1_2_9_5_7_3_6_21_75_4_9_68_3_9_35_8_6_2_4_7_146_8_3_7_1_9_52_3_2_4_9_56_8_17_5_8_61_273_491_7_9_48_3_2_6_5_",
		expectedComplexity: 79614,
	},
	{
		/*
			>>> Level: Nightmare
			>>> Complexity: 70197
		*/
		name:     "my board 7",
		board:    "C6A4C6C38E1A9D2B5B6C3C2B1B8A7A9H8B5F49A79C26A",
		solution: "7_9_3_61_45_2_8_65_2_7_381_9_4_8_4_12_95_7_3_6_21_7_54_9_68_3_9_35_8_6_24_7_14_6_83_71_95_2_3_2_4_9_5_6_81_7_58_6_1_2_7_3_491_794_8_3_265_",
	},
	{
		name:               "my board 8",
		board:              "C6F5B3A1A4F7C1B4C3A35B2A7A4K95A8D612D17C3B5",
		solution:           "7_9_3_61_4_5_2_8_6_52_7_38_19_48_4_1_2_9_5_73_6_2_17_5_49_6_8_39_358_6_24_71_46_8_3_7_1_9_5_2_3_2_4_956_81_7_5_8_6127_3_4_9_179_4_8_32_6_5",
		expectedLevel:      solver.LevelNightmare,
		expectedComplexity: 106650,
	},
	{
		name:               "my board 9",
		board:              "A1A97B45E37D7E3C195D59I6B1B2B7F4A5C6B9634C7",
		solution:           "3_18_972_6_459_2_5_4_6_378_1_4_6_78_5_1_9_2_36_3_4_1958_7_2_1_592_8_7_4_3_6_8_7_2_63_4_15_9_28_3_71_6_5_9_4_7_41_52_9_3_68_5_96348_2_1_7",
		expectedLevel:      solver.LevelNightmare,
		expectedComplexity: 96227,
	},
	{
		name:               "my board 10",
		board:              "9E3C4B2B59B5A8C23B65I8B3B93B5C3B6A4C8C7E74B2B",
		solution:           "92_8_5_7_4_31_6_7_43_1_26_8_591_6_59_83_7_4_237_2_651_9_8_4_4_5_6_7_9_81_2_38_1_934_2_56_7_2_31_8_69_47_5_5_84_2_3_76_9_1_6_9_741_5_23_8_",
		expectedLevel:      solver.LevelNightmare,
		expectedComplexity: 101193,
	},
	{
		name:          "my board 11",
		board:         "A67D3A2A5G8B6A7F92A81D879B1B63F6B8A2I9A3A2B67A",
		expectedLevel: solver.LevelDarkEvil,
		solution:      "4_679_2_1_5_38_21_58_7_3_4_9_6_3_89_4_65_71_2_6_7_4_5_923_815_2_3_1_8796_4_19_8_634_2_5_7_9_4_67_5_81_23_7_5_2_3_1_6_8_4_98_31_24_9_675_",
	},
	{
		name:          "my board 12",
		board:         "A5B1B3A4F7C6B72E1G4A5C258D6B76C5A1F25F8C4",
		expectedLevel: solver.LevelNightmare,
		solution:      "8_57_2_14_9_36_42_3_5_9_6_8_71_9_1_68_3_724_5_2_7_9_16_8_4_5_3_6_3_47_59_1_8_2581_4_2_3_69_7_762_9_4_53_18_1_4_8_3_7_256_9_3_9_5_6_81_7_2_4",
	},
	{
		name:          "my board 13",
		board:         "9E261F5D2B3A7C42B1F7B3A8A9A3C24A8C9A7E8C2C59A1A",
		expectedLevel: solver.LevelNightmare,
		solution:      "93_5_8_4_7_2617_4_6_1_9_2_58_3_1_8_25_6_34_79_3_7_428_5_19_6_6_2_1_9_74_8_35_85_96_31_7_4_241_83_2_6_95_75_9_3_7_1_86_2_4_26_7_4_593_18_",
	},
	{
		name:          "my board 14",
		board:         "A63F98D72E98C6B9B81H62A71B5A3C6A2H76C7B8D1",
		expectedLevel: solver.LevelNightmare,
		solution:      "5_632_4_7_1_8_9_984_5_6_1_723_1_2_7_3_984_5_6_64_5_93_2_817_8_3_9_7_1_4_5_622_716_8_59_34_4_5_61_29_3_7_8_3_1_8_4_762_9_5_79_2_85_3_6_4_1",
	},
	{
		name:          "my board 15",
		board:         "6B81A7B7F9C2E1C9D41A7A4F4B8A2A9A1B4A5H8A28C5B9",
		expectedLevel: solver.LevelNightmare,
		solution:      "64_9_812_73_5_71_3_4_5_6_8_92_8_5_23_9_7_6_4_13_6_8_92_1_5_7_412_75_43_9_6_8_5_9_46_7_81_23_97_12_8_43_56_4_3_5_1_6_9_2_87_286_7_3_54_1_9",
	},
	{
		name:               "my board 16",
		board:              "7E1B5E746A61G7C453C3A9D4C1A8C8A5B6E6A3A7I4",
		expectedLevel:      solver.LevelNightmare,
		expectedComplexity: 114620,
		solution:           "73_4_9_5_6_18_2_52_9_1_3_8_7468_614_7_2_3_9_5_9_78_2_6_4531_6_1_38_95_4_2_7_45_2_3_17_86_9_2_87_54_9_61_3_1_4_5_62_39_78_3_9_6_7_8_1_2_5_4",
	},
	{
		name:               "my board 17 (BlackHole)",
		board:              "F7A3A6A9C5D2J258A3A2B8A6B8B6I3B1A918C3A8D92A5",
		expectedLevel:      solver.LevelBlackHole,
		expectedComplexity: 144157,
		solution:           "9_1_8_4_5_6_72_32_67_93_1_4_58_5_4_3_28_7_9_1_6_6_7_4_3_1_2589_35_27_9_81_64_1_89_5_64_3_7_2_7_2_5_6_4_38_9_14_9182_5_6_37_83_6_1_7_924_5",
	},
}
