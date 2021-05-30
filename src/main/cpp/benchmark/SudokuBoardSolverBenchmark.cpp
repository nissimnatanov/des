#include <iostream>
#include "SudokuSolver.h"
#include "asserts.h"

using namespace std;

/*
>>> Game Hardest 28
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
Sudoku Result {
  action: Solve
  status: Succeeded
  level (complexity): BlackHole (26301)
  board:  6D894A9D61C7B4D2B61J2C89B2G6C5G3A8D16B
  elapsed (microseconds): 581739
  value count: 22
  steps: {
    Easy [Single In Square, 1] X 305: 305
    Easy [Single In Row, 1] X 17: 17
    Easy [Single In Column, 1] X 19: 19
    Medium [The Only Choice, 5] X 32: 160
    Hard [Identify Pairs, 20] X 15: 300
    Recursion1 [Trial & Error, 100] X 5: 500
    Recursion2 [Trial & Error, 1000] X 25: 25000
  }
}
*/
const string hardest28_str = "6D894A9D61C7B4D2B61J2C89B2G6C5G3A8D16B"; // complexity: 26301

/*
>>> Game Arto Inkala Hardest
╔═══════╦═══════╦═══════╗
║ 0.0.5 ║ 3 0.0.║ 0.0.0.║
║ 8 0.0.║ 0.0.0.║ 0.2 0.║
║ 0.7 0.║ 0.1 0.║ 5 0.0.║
╠═══════╬═══════╬═══════╣
║ 4 0.0.║ 0.0.5 ║ 3 0.0.║
║ 0.1 0.║ 0.7 0.║ 0.0.6 ║
║ 0.0.3 ║ 2 0.0.║ 0.8 0.║
╠═══════╬═══════╬═══════╣
║ 0.6 0.║ 5 0.0.║ 0.0.9 ║
║ 0.0.4 ║ 0.0.0.║ 0.3 0.║
║ 0.0.0.║ 0.0.9 ║ 7 0.0.║
╚═══════╩═══════╩═══════╝
Sudoku Result {
  action: Solve
  status: Succeeded
  level (complexity): BlackHole (23750)
  board:  B53E8F2B7B1A5B4D53C1B7C6B32C8B6A5D9B4D3F97B
  elapsed (microseconds): 42438
  value count: 23
  steps: {
    Easy [Single In Square, 1] X 697: 697
    Easy [Single In Row, 1] X 17: 17
    Easy [Single In Column, 1] X 16: 16
    Medium [The Only Choice, 5] X 12: 60
    Hard [Identify Pairs, 20] X 3: 60
    Recursion1 [Trial & Error, 100] X 9: 900
    Recursion2 [Trial & Error, 1000] X 22: 22000
  }
}
*/
const string arto_inkala_hardest_str = "B53E8F2B7B1A5B4D53C1B7C6B32C8B6A5D9B4D3F97B"; // complexity: 23750

// Other known hardest games.
const string hardest1_str = "8J36F7B9A2C5C7G457E1C3C1D68B85C1B9D4B";          // complexity: 18561
const string al_escargot_str = "1D7A9B3B2C8B96B5D53B9C1B8C26D4C3F1B4F7B7C3B"; // complexity: 13139

