@echo off
echo generating...
set html_name=abcdefghijklmn.html
more config.json | go run . > %html_name%
set html_loc=%cd%\%html_name%
set pdf_loc=%cd%\plan.pdf
"C:\Program Files\Google\Chrome\Application\chrome.exe" --headless --print-to-pdf-no-header --print-to-pdf=%pdf_loc% --no-margins %html_loc%
del %html_name%
echo done!