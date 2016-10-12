# bulletlist
A simple tool to generate an ordered list in text

## Syntax
The basic syntax is simple. Each separate argument creates a part of the master list. For example, `bulletlist "2:10"` creates the following list:
```
$ bulletlist 2:10
1. 1.  
   2.  
   3.  
   4.  
   5.  
   6.  
   7.  
   8.  
   9.  
   10.
```
As you can see, the first number (2) is the number for the list. The second number is the length of the list.

### Other parameters
You can give other parameters after the number and length. For example, if you wanted left-padded lowercase roman numerals as the list keys, you'd give `leftpad` and `type=romansmall`. The order of the extra parameters doesn't matter.
```
$ bulletlist 1:10:type=romansmall:leftpad
i.    i.
     ii.
    iii.
     iv.
      v.
     vi.
    vii.
   viii.
     ix.
      x.
```

### Child lists
You can nest lists too. To nest a list, make a list definition like usual and put it in brackets after a semicolon.
```
$ bulletlist '1:10:type=romansmall:leftpad;[2:5]'
i.    i.
     ii. 1.
         2.
         3.
         4.
         5.
    iii.
     iv.
      v.
     vi.
    vii.
   viii.
     ix.
      x.
```
To use these in bash or other similar shells, you'll probably want to enclose the whole string in single quotes (`'`), since it'll otherwise interpret the semicolon as a command separator.
