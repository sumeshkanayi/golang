# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include ../../Make.inc
O:=$(HOST_O)

TARG=8a

HFILES=\
	a.h\
	y.tab.h\
	../8l/8.out.h\

OFILES=\
	y.tab.$O\
	lex.$O\
	../8l/enam.$O\

YFILES=\
	a.y\

include ../../Make.ccmd

lex.$O:	../cc/macbody ../cc/lexbody