const vector<string> my_blackholes = {

    /*
    ╔═══════╦═══════╦═══════╗
    ║ 0.0.1 ║ 0.6 0.║ 5 0.0.║
    ║ 4 0.8 ║ 0.0.1 ║ 6 0.0.║
    ║ 0.5 0.║ 0.0.0.║ 0.1 2 ║
    ╠═══════╬═══════╬═══════╣
    ║ 0.0.4 ║ 0.0.0.║ 0.0.3 ║
    ║ 0.0.0.║ 9 0.0.║ 0.0.0.║
    ║ 0.3 0.║ 0.7 0.║ 0.6 0.║
    ╠═══════╬═══════╬═══════╣
    ║ 2 0.0.║ 0.3 0.║ 0.0.1 ║
    ║ 0.1 5 ║ 6 0.0.║ 0.7 0.║
    ║ 8 0.0.║ 7 0.0.║ 0.0.0.║
    ╚═══════╩═══════╩═══════╝
    Sudoku Result {
        action: Solve
        status: Succeeded
        level (complexity): BlackHole (41655)
        board:  B1A6A5B4A8B16C5E12B4E3C9F3B7B6A2C3C1A156C7A8B7E
        elapsed (microseconds): 113602
        value count: 25
        steps: {
            Easy [Single In Square, 1] X 902: 902
            Easy [Single In Row, 1] X 28: 28
            Easy [Single In Column, 1] X 15: 15
            Medium [The Only Choice, 5] X 10: 50
            Hard [Identify Pairs, 20] X 8: 160
            Recursion1 [Trial & Error, 100] X 15: 1500
            Recursion2 [Trial & Error, 1000] X 39: 39000
        }
    }
    */
    "B1A6A5B4A8B16C5E12B4E3C9F3B7B6A2C3C1A156C7A8B7E", // 41655

    /*
    ╔═══════╦═══════╦═══════╗
    ║ 0.0.0.║ 0.0.0.║ 0.0.0.║
    ║ 1 0.0.║ 3 8 0.║ 5 0.0.║
    ║ 0.9 0.║ 0.5 0.║ 3 1 0.║
    ╠═══════╬═══════╬═══════╣
    ║ 7 0.0.║ 5 4 0.║ 0.6 0.║
    ║ 0.3 6 ║ 0.0.0.║ 4 0.0.║
    ║ 0.0.0.║ 0.6 0.║ 0.0.8 ║
    ╠═══════╬═══════╬═══════╣
    ║ 0.0.0.║ 0.0.0.║ 0.0.0.║
    ║ 5 0.0.║ 0.3 0.║ 0.4 1 ║
    ║ 0.7 0.║ 0.1 6 ║ 8 3 0.║
    ╚═══════╩═══════╩═══════╝
    Sudoku Result {
        action: Solve
        status: Succeeded
        level (complexity): BlackHole (34777)
        board:  I1B38A5C9B5A31A7B54B6B36C4F6C8I5C3B41A7B1683A
        elapsed (microseconds): 78396
        value count: 26
        steps: {
            Easy [Single In Square, 1] X 988: 988
            Easy [Single In Row, 1] X 37: 37
            Easy [Single In Column, 1] X 12: 12
            Medium [The Only Choice, 5] X 4: 20
            Hard [Identify Pairs, 20] X 6: 120
            Recursion1 [Trial & Error, 100] X 26: 2600
            Recursion2 [Trial & Error, 1000] X 31: 31000
        }
    }
    */
    "I1B38A5C9B5A31A7B54B6B36C4F6C8I5C3B41A7B1683A", // 34777

    /*
    ╔═══════╦═══════╦═══════╗
    ║ 0.0.3 ║ 0.8 5 ║ 0.0.0.║
    ║ 6 0.0.║ 1 0.0.║ 0.0.0.║
    ║ 0.0.2 ║ 7 0.0.║ 0.5 0.║
    ╠═══════╬═══════╬═══════╣
    ║ 0.0.0.║ 0.1 7 ║ 0.8 0.║
    ║ 0.6 0.║ 0.3 0.║ 5 1 0.║
    ║ 0.0.0.║ 0.0.6 ║ 0.0.7 ║
    ╠═══════╬═══════╬═══════╣
    ║ 0.9 0.║ 3 0.8 ║ 6 0.0.║
    ║ 1 0.0.║ 0.7 0.║ 0.0.0.║
    ║ 0.0.5 ║ 0.0.0.║ 8 0.0.║
    ╚═══════╩═══════╩═══════╝
    Sudoku Result {
        action: Solve
        status: Succeeded
        level (complexity): BlackHole (33434)
        board:  B3A85C6B1G27C5E17A8B6B3A51F6B7A9A3A86B1C7F5C8B
        elapsed (microseconds): 214043
        value count: 25
        steps: {
            Easy [Single In Square, 1] X 676: 676
            Easy [Single In Row, 1] X 27: 27
            Easy [Single In Column, 1] X 11: 11
            Medium [The Only Choice, 5] X 4: 20
            Recursion1 [Trial & Error, 100] X 27: 2700
            Recursion2 [Trial & Error, 1000] X 30: 30000
        }
    }
    */
    "B3A85C6B1G27C5E17A8B6B3A51F6B7A9A3A86B1C7F5C8B", // 33434

    /*
    ╔═══════╦═══════╦═══════╗
    ║ 0.0.1 ║ 0.0.0.║ 0.4 0.║
    ║ 0.0.9 ║ 0.4 0.║ 5 0.1 ║
    ║ 4 0.0.║ 0.2 0.║ 8 3 0.║
    ╠═══════╬═══════╬═══════╣
    ║ 0.4 0.║ 0.3 2 ║ 0.7 0.║
    ║ 6 0.0.║ 4 0.0.║ 0.8 0.║
    ║ 7 0.0.║ 0.0.0.║ 0.0.4 ║
    ╠═══════╬═══════╬═══════╣
    ║ 0.8 0.║ 3 7 0.║ 0.2 0.║
    ║ 0.0.0.║ 0.6 0.║ 7 0.0.║
    ║ 0.0.0.║ 2 0.8 ║ 0.0.0.║
    ╚═══════╩═══════╩═══════╝
    Sudoku Result {
        action: Solve
        status: Succeeded
        level (complexity): BlackHole (27270)
        board:  B1D4C9A4A5A14C2A83B4B32A7A6B4C8A7G4A8A37B2E6A7E2A8C
        elapsed (microseconds): 138243
        value count: 27
        steps: {
            Easy [Single In Square, 1] X 572: 572
            Easy [Single In Row, 1] X 13: 13
            Easy [Single In Column, 1] X 15: 15
            Medium [The Only Choice, 5] X 14: 70
            Hard [Identify Pairs, 20] X 5: 100
            Recursion1 [Trial & Error, 100] X 25: 2500
            Recursion2 [Trial & Error, 1000] X 24: 24000
        }
    }
    */

    "B1D4C9A4A5A14C2A83B4B32A7A6B4C8A7G4A8A37B2E6A7E2A8C", // 27270
    "5A12B3A4G6A3H1E5D5A2C3A4A5B1A6C47F61E4B396B5",        // 26471
    "9D4B1C72A6C2A1A9A8D3A71C7E4A3B94D1B29A3B7B4A3D3E2A",  // 22162
    "J6C9A1C3C8A7A8D1D9B874C6A3B8C8A572C74B35C9C4C",       // 21106
    "1B5A64I7B5F1B16E4C79A5A79A38E16B894C4C6B9C6B17",      // 20910
    "6B5D7B9B63A8A78I3G7A68A5A3A6A578A9C84C59I6B259B",     // 20171
    "B3A8C7B5A1A38A8E1E34A7I6E764B3147B6B7A8B4B19D1C",     // 19037
    "B9E1D9B7347B2A9B5C4A7E8A6A9H8C7A84A2A92A5F4D5B",      // 18489
    "B23C1A9I1B86C8C35A2A5B8B13D61C82A65A8G2B6A49E52",     // 17469
    "6B3B5C3C6B4B2A9E4B5A8A25D8C8A64D3C9B4A8C8D5B8A3A21A", // 17320
    "J71E658C1A2A3C8E125A3A8D1A23C95A1A8A3A3A9D2C3A5A9A",  // 16304
    "7A568C1B9B2B86D12D38C25E5E62B7D2B7A6A96F7A4D5B",      // 14051
};

