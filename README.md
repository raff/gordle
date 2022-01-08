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
- To match letters in a specific position (green) add them (in uppercase) to the pattern using a symbol (like `-`, `*` or `#`) to mark letters you don't know.
    For example `##AB#` means the letters `AB` should be the 3rd and 4th letters in the word.

It is also possible to specify letters that should be in the word but are in the wrong place by adding them in lowercase.

For example in the following pattern, the letter `C` is in the right place (first letter of the word) but the letters `ra` are in the wrong place:

	C--ra



    
