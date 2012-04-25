The Burrows-Wheeler-Scott transform (called the Burrows-Wheeler transform
"Scottified" in existing literature, but that sounds silly) sorts together all
infinitely repeated cycles of each Lyndon word of the input, then takes the
last character of each rotation of each Lyndon word in the overall sorted
order. Of course, this description is about as intelligible as any of the
existing literature on it for someone not intimately familiar with the
concepts involved, and incomplete for someone who is.

Lyndon words are sequences which are less than any of the rotations of that
sequence. "Less than", that is, the order, is defined in "the usual
lexicographical way": 'a' < 'b' by definition; 'aa' < 'ab' because
the first positions are the same and the second position is lesser in 'aa';
'ab' < 'ba' because 'ab' is lesser in the first position. For now, the order
of sequences whose lengths are not equal is left undefined.

A Lyndon word is a sequence such that no matter how many times you take the
rightmost (or leftmost) element and attach it to the left (or right,
respectively), the result is never less than the word. By the Chen-Fox-Lyndon
theorem, every ordered sequence has a unique "Lyndon factorization" of Lyndon
words, such that each word in the factorization is never greater than its
predecessor, where order is _here_ (not for the BWST algorithm) defined for
words of unequal length such that if a is a prefix of b, then a < b.

The concept is best illustrated by example. The Lyndon factorization of the
sequence 'FOOBAR2000' (assume ASCII - digits precede letters) is the words
'FOO', 'B', 'AR', '2', '0', '0', and '0'. 'F' is a Lyndon word because it has
no rotations, but it's not part of the factorization because 'FO' is also a
Lyndon word: 'FO' < 'OF'. 'FOO' is a Lyndon word because
'FOO' < 'OFO' < 'OOF'. 'FOOB' is _not_ a Lyndon word; 'FOOB' is not less than
its rotation 'BFOO'.

Duval (1983) gives an algorithm for finding the Lyndon factorization of a
sequence in linear time. Wikipedia's article on Lyndon words now thankfully
has a description of the algorithm, but it does a poor job of explaining what
it actually does, so I'll describe it myself. It is key to realize that all
Lyndon words of length greater than 1 end with a character greater than the
one with which it starts. Knowing this, it's easy to realize that if, while
scanning the string for the Lyndon words, a character is encountered that is
less than the character at the start of the current word, then the word has
ended. Note, however, that this is not the only case; an illustrative example
here is 'ABCA', which factorizes to 'ABC' and 'A'. When the algorithm
encounters a character equal to the first, it has to start comparing to the
second character. Or, more generally, while only one of the two indices into
the string the algorithm holds is incremented on each step, in the case of the
compared characters being equal, both indices are incremented. Since this
causes the algorithm to treat repeated strings as equal, word boundaries are
determined according to the difference of the indices, and the lower one is
reset to the start each time the comparison yields lower earlier.

For my implementation of Lyndon factorization, see lyndon.go:7. Note that it
does not in all cases produce the true factorization to increase efficiency of
the BWST; see the comment on line 24 of that file.

Now that I think the Lyndon factorization is satisfactorily explained, the
BWST itself can be introduced. Recall that the BWST sorts the infinitely
repeated rotations of all Lyndon words of the input. Let's take a word that
David Scott, the person who developed BWST, actually used with respect to the
algorithm, which illustrates not only this concept but also one of the
problems involved in learning about the transform: 'SCOTTIFACATION'. Its
Lyndon factorization produces 'S', 'COTTIF', and 'ACATION'. All rotations of
these words are:

    S
	COTTIF
	FCOTTI
	IFCOTT
	TIFCOT
	TTIFCO
	OTTIFC
	ACATION
	NACATIO
	ONACATI
	IONACAT
	TIONACA
	ATIONAC
	CATIONA

These rotations are not sorted according to the usual lexicographical order.
In particular, strings of different lengths are compared as if both are
repeated infinitely. A shorter length to compare is each word repeated as many
times as the other has characters. Even shorter is to compare them with
indices modulus the lengths until either one is determined lesser or both are
repeated at least once. See function CyclicLess at sort.go:10.