const vector<string> my_nightmares = {
    "G3F761A1A76D2C35C94B9B3G8B67A1C2B6B1B9A339B7C1",     // 13576
    "A3C9A8C1A5G3B6C8B36A54D9A1D64A8A3A9C6B1C2C8C139B5B", // 13449
    "41E3B37C6A19G7C14A9I7A6A4B9G37B57A541D3B5A2A4A",     // 12908
    "A1C2B76A21B3E6C1A3E2C7A9C3A2A6B3A79I924B7C13C69A8",  // 9823
    "J2B4A93A7C96428D2B8G6B3D8A95I24B59C9A5A8A2A3",       // 7182
    "A7B5I4B55B6B23G9D7A69A4C4A8A3C2E5A3C7C9B58A3A26",    // 6548
    "D9A57C956B1H8B4D9B9B7C2A23B5C742A9F53B24B7D5B2",     // 6384
};

const vector<string> boards_19_or_less = {
    "E6B4B9G7C4B3F5F2B1B46B7E8A1B9E5C76F5A", // values: 19, complexity: 524
    "C1D7B2A8I3N4A52C6F91B5C8B1A39E2E56A",   // values: 19, complexity: 408
    "E4D59E2A6E7D1D9D5A6B3D74E2C8D9D54A7F",  // values: 19, complexity: 304
};

