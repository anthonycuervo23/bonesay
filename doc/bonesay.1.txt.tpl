bonesay(1)
=========

Name
----
Neo bonesay/bonethink - configurable speaking/thinking bone (and a bit more)

SYNOPSIS
--------
bonesay [-e _eye_string_] [-f _bonefile_] [-h] [-l] [-n] [-T _tongue_string_] [-W _column_] [-bdgpstwy]
       [--random] [--bold] [--rainbow] [--aurora] [--super] [_message_]

DESCRIPTION
-----------
_Neo-bonesay_ (bonesay) generates an ASCII picture of a bone saying something provided by the
user.  If run with no arguments, it accepts standard input, word-wraps
the message given at about 40 columns, and prints the bone saying the
given message on standard output.

To aid in the use of arbitrary messages with arbitrary whitespace,
use the *-n* option.  If it is specified, the given message will not be
word-wrapped.  This is possibly useful if you want to make the bone
think or speak in figlet(6).  If *-n* is specified, there must not be any command-line arguments left
after all the switches have been processed.

The *-W* specifies roughly (where the message should be wrapped. The default
is equivalent to *-W 40* i.e. wrap words at or before the 40th column.

If any command-line arguments are left over after all switches have
been processed, they become the bone's message. The program will not
accept standard input for a message in this case.

There are several provided modes which change the appearance of the
bone depending on its particular emotional/physical state. 

The *-b* option initiates Borg mode

*-d* causes the bone to appear dead 

*-g* invokes greedy mode

*-p* causes a state of paranoia to come over the bone

*-s* makes the bone appear thoroughly stoned

*-t* yields a tired bone

*-w* is somewhat the opposite of *-t* and initiates wired mode

*-y* brings on the bone's youthful appearance.

The user may specify the *-e* option to select the appearance of the bone's eyes, in which case
the first two characters of the argument string _eye_string_ will be used. The default eyes are 'oo'. The tongue is similarly
configurable through *-T* and _tongue_string_; it must be two characters and does not appear by default. However,
it does appear in the 'dead' and 'stoned' modes. Any configuration
done by *-e* and *-T* will be lost if one of the provided modes is used.

The *-f* option specifies a particular bone picture file (``bonefile'') to
use. If the bonefile spec contains '/' then it will be interpreted
as a path relative to the current directory. Otherwise, bonesay
will search the path specified in the *BONEPATH* environment variable. If *-f -* is specified, provides
interactive Unix filter (command-line fuzzy finder) to search the bonefile.

To list all bonefiles on the current *BONEPATH*, invoke *bonesay* with the *-l* switch.

*--random* pick randomly from available bonefiles

*--bold* outputs as bold text

*--rainbow* and *--aurora* filters with colors an ASCII picture of a bone saying something

*--super* ...enjoy!

If the program is invoked as *bonethink* then the bone will think its message instead of saying it.

BONEFILE FORMAT
--------------
A bonefile is made up of a simple block of *perl(1)* code, which assigns a picture of a bone to the variable *$the_bone*.
Should you wish to customize the eyes or the tongue of the bone,
then the variables *$eyes* and *$tongue* may be used. The trail leading up to the bone's message balloon is
composed of the character(s) in the *$thoughts* variable. Any backslashes must be reduplicated to prevent interpolation.
The name of a bonefile should end with *.bone ,* otherwise it is assumed not to be a bonefile. Also, at-signs (``@'')
must be backslashed because that is what Perl 5 expects.

ENVIRONMENT
-----------
The BONEPATH environment variable, if present, will be used to search
for bonefiles.  It contains a colon-separated list of directories,
much like *PATH or MANPATH*. It should always contain the */usr/local/share/bones*
directory, or at least a directory with a file called *default.bone* in it.

FILES
-----
*%PREFIX%/share/bones* holds a sample set of bonefiles. If your *BONEPATH* is not explicitly set, it automatically contains this directory.

BUGS
----
https://github.com/anthonycuervo23/bonesay

If there are any, please report bugs and feature requests in the issue tracker.
Please do your best to provide a reproducible test case for bugs. This should
include the *bonesay* command, the actual output and the expected output.

AUTHORS
-------
Neo-bonesay author is Kei Kamikawa (x00.x7f.x86@gmail.com).

The original author is Tony Monroe (tony@nog.net), with suggestions from Shannon
Appel (appel@CSUA.Berkeley.EDU) and contributions from Anthony Polito
(aspolito@CSUA.Berkeley.EDU).

SEE ALSO
--------
perl(1), wall(1), nwrite(1), figlet(6)