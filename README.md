# gordle
gordle is a helper for wordle

## Usage

   gordle [-words=word file] [-skip=letters] [-contain=letters] pattern

### Where:

    -all
        show all words (answer and allowed) - false: only show valid answer words
    -list
        list mode: pass a list of word attempts and results
    -contain string
        accept only words that contain all these letters
    -skip string
        skip words containing any of these letters
    -sort
        sort by frequency (default true)
    -words string
        external file with words

This commands prints a list of 5 letter words matching the specified parameters.

- To skip letters that don't match (gray), add them to the -skip list.
- To select letters that match in any position (yellow), add them to the -contain list.
- To match letters in a specific position (green) add them (in uppercase) to the pattern using a symbol (like `-`, `*` or `#`) to mark letters you don't know.
    For example `##AB#` means the letters `AB` should be the 3rd and 4th letters in the word.

It is also possible to specify letters that should be in the word but are in the wrong place by adding them in lowercase.

For example in the following pattern, the letter `C` is in the right place (first letter of the word) but the letters `ra` are in the wrong place:

	C--ra

In list mode you can pass a list word attempts and results instead of adding missing letters to skip list.
This is an example where after the last pass the word pick should be ROYAL:

        > gordle -list ADOBE a.o..

            M 10 [MACHO MACRO MAJOR MANGO MANOR MASON MAYOR MOCHA MOLAR MORAL]
            C 10 [CACAO CANON CARGO CAROL CHAOS COACH COAST COCOA COMMA CORAL]
            T 8 [TALON TANGO TAROT TOAST TONAL TONGA TOPAZ TOTAL]
            S 7 [SALON SALVO SAVOR SAVOY SOAPY SOLAR SONAR]
            R 6 [RATIO RAYON RAZOR ROACH ROAST ROYAL]
            L 5 [LASSO LOAMY LOATH LOCAL LOYAL]
            P 4 [PATIO PIANO POLAR POLKA]
            O 4 [OCTAL OFFAL ORGAN OVARY]
            V 4 [VALOR VAPOR VOCAL VOILA]
            F 4 [FAVOR FOAMY FOCAL FORAY]
            W 2 [WAGON WOMAN]
            K 1 [KOALA]
            H 1 [HAVOC]
            Z 1 [ZONAL]

        > gordle -list ADOBE a.o.. MANGO .a..o

            C 5 [CHAOS COACH COAST COCOA CORAL]
            L 3 [LOATH LOCAL LOYAL]
            T 3 [TOAST TOPAZ TOTAL]
            O 3 [OCTAL OFFAL OVARY]
            R 3 [ROACH ROAST ROYAL]
            V 2 [VOCAL VOILA]
            S 2 [SOAPY SOLAR]
            F 2 [FOCAL FORAY]
            P 2 [POLAR POLKA]
            K 1 [KOALA]

        > gordle -list ADOBE a.o.. MANGO .a..o CHAOS ..ao

            P 2 [POLAR POLKA]
            T 2 [TOPAZ TOTAL]
            F 1 [FORAY]
            L 1 [LOYAL]
            O 1 [OFFAL]
            R 1 [ROYAL]
            V 1 [VOILA]

        > gordle -list ADOBE a.o.. MANGO .a..o CHAOS ..ao. FORAY .OrAy

            R 1 [ROYAL]