const vector<string> my_other_games = {
    "A7A5D1C8A9C5D17B9A53C7E8B9A8E523C9B3B6A1G3C76B",
    "I2F86B4B67A9D674G96D98C3A78C29C467A8A1A1C4D",
    "A6A2C4C8A35A6G8B59F74C8E23A4C9B5F68B5A3G69B",
    "B2F47C1G7B96C6E32D5B19A3B2A7B7C1G4B92B1A7B3",
    "74A2A8B1A1B35C5F8H3D541A2F3B5B6B97A319H8A1B6B",
    "A1C3B57D9C9B716B2A4A1B2H6A7D97D9A37A5G4C37B2C1",
    "7B92A5J36B7C24D6B3A8A7B3D6A4D1B6F3A1B4B645B7D",
    "97F2B5B3A4D92D6512F4C9B1D5J7C1D4A9D7652A",
    "3D91E58E49E79D861C71I327A571A493E4B5D93F",
    "9A4A7F12A43C2C8A1F57B8D29D271C54B1F5E49F87A",
    "6C279I7D4A98C1D4G8A29D32B898F146C138C3C2B",
    "65A2H98A1D95A6B2F8B59A8A2A16A76A1A2C3H87C59H7A",
    "1A9B6C3B4G5D248E5A34B13C7E8G7E5A6914E2C6A",
    "8C75C9E8C5A6A9A3K47B69F9A7A31F5E1C2A68A5A4A1",
    "A69A4A1D79A36D8B7C7E21B4E8A1A6I7A5D2D75D4B32A",
    "E3A968A6G1C4C3C6A8A54C9C7A7B3H7B3C5F1D8A2A",
    "C2B94C5A4E42A87B5B842A1E17623F3C6B3F7K6B9A",
    "D4I87C56D23A6F9A3A52A1F9B3D8D2B7C3B9A4B2B6",
    "A4D31C2A8A97B7A2D4C9I7D5D683F5C48F12A37C4A",
    "B1C5B5C3C2A6B4937A7E8A6B3B8B7B862K5D7B6A3A9D7B",
    "B2A4E4A8C5A81C54E4D2F8A3B59B14A25C9C13A26J6B",
    "5D38G154D6A5B3B32B65B8I7E4D864C75C3A1B6F",
    "A6C2B1C4D875E6G3B9B5A4D7C3C6B1C89B59B7G6A2A",
    "9I1A2A3A4E9A76A8E3A674A6A2B5C5I2A8E4C5B3C61B",
    "C63B2C2B76A53D1C1C8A9C8C2K1C4C67A2D81A7D62B",
    "94A6C717O6A57E3C1C2A4C9F2B57B61D632C4B9C",
    "A48C1D52A489C7D6D7M27B435B8B25A1B341A8B5K",
    "F537A73C1A2B8H71A3H9B3A26B8C6B84F4D2B59F",
    "8B2A34B2B7J5A6B6B8A3C3A1C4C8C21A6C4D9G5B7A5B8A",
    "E2A58A2B5F41C2A4A1C7C5A7A3C6F923D7A89A7E4E68C",
    "B75F6C235F4A7E3B8B51A2C43B7E9G8F9D4B12B",
    "A71B2D9D8B34E2E237A4C4C591B9D3A19E7A2D4B4D7C",
    "B8D7A26C4A1A3F28B3C46E63F4D1A5A8M249B257A3",
    "B4A9A7C8A3E1C8A32A7I5C6E34C5A24A6B571E246J",
    "B4A8F37B621K68B57A4G6A7A2B8B6B15E45A7F2C16A",
    "I2C378A5A7A6C3H5A3B9A6E15B7D7A2E6F9A98D1A",
    "B4E6A3B5D17C3A2F2G9C2C8A1A439F6A5E4B62A7B8B",
    "I2C378A5A7A6C3H5A3B9A6E15B7D7A2E6F9A98D1A",
    "C5C325A4C9B92E8A76J4F8C762B4A2D1A1C8B4B8D7A",
    "8A71J8E43D6B8A6A9A592D8N91A7C7C265A3C87C",
    "B4A8F37B621K68B57A4G6A7A2B8B6B15E45A7F2C16A",
    "B4E6A3B5D17C3A2F2G9C2C8A1A439F6A5E4B62A7B8B",
    "C6G5A7A6A9G4A46C13A8A52D1K3B8D2B3A68D42C5A",
    "F6A8A46A7E9B8B2A2B1C56D24H3E5B91C3A2C6A1A9D8A",
    "92A3H4F31A27B5C2A1A34A5B7B1I7D9D31C8C67C38B",
    "C2B7B2C73D1C6D3B15C7E81I3B49D86A8C59C2A6D",
    "C7A4A1D3E48B2C5B45D11F985D3A2G86F2C29E5A",
    "A69D4D5D18G7B4A7B3B38B6H4B5A2C9A54I9B8B2B",
    "78K7A3A6B9C4E1B8C73A49E8G9B1B5A1A5B6B3E3B9",
    "B6E4C4A9D4C812A7A5A2I3B9B9G917B5D7A5D5F68",
    "D9B45E2C3A6A541I871A7D9A4C6I6D5B3B14B2B13B",
    "5D7C6A4A5E916E16A9D5A2D81F52D3B94A2C74J98A",
    "9D1B8A7A3I9C42B4A89C3G8B7A3B1B9C867A1G3A79B4B",
    "C89D2C13A6C36B2B16C7B8K4A8931A8E1E9C53A7B4D",
    "F7A9A37D8A5A2E11D8A3E721C9C3F4A6A2C5B7I6A5",
    "C32B9E78B67H5D29D15D4B2D8A95B1B4781C5F9E",
    "7D62B5E78E87A6A3B4B8C1A3D287B5A1B4M2B5A28D1A",
    "B4B2H3A726C7A1B7G9B5A4E4C956C18B6F2C5A3C6A8A",
    "A3A4A25F6I146A2G7A8B7F61A4C869B3A2H357F9A",
    "B8A7A3B4A56C9A2E64A69A5D2C8C7C7A2C1B62E1D82B7C4D",
    "B51B9A73C4E9A7A2H4F26B5A236D1B5C98B9C2B45C4B1B",
    "C3E8C1B96B3E41B24E6B3C5A2C1B84B7698A3E3A7G5B",
    "5F12C8D7B9F7B3A5G6A7B9D2A3C4E8A6A28B53B5A7A6B",
    "B4A89A1B582A1C2B4H1B6C463D7A2C4G7A1I5A79E62",
    "A2A3495E6B9A1A4I2D5A65A9B1C8F6A9L523A5C74C",
    "E9C6A4A5E8A1D21A93D7D2F89B4A342E6B5C3A84E6C",
    "32B915C8G9E3D4B7B961A8B2D8B4C16B8B7C2E6D6A9B",
    "A574I9D8A9A2E1E5C2A1C6A9C7B2C8C4A5F6C6B137A",
    "A3C12B1A5A2A7A37E5G34B38C4D9A21A6C5A8D6D6A1I2A",
    "C9A6B4D3B1C78A4A365B1C4C8C1A5E36C4A72A3C2C5B87A3F",
    "B5C6B82H19A3H4A37C2B7A8B3B9C6E6B1A8A2B5C5D42A",
    "F5C82A9B6C93C2J3B78A9D8A59B1A93A7B15A1C6A3A8A65D2",
    "B6F1B2A8D259C1D1B2B6C4C3B4D5A8C9B2A5A3B7B8F3A6",
    "5A4E7A37G89A57I3A5A7B1B6D234A7D8C24A1B75B3C1E",
    "A1B36B2C5B7D32A41B2B3C8A9A1E4B5G5B47B3A6D52J",
    "8C125A95C6A2I7A185E7B2A5A9B4C8G6E7C95B1A6B71C",
    "58A1B7F6D9B8A5J3B3B51A2B5A9A38E3A42D2E514E7A",
    "D1B4A4B6A5E38C9F85B6A1B2H697A2A4D6A38F4B6C38A",
    "C3A85B7E43E7C9A9G3B5F85B6C5G41A27D3B94A52B",
    "7A3D2G8A3F9D9A8B422B4B7A9C6B13A4B5C7C1B7D5B24C",
    "B2C71A4A1A72B6F2A4B4B31E9E3B18F935C772B1A34H9A",
    "E951B3A25A9C6E7A7E2A1C7A2F6A8A5A3A8G9C1B5A5D76A",
    "A6A459B3D8C6B5C9D6A73D1B4C78E6D3C7A5B83A5C1F4A",
    "B1A4B2A7B69C3B4B7B53B4B57A9E16L9A35D8A2C57A2G",
    "E9D9A8B37B3C15D3A6D7A9B4A6C6B7A1A5C4A1A9A1A29A7J8",
    "C42F9A8A4A7E1C4B21A8B5D6D9D2A67D2C845H2C51A",
    "C3D6B4A2D6A859A4C9D5A4B7F5B7A9A8B8B6A3D6B5A1A24F5",
    "E78C75D4A9B453H4A787C2A3B8A93A5C5D8A97G8E921B",
    "B9B4A5D69H2C1B1B58C4A9C1A7E2C85H4C196B37C8A",
    "B3B62C92A15B8C3D9B58A1B691E5B48C3G8F15C34A16A3C",
    "G6C65A92I5E184D1E88A24A3A1A1A863A5D5A9C332C5A8A",
    "A3A6C5B9B8A4F7B68B9B5A4D21A6G4A2E275A424A1C7B87C2B",
    "E79B31H7E62B62B3A5C57B9B5B6B27C3A1B62C8D1C4B7A",
    "B32D62A74B9B8J8D47I416B78A2J841B7A37C39A84",
    "H21B3B497C4A25A3463A21H82G9B16F7A7A9A6A1D5C9B",
    "C76D7D5D9F4A79B8B3B65B9A853A9C76A6D18B21F3B1D7",
    "A8D7B4A5A1A23B2C356F8D4A1C7I6A9I19A6A5A6A84A2C",
    "8D64E39E2F66B2B5D3A897D5B1B2A9F41C4B2E126A9",
    "B7A9E5A8D768A2B1D39A2D4B3B91B8D269P48B1C5A",
    "1E5A9A2D76B8C9E4C9A5D52C2B4B6C78A4D5A32A8F3E",
    "A6E3B1B478A54A5E1A2B9A6G2F781A52G8B51824C9H",
    "B1A4D8B2K7A6B48D5B9C6C8B9B1B4B3I6A746A8B129A",
    "6B34B92A34E1C8B3D7C91B6B7D2B9C87B6A9E2A4A7I4A8",
    "A1E4A6C9A1D4A7C2B35A8A69B53E48F37A893C139B57L",
    "B62E5C9F7A8A95A2D783B8A53C29B8D117A3C8C9A7A3K",
    "B68B3B7D6A94A3A9A7A8B6A7B4C7E53B1A85J41E45B1B57A9B",
    "B38A29B96F2B1D5A1C295B49F8C18D3A9A75A1A2B3B7D5D3A",
    "B8A9A5C1G5B64A1B6D5A17D6B5D9A28B7A28A6G2B9C6C7B",
    "C5B9D47A638A56B83A4B3D4D6C8B94A8D265E9F7B3B3A5A6B",
};

