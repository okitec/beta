Beta – Betacode converter
=========================

**Version 0.1: Incomplete and subject to change**

Beta is a library for converting [TypeGreek](http://www.typegreek.com)-flavoured Betacode to polytonic Greek.
It can generate UTF-8 text with precomposed characters (NFC) or combining diacritics (NFD). Read this
[Go blog entry][1] for more information about text normalisation.

This implementation is independent from TypeGreek, but implements the same rules: diacritics can appear in any
order after the base character, capital Betacode letters form capital Greek letters (no asterisk), and a sigma
followed by whitespace or punctuation becomes a terminal sigma automatically.

The Greek subset of [Standard Betacode](https://www.tlg.uci.edu/encoding/) will also be implemented. This
means interpreting asterisks and putting up with diacritics after the asterisk and before the base character.

[1]: https://blog.golang.org/normalization
