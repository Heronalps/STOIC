#!/bin/bash

FILE=$1

SZ=`wc -l $FILE | awk '{print $1}'`
CNT=1

while ( test $CNT -le $SZ ) ; do
	grep -v "epoch" $FILE | awk -F ',' '{print $3,$2,$4,$8,$1,$7,$9,$10,$11}' | head -n $(($CNT+3)) | tail -n 4 > zzz222
	MCNT=1
	while ( test $MCNT -le 4) ; do
		LINE=`head -n $MCNT zzz222 | tail -n 1`
		TYPE=`echo $LINE | awk '{print $1}'`
		PRED=`echo $LINE | awk '{print $3}'` 
		ACT=`echo $LINE | awk '{print $4}'` 
		if ( test "$TYPE" = "edge" ) ; then
			ETIME=$PRED
			AMIN=$ACT
			break
		fi
		MCNT=$(($MCNT+1))
	done

	MCNT=1
	MIN=$ETIME
	METHOD="edge"
	AMETHOD="edge"
	RMAIN=0
	PPTIME=0
	ELAPSED=$ETIME
	BSIZE=`head -n 1 zzz222 | awk '{print $2}'`
	EPOCH=`head -n 1 zzz222 | awk '{print $5}'`
	while ( test $MCNT -le 4) ; do
		LINE=`head -n $MCNT zzz222 | tail -n 1`
		PRED=`echo $LINE | awk '{print $3}'` 
		ACTUAL=`echo $LINE | awk '{print $4}'` 
		TYPE=`echo $LINE | awk '{print $1}'`
		if ( test "$TYPE" != "edge" ) ; then
			STIME=`echo $LINE | awk '{print $7+$8}'`
			RTIME=`echo $LINE | awk '{print $9}'`
			PRTIME=`echo $LINE | awk '{print $6}'`
			REMAIN=`echo $ETIME $STIME | awk '{print $1 - $2}'`
			TTEST=`echo "$REMAIN $PRTIME" | awk '{if($1 < $2) {print 1}else{print 0}}'` 
			if ( test $TTEST -eq 0 ) ; then
				PRED=`echo $STIME $PRTIME | awk '{print $1+$2}'`
				MTEST=`echo $MIN $PRED | awk '{if($2 < $1) {print 1}else{print 0}}'` 
				if ( test $MTEST -eq 1 ) ; then
					MIN=$PRED
					RMAIN=$REMAIN
					PPTIME=$PRTIME
					ELAPSED=$STIME
					METHOD=`echo $LINE | awk '{print $1}'`
				fi
			fi
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
#		echo "$BSIZE predicted: $METHOD $MIN actual: $AMETHOD $AMIN etime: $ETIME elapsed: $ELAPSED remain: $RMAIN predrtime: $PPTIME $EPOCH"
		echo "$BSIZE predicted: $METHOD $MIN actual: $AMETHOD $AMIN diff: $DIFF"
		fi
	fi
	CNT=$(($CNT+4))
done
	

