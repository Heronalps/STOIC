#!/bin/bash

FILE=$1

SZ=`wc -l $FILE | awk '{print $1}'`
CNT=1

while ( test $CNT -le $SZ ) ; do
	grep -v "epoch" $FILE | awk -F ',' '{print $3,$2,$4,$8,$1}' | head -n $(($CNT+3)) | tail -n 4 > zzz111
	MCNT=1
	MIN=999999999.0
	AMIN=999999999.0
	BSIZE=`head -n 1 zzz111 | awk '{print $2}'`
	EPOCH=`head -n 1 zzz111 | awk '{print $5}'`
	while ( test $MCNT -le 4) ; do
		LINE=`head -n $MCNT zzz111 | tail -n 1`
		PRED=`echo $LINE | awk '{print $3}'` 
		ACTUAL=`echo $LINE | awk '{print $4}'` 

		MTEST=`echo $MIN $PRED | awk '{if($2 < $1) {print 1}else{print 0}}'` 
		if ( test $MTEST -eq 1 ) ; then
			MIN=$PRED
			METHOD=`echo $LINE | awk '{print $1}'`
		fi
		ATEST=`echo $AMIN $ACTUAL | awk '{if($2 < $1) {print 1}else{print 0}}'` 
		if ( test $ATEST -eq 1 ) ; then
			AMIN=$ACTUAL
			AMETHOD=`echo $LINE | awk '{print $1}'`
		fi
		MCNT=$(($MCNT+1))
	done
	if ( test "$METHOD" != "$AMETHOD" ) ; then
		LTEST=`echo $MIN $AMIN | awk '{if($2 < $1) {print 1}else{print 0}}'` 
		if ( test $LTEST -eq 0 ) ; then
		DIFF=`echo $MIN $AMIN | awk '{if($1 < $2){print $2-$1}else{print $1-$2}}'`
		PCT=`echo $DIFF $AMIN | awk '{print $1/$2}'`
		echo $BSIZE $MIN $METHOD $AMIN $AMETHOD $DIFF $PCT $EPOCH
		fi
	fi
	CNT=$(($CNT+4))
done