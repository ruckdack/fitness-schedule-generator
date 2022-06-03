echo generating...
html_name="abcdefghijklmn.html"
cat config.json | go run . > $html_name
google-chrome --headless --print-to-pdf-no-header --print-to-pdf=plan.pdf --no-margins $html_name
rm $html_name
echo done!