const vector<string> timeouts = {
    "A9B1A8C14D75D7C9A527A1A48A87A94H8E1B6D38B97C6C7C",
    "65A91E2C3A1A91E87A421D5A6E2E261H4F8B56A76F",
    "H146B712D1A3A67B8C7H3B7C6H3B7B7D5D14789B2",
    "A7B24D3A9A76A89A1A8I38G8C3B4F9C1B7C5A9C5D2A8A",
    "38A2G1C2C57J7H4B8E9A5A3E56A8D19B4845D1A",
    "2B3C46E4B17A42D9B6A4D4H5C7E5D8A28E5F6D",
    "B5E446B8K7D1C6B2B5D5A3A7A48A3B6A9E8C12A6G3",
    "6C91F2G1A3B5B5B6B82E9E3C1D9E3B7B4B121A38D",
    "5D318F5F27C57D4A53C5B4F6A9B654F38B1B9C9A5D",
    "C6C3B7F92H64B1A7A8B98K29A9H4B175C1C4A8B",
    "5D318M27C57D4A53C5B4F6A9B654F38B1B9C9A5D",
    "5D318F5F27C57D4A53C5B4F6A9B654F38B1B9C9A5D",
    "2A7J3A7C6B98A4M5F8B4B51E356A3H1A2A7C8",
    "D8A3D1B4D3A9D4C4B92A4A9D3F3A8A9A6E12C6K5A",
    "A28A4A915B7H52D7B3A5B49A9A8B27N2A36B5F2D6D",
    "A63A2A4C4A8A3A6B2N81B6I3B5C92431A56E51A4F3B",
    "1B3A4A5B4A7B68E9J3C9D56A8A3B6G19D5624M",
    "1B3A4A5B4A7B68E9H2A3C9D56A8A3B6H9D5624M",
    "5D7C9D6A4E4B96A3A7D2C21E276G3D2B5E1C4D67",
    "6D9G215C4E8H7A1A4B6C5G6B7A9B5D854B2C7B3B",
    "A93B2B1D5A3B8A7A1B9F48D4C937H6K9175E5C67A",
    "C6F971A8D2E732A1A6A8C7A24A3D4L245C7H9E",
    "A4E691H9F8E5C6A6B78A95A57H38D1E7D9C36A4",
    "D24A5D31D2B9F3D7A496C78D7G9A6E714B8B2F5B",
    "4A526C87G98B5B1H7B5B6C8E35F9A1D8B4A9D2B6C",
    "4I2A8E5A7A9B3H1A3B6A2C1A4C5C8C31H3657E9B",
    "J7C9A3A5J4A9A31A6C54A2D2C4A2B94A8D5A3C61C7A4A3",
    "A9A4B26B169E4D6I6D3A75B958J51A928C2A3D5F6",
};

