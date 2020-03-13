#!/bin/bash

FILE=$1

SZ=`wc -l $FILE | awk '{print $1}'`
CNT=1

while ( test $CNT -le $SZ ) ; do
	grep -v "epoch" $FILE | awk -F ',' '{print $3,$2,$4,$8,$1,$7,$9,$10,$11}' | head -n $(($CNT+3)) | tail -n 4 > zzz222
	MCNT=1
	OMIN=99999999.9
	TMIN=9999999.9
	while ( test $MCNT -le 4) ; do
		LINE=`head -n $MCNT zzz222 | tail -n 1`
		TYPE=`echo $LINE | awk '{print $1}'`
		if ( test "$TYPE" = "cpu" ) ; then
			MCNT=$(($MCNT+1))
			continue
		fi
		if ( test "$TYPE" = "gpu1" ) ; then
			MCNT=$(($MCNT+1))
			continue
		fi
		PRED=`echo $LINE | awk '{print $3}'` 
		ACT=`echo $LINE | awk '{print $4}'` 
# get edge prediction time and assume we are doing edge
		if ( test "$TYPE" = "edge" ) ; then
			EDGEPRED=$PRED
			EDGEACT=$ACT
		fi
# record the minimum predicted time as the one we report to the user
		MTEST=`echo "$PRED $OMIN" | awk '{if($1 < $2){print 1}else{print 0}}'`
		if ( test $MTEST -eq 1 ) ; then
			OMIN=$PRED
			OMETHOD=$TYPE
			OACTUAL=$ACT
		fi
# record the actual minimum
		MTEST=`echo $TMIN $ACT | awk '{if($1 > $2){print 1}else{print 0}}'`
		if ( test $MTEST -eq 1) ; then
			TMIN=$ACT
			TMETHOD=$TYPE
		fi
		MCNT=$(($MCNT+1))
	done

# start out running on the edge
	MCNT=1
# current prediction is EDGE
	PRED=$EDGEPRED
# the actual min for the edge is EDGEACT
	AMIN=$EDGEACT
	AMETHOD="edge"
	APRED=$EDGEPRED
	RMAIN=0
	PPTIME=0
	BSIZE=`head -n 1 zzz222 | awk '{print $2}'`
	EPOCH=`head -n 1 zzz222 | awk '{print $5}'`
	while ( test $MCNT -le 4) ; do
		LINE=`head -n $MCNT zzz222 | tail -n 1`
		PRED=`echo $LINE | awk '{print $3}'` 
		ACTUAL=`echo $LINE | awk '{print $4}'` 
		TYPE=`echo $LINE | awk '{print $1}'`
		if ( test "$TYPE" = "cpu" ) ; then
			MCNT=$(($MCNT+1))
			continue
		fi
		if ( test "$TYPE" = "gpu1" ) ; then
			MCNT=$(($MCNT+1))
			continue
		fi
		if ( test "$TYPE" != "edge" ) ; then
# actual start time for non edge
			STIME=`echo $LINE | awk '{print $7+$8}'`
# actual run time for non edge 
			RTIME=`echo $LINE | awk '{print $9}'`
# predicted run time for non edge
			PRTIME=`echo $LINE | awk '{print $6}'`
# predicted remaining time for edge
			REMAIN=`echo $APRED $STIME | awk '{print $1 - $2}'`
			ACTREMAIN=`echo $AMIN $STIME | awk '{print $1-$2}'`
echo "$BSIZE $TYPE |  $AMETHOD remain: $REMAIN ($PRED $STIME) actremain: $ACTREMAIN $TYPE pred run: $PRTIME"
#echo "$BSIZE reported: $OMETHOD $OMIN bestsofar: $AMETHOD $AMIN best: $TMETHOD $TMIN"

			TTEST=`echo "$ACTREMAIN" | awk '{if($1 <= 0) {print 1}else{print 0}}'` 
# have we finished already?
			if ( test $TTEST -eq 1 ) ; then
				MCNT=$(($MCNT+1))
				continue
			fi
# do we predict we should switch?
			TTEST=`echo "$REMAIN $PRTIME" | awk '{if($1 <= $2) {print 1}else{print 0}}'` 
# if predictd remaining time on edge is bigger than predicted run time for
# non edge
			if ( test $TTEST -eq 0 ) ; then
# record prediction of time to completion for non edge from this point
				APRED=$PRED
				AMIN=$ACTUAL
				AMETHOD=$TYPE
			fi
		fi
		MCNT=$(($MCNT+1))
	done

# did the edge finish it before the time given to the user
	FTEST=`echo $OACTUAL $EDGEACT | awk '{if($1 > $2){print 1}else{print 0}}'`
	if ( test $FTEST -eq 1 ) ; then
		DIFF=`echo $OACTUAL $EDGEACT | awk '{if($1 < $2){print $2-$1}else{print $1-$2}}'`
		echo "EDGE SUCCESS: $BSIZE predicted: $OMETHOD $OMIN $OACTUAL edge: $EDGEACT diff: $DIFF"
		CNT=$(($CNT+4))
		continue
	fi

# did the edge finish before the STOIC time
	FTEST=`echo $AMIN $EDGEACT | awk '{if($1 > $2){print 1}else{print 0}}'`
	if ( test $FTEST -eq 1 ) ; then
		DIFF=`echo $AMIN $EDGEACT | awk '{if($1 < $2){print $2-$1}else{print $1-$2}}'`
		echo "CHOSE SUCCESS: $BSIZE switched: $OMETHOD $OMIN $OACTUAL edge: $EDGEACT chosen: $AMETHOD $AMIN diff: $DIFF"
		AMIN=$EDGEACT
		AMETHOD="edge"
		OMETHOD="edge"
		OMIN=$EDGEACT
		OACTUAL=$EDGEACT
	fi
	
# if we swiitched
	if ( test "$OMETHOD" != "$AMETHOD" ) ; then
		DIFF=`echo $OACTUAL $AMIN | awk '{if($1 < $2){print $2-$1}else{print $1-$2}}'`
# if actual time the user would have got is bigger that what happened, we failed to make it better
		LTEST=`echo $OACTUAL $AMIN | awk '{if($2 > $1) {print 1}else{print 0}}'` 
		if ( test $LTEST -eq 1 ) ; then
			echo "FAIL: $BSIZE reported: $OMETHOD $OMIN $OACTUAL actual: $AMETHOD $AMIN diff: $DIFF $EPOCH"
		else
			echo "SUCCESS: $BSIZE predicted: $OMETHOD $OMIN $OACTUAL actual: $AMETHOD $AMIN diff: $DIFF"
		fi 
	else
# if the actual of the predicted > than the true min, we failed
		DIFF=`echo $OACTUAL $TMIN | awk '{if($1 < $2){print $2-$1}else{print $1-$2}}'`
		LTEST=`echo $OACTUAL $TMIN | awk '{if($1 <= $2) {print 1}else{print 0}}'`
		if ( test $LTEST -eq 0 ) ; then
			echo "NOSWITCH FAIL: $BSIZE reported: $OMETHOD $OMIN $OACTUAL actual: $TMETHOD $TMIN diff: $DIFF $EPOCH"
		else
			echo "NOSWITCH SUCCESS: $BSIZE predicted: $OMETHOD $OMIN $OACTUAL actual: $TMETHOD $TMIN diff: $DIFF"
		fi
	fi
	CNT=$(($CNT+4))
done
	