So, if we sort the rotations, we get:

    S        ACATION
	COTTIF   ATIONAC
	FCOTTI   CATIONA
	IFCOTT    COTTIF
	TIFCOT    FCOTTI
	TTIFCO    IFCOTT
	OTTIFC   IONACAT
	ACATION  NACATIO
	NACATIO  ONACATI
	ONACATI   OTTIFC
	IONACAT        S
	TIONACA   TIFCOT
	ATIONAC  TIONACA
	CATIONA   TTIFCO

The BWST is now the last character of each rotation in the sorted output:
'NCAFITTOICSTAO'.

This was perhaps a bad example; entropy was not reduced, and the special
cyclic order never came into play. The basic concepts have been explained,
however, and that is my aim.

Now, the entire point of the BWST is that it has an inverse - you can get
'SCOTTIFACATION' back out of 'NCAFITTOICSTAO'. To do this, we need to compare
the BWST output with its sorted order. Sorting gives 'AACCFIINOOSTTT'. Now we
build a table thus:

    Index  Sorted   BWST   Start + Count = Sum   Map
	0      A        N      7       0       7     2
	1      A        C      2       0       2     12
	2      C        A      0       0       0     1
	3      C        F      4       0       4     9
	4      F        I      5       0       5     3
	5      I        T      11      0       11    4
	6      I        T      11      1       12    8
	7      N        O      8       0       8     0
	8      O        I      5       1       6     7
	9      O        C      2       1       3     13
	10     S        S      10      0       10    10
	11     T        T      11      2       13    5
	12     T        A      0       1       1     6
	13     T        O      8       1       9     11

 - Index is the zero-based index into the sorted sequence for each character.
 - Sorted is the sorted sequence.
 - BWST is the input into the inverse function.
 - Start is the first index in the sorted string at which the corresponding
   BWST character is found.
 - Count is the number of times the corresponding BWST character already has
   been found in the sequence.
 - Sum is the sum of the starts and counts.
 - Map is the line number whose sum equals the current index.

We start at index 0 and follow the map, outputting from the sorted sequence as
we go:

    Index   Sorted   Map   Output
	0       A        2     A
	2       C        1     AC
	1       A        12    ACA
	12      T        6     ACAT
	6       I        8     ACATI
	8       O        7     ACATIO
	7       N        0     ACATION

But the map at line 7 points to an index we've already visited. We've now
retrieved the lexicographically least Lyndon word from the input. Next, we
move to the lowest index we have not yet visited, which is 3.

    Index   Sorted   Map   Output
	3       C        9     C
	9       O        13    CO
	13      T        11    COT
	11      T        5     COTT
	5       I        4     COTTI
	4       F        3     COTTIF

We know the next greatest Lyndon word. The only unvisited index so far is 10,
so S at 10 is the greatest Lyndon word in the original input. Concatenating
the words in nonincreasing order yields SCOTTIFACATION. BWST inverted.

Resources:
 - http://groups.google.com/group/comp.compression/msg/a0236d754e869212 - This
   is an old post, so it misses some connections which now are known, but it
   is the only plain-English description of the BWST and UNBWST I could find.
   I'm not sure why, but my BWST inverse algorithm is different from and
   somewhat inferior to the one given here.
 - http://bijective.dogma.net/00yyy.pdf - This paper is a fairly intelligible
   description and implementation of the algorithms, but it contains several
   significant errors. Its examples are more optimized than the ones I gave;
   if you want to learn more, check it out, but not until you understand the
   algorithms enough to be able to recognize the errors.
 - http://arxiv.org/abs/0908.0239 - This paper is for those with advanced
   degrees. It is correct, but almost impossible to understand: "Let k ∈ ℕ.
   Let ⋃_{i=1}^s[v_i] = {w_1, ..., w_n} ⊆ ∑^+ be a multiset built from
   conjugacy classes [v_i]. Let M = (w_1, ..., w_n) satisfy context_k(w_1) ≤
   ··· ≤ context_k(w_n) and let L = last(w_1) ··· last(w_n) be the sequence of
   the last symbols. Then context_k(w_i) = λ_Lπ_L(i)·λ_Lπ_L^2(i)···λ_Lπ_L^k(i)
   where π_L^t denotes the t-fold application of π_L and λ_Lπ_L(i) =
   λ_L(π_L(i))."
