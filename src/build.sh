# Crawl site
wget --recursive --no-host-directories http://localhost:3100/index.html -o src/wget.log

# Now deal with favicons that are searched for by browsers
cp -r src/images/fav/* images/fav/.
cp src/browserconfig.xml .

cd docs
for f in `ls`
do
    mv ${f} ${f}.html
done

for f in `find . -type f | grep -v "^\./src" | grep -v "^\./\.git"`
do
    git add ${f}
done

