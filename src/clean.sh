for i in `ls | egrep -v ^src`
do
    rm -rf ${i}
done
