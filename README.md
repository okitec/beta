Beta – Betacode converter
=========================

**Version 0.2: Incomplete and subject to change**

Beta is a library for converting [TypeGreek](http://www.typegreek.com)-flavoured Betacode to polytonic Greek.
It can generate UTF-8 text with precomposed characters (NFC) or combining diacritics (NFD). Read this
[Go blog entry][norm] for more information about text normalisation.

This implementation is independent from TypeGreek, but implements the same rules: diacritics can appear in any
order after the base character, capital Betacode letters form capital Greek letters (no asterisk), and a sigma
followed by whitespace or punctuation becomes a terminal sigma automatically.

Sort-of parsing the Greek subset of [Standard Betacode](https://www.tlg.uci.edu/encoding/) is now implemented.
This means interpreting asterisks and putting up with diacritics after the asterisk and before the base
character, but unlike Standard Betacode, this implementation is never case-insensitive.

[norm]: https://blog.golang.org/normalization

The directory `beta` contains a minimal example program that uses `beta.Writer`.