void run(SudokuBoardConstShared board, SudokuSolverOptions options, string name, SudokuSolverStatus expectedStatus)
{
    SudokuSolverShared solver = createSolver();

    SudokuResultConstShared result = solver->run(options, board);

    if (expectedStatus == SudokuSolverStatus::SUCCEEDED &&
        result->getStatus() != SudokuSolverStatus::SUCCEEDED)
    {
        cerr << "FAILED: " << serializeBoard(board.get()) << endl;
    }
    assertThat(result->getStatus()).isEqualTo(expectedStatus);
    if (expectedStatus != SudokuSolverStatus::SUCCEEDED)
    {
        return;
    }

    if (result->getOptions().getAction() != SudokuSolverAction::SOLVE ||
        result->getLevel() >= SudokuLevel::NIGHTMARE ||
        result->getOriginalBoard()->getFreeCellCount() > 61)
    {
        cerr << ">>> Game " << name << endl;
        if (result->getOptions().getAction() == SudokuSolverAction::SOLVE &&
            result->getLevel() >= SudokuLevel::BLACKHOLE)
        {
            cerr << result->getOriginalBoard();
        }
        cerr << *result;
    }
}

void run(SudokuBoardConstShared board, SudokuSolverOptions options, string name)
{
    run(board, options, name, SudokuSolverStatus::SUCCEEDED);
}

