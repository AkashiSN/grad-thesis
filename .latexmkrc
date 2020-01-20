#!/usr/bin/env perl
$latex     = "find . -type f -name '*.tex' -print0 | xargs -0 sed -i -e 's/、/，/g' -e 's/。/．/g' && uplatex %O -synctex=1 -halt-on-error -interaction=nonstopmode -shell-escape -file-line-error %S";
$pdflatex  = "find . -type f -name '*.tex' -print0 | xargs -0 sed -i -e 's/、/，/g' -e 's/。/．/g' && pdflatex %O -synctex=1 -halt-on-error -interaction=nonstopmode -shell-escape -file-line-error %S";
$lualatex  = "find . -type f -name '*.tex' -print0 | xargs -0 sed -i -e 's/、/，/g' -e 's/。/．/g' && lualatex %O -synctex=1 -halt-on-error -interaction=nonstopmode -shell-escape -file-line-error %S";
$xelatex   = "find . -type f -name '*.tex' -print0 | xargs -0 sed -i -e 's/、/，/g' -e 's/。/．/g' && xelatex %O -synctex=1 -halt-on-error -interaction=nonstopmode -shell-escape -file-line-error  %S";
$biber     = 'biber %O --bblencoding=utf8 -u -U --output_safechars %B';
$bibtex    = 'upbibtex -kanji=utf8 %O %B';
$makeindex = 'upmendex %O -o %D %S';
$dvipdf    = 'dvipdfmx %O -o %D %S';
$dvips     = 'dvips %O -z -f %S | convbkmk -u > %D';
$ps2pdf    = 'ps2pdf %O %S %D';
$pdf_mode  = 3;
$out_dir   = "out"