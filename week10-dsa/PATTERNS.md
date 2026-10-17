# Week 10 — Pattern library

Before coding, match the problem statement against the **Signal** of each
pattern. A signal is a feature of the statement (input shape, constraints,
what must be returned), not a restatement of the problem.

Before submitting, every time:

- 2-3 test cases of my own, at least one designed to break the solution.
- Big-O with every variable defined first, and the space convention stated
  (auxiliary, excluding the output, or total).

---

## Hash set — "have I seen this before?"

- **Signal:** duplicates, membership or counting over **unsorted** input,
  where comparing every pair would be O(n²).
- **Core idea:** one pass; look up, then insert.
- **Go tools:**
  - `map[T]struct{}` or `map[T]bool`. Since Go 1.24 (Swiss tables) both use
    the same slot size, so pick by meaning, not by memory.
  - `make(map[K]V, n)` to pre-size when n is known.
  - Comma-ok `v, ok := m[k]` whenever the zero value is a valid stored value.
- **Traps:**
  - Storing an index in `map[int]int` and testing `m[k] != 0`: index 0 is
    the zero value.
  - Hash operations are O(1) *expected*, not guaranteed.
- **Problems:** 217 Contains Duplicate, 1 Two Sum (complement lookup).

## Canonical key — "group things that are equivalent"

- **Signal:** group or compare items that are "the same up to" a
  transformation (letter order, case, rotation).
- **Core idea:** map every item to a canonical form. The key must work in
  both directions: equivalent items ⇒ same key, and different items ⇒
  different keys.
- **Go tools:**
  - Sorted key: `[]byte(s)` + `slices.Sort` + `string(...)`. This is only
    valid for ASCII input; for Unicode, sort runes.
  - Count key: `[26]int` works as a map key because arrays are comparable
    (slices are not).
  - A map from key to an index in the output keeps the output order
    deterministic.
- **Traps:**
  - `g := groups[i]; g = append(g, w)` does **not** update `groups[i]`. A
    slice is a header `{ptr, len, cap}` copied by value, so always store
    the result of `append` back: `groups[i] = append(groups[i], w)`.
  - Space is O(n·k), not O(n): every key is a new string of up to k bytes.
- **Problems:** 49 Group Anagrams, 242 Valid Anagram.

## Two pointers, opposite ends

- **Signal:** sorted (monotonic) input, symmetry (palindromes), a pair with
  a target, or O(1) extra space / in-place required.
- **Core idea:** `left` and `right` move inward. Every step discards
  candidates that can no longer be the answer (state the invariant).
- **Go tools:**
  - `switch {}` with no expression for a 3-way branch.
  - `a, b = b, a` to swap.
  - Bounds check first in `&&`, so short-circuit protects the index.
  - Byte arithmetic such as `c - 'a'`.
- **Traps:**
  - Not reading the statement (digits count in 125).
  - Moving both pointers in one iteration with independent `if`s.
  - 1-indexed output: convert at the boundary.
  - Unsorted input with original indices required → use a hash map instead.
- **Problems:** 125 Valid Palindrome, 167 Two Sum II.

## Length-prefix framing — "pack a list into one string and back"

- **Signal:** serialize and deserialize content where **any character may
  appear**, so a bare delimiter is ambiguous.
- **Core idea:** write `<len>;<content>`. The decoder searches for `;` only
  inside the digit field, then skips exactly `len` bytes, so it never scans
  the content.
- **Alternative:** escape both the delimiter and the escape character
  itself. This is better for streams (length unknown up front) and for
  human-readable formats.
- **Go tools:**
  - `strings.Builder` (plus `Grow` when the final size is known) and
    `WriteByte`.
  - `strconv.Itoa` / `strconv.Atoi` for multi-digit lengths.
  - `strings.IndexByte`: the position is relative to the slice you pass.
  - Substrings share bytes with the original, so slicing copies nothing.
- **Traps:**
  - A fixed 1-digit length breaks at len ≥ 10.
  - A separator *between* words makes `[]` and `[""]` identical; use a
    terminator after every word.
  - Escaping `;` without escaping `\`.
  - One wrong length silently corrupts everything after it.
  - `strconv.Atoi` accepts `"-1"`.
- **Complexity** (k = words, S = total characters, n = len(encoded)):
  - Encode: O(S + k) time. O(1) extra if the output is not counted, but
    without `Grow` the builder doubles and leaves garbage buffers (peak
    close to 2x the output). To pre-size, count the digits with a loop,
    not `len(strconv.Itoa(...))`, which allocates for lengths ≥ 100.
  - Decode: lead with **O(n) time and space** (the input size). Refine
    if asked: O(k + prefix digits) time, because it reads only the
    prefixes and skips the contents, and O(k) space for string headers,
    because substrings share bytes. Worst case: many empty words
    (`"0;"` is 2 bytes, so k = n/2).
- **Tests:** aim at the weak points of the chosen design (delimiter inside
  a word, multi-digit length) and always include `[]` and `[""]`.
- **Problems:** 271 Encode and Decode Strings.

---

## Re-solve log

| Problem | First attempt | Re-solve | Gap | Notes |
|---|---|---|---|---|
| 271 Encode/Decode | > 30 min, not solved | 12:31, 1 submit | ~18 h | Time inflated by recall. Decode index math was rebuilt, not remembered. Re-solved again on 16/10 (row below). |
| 49 Group Anagrams | 15:00, 1 submit | 11:31, 1 submit | 2 days | Inserting into `[][]string` is no longer a blocker. |
| 271 Encode/Decode (2nd) | 12:31, 1 submit (09/10) | 13:27, 1 submit | 7 days | Rebuilt alone after a week: real retention, not recall. Weak spot is now the Big-O *reasons* (said decode space comes from copying substrings; it does not). Own tests missed `[]` and `[""]`. |