void runTimed(SudokuBoardConstShared board, SudokuSolverOptions options, string name)
{
    bool prevIntegrityChecks = setIntegrityChecks(false);

    run(board, options, name);

    setIntegrityChecks(prevIntegrityChecks);
}

void runTimed(const vector<string> &boards, SudokuSolverOptions options, const string &name_prefix)
{
    for (int i = 0; i < boards.size(); i++)
    {
        SudokuBoardConstShared board = cloneAsImmutable(deserializeSudokuBoard(boards[i]));
        string name = name_prefix + to_string(i);
        runTimed(board, options, name);
    }
}

void testTimedOutBoard(SudokuBoardConstShared board)
{
    SudokuSolverOptions defaultSolveOptions(SudokuSolverAction::SOLVE);
    SudokuSolverOptions defaultSolveFastOptions(SudokuSolverAction::SOLVE_FAST);
    SudokuSolverOptions defaultProveOptions(SudokuSolverAction::PROVE);

    // Boards should never timeout with Prove and SolveFast!
    SudokuSolverShared solver = createSolver();
    SudokuResultConstShared result = solver->run(defaultProveOptions, board);
    if (result->getStatus() == SudokuSolverStatus::TIMEOUT)
    {
        cerr << *result.get();
        throw logic_error("Board timed out on Solve Fast (for previously proven board)");
    }
    result = solver->run(defaultSolveFastOptions, board);
    if (result->getStatus() == SudokuSolverStatus::TIMEOUT)
    {
        cerr << *result.get();
        throw logic_error("Board timed out on Prove");
    }

    if (result->getStatus() == SudokuSolverStatus::SUCCEEDED)
    {
        result = solver->run(defaultProveOptions, board);
        if (result->getStatus() == SudokuSolverStatus::TIMEOUT)
        {
            cerr << *result.get();
            throw logic_error("Board timed out on Solve (for previously proven board)");
        }
    }
}
void testTimedOutBoards(const vector<string> &boards)
{
    for (int i = 0; i < boards.size(); i++)
    {
        SudokuBoardConstShared board = cloneAsImmutable(deserializeSudokuBoard(boards[i]));
        testTimedOutBoard(board);
    }
}

