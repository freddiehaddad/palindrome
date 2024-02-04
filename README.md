# Multithreaded Palindrome Detector

Program to check if a string is a palindrome using channels and concurrency.

The program works by spawning three threads with specific responsibilities. One
thread reads characters starting from the left and feeds palindrome characters
to a unique channel. The second thread reads characters from the right also
feeding palindrome characters to another unique channel. The third thread reads
characters from these two channels and compares them until the indices cross or
two non-matching characters are encountered and reports the result of comparing
to another separate channel.

Palindrome characters are alphanumeric. Anything else is ignored. For example,
in the string `"A man, a plan, a canal: Panama"`, the spaces, commas and the
colon are ignored and the palindrome comparison occurs as if the string was
`"amanaplanacanalpanama"`.

Also note that the strings are case insensitive.

```text
+-------------------------+
| "Madam, I'm Adam"     G |
|                         |
|                         +--+---------------------+
| isPalindrome()          |  |                     |
+--+----------------------+  |                     |
   |                         v                     v
   |   +--------------------------+  +--------------------------+
   |   | "Madam, I"             G |  | "I'm Adam"             G |
   |   |                          |  |                          |
   |   |                          |  |                          |
   |   | feedCharactersForward()  |  | feedCharactersReverse()  |
   |   +----+---------------------+  +-------------+------------+
   |        |                                      |
   |        |            +-------------------------+
   |        |            |
   |        v            v
   |   +--------------------------+
   |   | "madami"     "madami"  G |
   |   |                          |
   +-->|                          |
       | compareCharacters()      |
       +--------------------------+

G: Go Routine
```
