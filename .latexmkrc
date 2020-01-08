#!/usr/bin/env perl
$lualatex                    = 'lualatex -halt-on-error -synctex=1 -shell-escape';
$bibtex                      = 'pbibtex -kanji=utf8 %O %B';
$biber                       = 'biber --bblencoding=utf8 -u -U --output_safechars';
$makeindex                   = 'mendex %O -o %D %S';
$max_repeat                  = 5;
$pvc_view_file_via_temporary = 0;
$pdf_mode                    = 4;
$postscript_mode = $dvi_mode = 0;