void myBoards()
{
    SudokuSolverOptions defaultSolveOptions(SudokuSolverAction::SOLVE);
    SudokuSolverOptions defaultSolveFastOptions(SudokuSolverAction::SOLVE_FAST);
    SudokuSolverOptions defaultProveOptions(SudokuSolverAction::PROVE);

    /*
    SudokuSolverOptions largeTimeoutSolveOptions(SudokuSolverAction::SOLVE);
    largeTimeoutSolveOptions.setMaxSolverTime(600s);
    SudokuSolverOptions largeTimeoutProveOptions(SudokuSolverAction::PROVE);
    largeTimeoutProveOptions.setMaxSolverTime(600s);
*/

    runTimed(my_blackholes, defaultSolveOptions, "My Blackhole Board ");
    runTimed(my_blackholes, defaultProveOptions, "My Blackhole Board (Prove) ");

    runTimed(my_nightmares, defaultSolveOptions, "My Nightmare Board ");
    runTimed(my_nightmares, defaultProveOptions, "My Nightmare Board (Prove) ");

    /*
    runTimed(boards_19_or_less, defaultSolveOptions, "My 19 or Less Board ");
    runTimed(boards_19_or_less, defaultProveOptions, "My 19 or Less Board (Prove) ");

    cout << "other games" << endl;
    runTimed(my_other_games, defaultSolveOptions, "My Game ");
    cout << "other games (prove)" << endl;
    runTimed(my_other_games, defaultProveOptions, "My Game (Prove) ");*/
}

void testTimedOutBoards()
{
    testTimedOutBoards(timeouts);
}

void runSudokuBoardSamplesBenchmark()
{
    cerr << "-----------------------------" << endl;
    cerr << "Running Sudoku Board Samples Benchmark..." << endl;

    SteadyTimePoint startTime = SteadyClock::now();

    SudokuBoardConstShared arto_inkala_hardest = cloneAsImmutable(deserializeSudokuBoard(arto_inkala_hardest_str));
    SudokuBoardConstShared hardest28 = cloneAsImmutable(deserializeSudokuBoard(hardest28_str));
    SudokuBoardConstShared hardest1 = cloneAsImmutable(deserializeSudokuBoard(hardest1_str));
    SudokuBoardConstShared al_escargot = cloneAsImmutable(deserializeSudokuBoard(al_escargot_str));

    runTimed(arto_inkala_hardest, SudokuSolverAction::SOLVE, "Arto Inkala Hardest");
    runTimed(arto_inkala_hardest, SudokuSolverAction::PROVE, "Arto Inkala Hardest");
    runTimed(arto_inkala_hardest, SudokuSolverAction::SOLVE_FAST, "Arto Inkala Hardest");
    runTimed(hardest28, SudokuSolverAction::SOLVE, "Hardest 28");
    runTimed(hardest28, SudokuSolverAction::PROVE, "Hardest 28");
    runTimed(hardest28, SudokuSolverAction::SOLVE_FAST, "Hardest 28");
    runTimed(hardest1, SudokuSolverAction::SOLVE, "Hardest 1");
    runTimed(al_escargot, SudokuSolverAction::SOLVE, "AL Escargot");

    myBoards();
    testTimedOutBoards();

    SteadyTimePoint endTime = SteadyClock::now();
    int elapsed_milliseconds = chrono::duration_cast<std::chrono::milliseconds>(endTime - startTime).count();
    cerr << "Elapsed (milliseconds): " << elapsed_milliseconds << endl;
}
