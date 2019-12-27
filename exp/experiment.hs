import Data.List

adder0 :: Int -> Int
adder0 n = sum [0..n]

--

adder1 :: Int -> Int
adder1 n = foldl (+) 0 [0..n]

--

adder2 :: Int -> Int
adder2 1 = 1
adder2 n = n + adder2 (n - 1)

--
--

counter0 :: String -> Char -> Int
counter0 [] c = 0
counter0 (x:xs) c
    | x == c = 1 + counter0 xs c
    | otherwise = counter0 xs c

--

counter1 :: String -> Char -> Int
counter1 s c = length (filter (==c) s)
