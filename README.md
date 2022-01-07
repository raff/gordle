# gordle
gordle is a helper for wordle

## Usage

   gordle [-words=word file] [-skip=letters] [-contain=letters] pattern

### Where:

    -words string
        external file with words (i.e. "/usr/share/dict/words" on *nix like environments)
    -skip string
        skip words containing any of these letters
    -contain string
        accept only words that contain all these letters

This commands prints a list of 5 letter words matching the specified parameters.

- To skip letters that don't match (gray), add them to the -skip list.
- To select letters that match in any position (yellow), add them to the -contain list.
- To match letters in a specific position (green) add them to the pattern using `-`, `*` or `#` to mark letters you don't know.
    For example `##AB#` means the letters `AB` should be the 3rd and 4th letters in the word.
